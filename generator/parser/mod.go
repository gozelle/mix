package parser

import (
	"fmt"
	"github.com/gozelle/fs"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func PrepareMod() (mod *Mod, err error) {
	var parent string
	parent, err = os.Getwd()
	if err != nil {
		// A nonexistent working directory can't be in a module.
		return
	}
	
	defer func() {
		if err == nil {
			if mod.file.Module == nil {
				err = fmt.Errorf("mod file not defined module")
				return
			}
		}
	}()
	
	for {
		var f *os.File
		if f, err = os.Open(fs.Join(parent, "go.mod")); err == nil {
			var d []byte
			d, err = io.ReadAll(f)
			if err != nil {
				return
			}
			_ = f.Close()
			if err == nil || err == io.EOF {
				// go.mod exists and is readable (is a file, not a directory).
				var mf *modfile.File
				mf, err = modfile.Parse("go.mod", d, nil)
				if err != nil {
					return
				}
				mod = &Mod{file: mf, root: parent}
				break
			}
		}
		d := filepath.Dir(parent)
		if len(d) >= len(parent) {
			return
		}
		parent = fs.Join(parent, "../")
	}
	return
}

type Mod struct {
	root     string
	file     *modfile.File
	packages map[string]*Package
	loaded   map[string]bool
}

func (m Mod) ModuleName() string {
	if m.file == nil {
		return ""
	}
	return m.file.Module.Mod.Path
}

func (m Mod) Gopath() string {
	return os.Getenv("GOPATH")
}

func (m *Mod) cachePackage(realPath string, pkg *Package) {
	if m.packages == nil {
		m.packages = map[string]*Package{}
	}
	
	m.packages[realPath] = pkg
}

// GetPackagePath 获取包的真实路径
// 1. 首先判断是否被本地替换
// 2. 然后判断是否是直接依赖的包
// 3. 最后判断是否为系统包
func (m Mod) GetPackagePath(pkg string) string {
	if m.file == nil {
		return ""
	}
	if m.ModuleName() != "std" {
		for _, v := range m.file.Require {
			if v.Mod.Path == pkg {
				return fs.Join(m.Gopath(), "pkg/mod", fmt.Sprintf("%s@%s", pkg, v.Mod.Version))
			} else {
				c := pkg
				valid := false
				for c != "." && c != "/" {
					if v.Mod.Path == c {
						valid = true
						break
					}
					c = filepath.Join(c, "../")
				}
				if valid {
					// 如果有根项目定义包版本升级
					for _, pv := range m.file.Require {
						if pv.Mod.Path == c {
							return fs.Join(m.Gopath(), "pkg/mod", fmt.Sprintf("%s@%s%s", c, pv.Mod.Version, strings.TrimPrefix(pkg, c)))
						}
					}
					return fs.Join(m.Gopath(), "pkg/mod", fmt.Sprintf("%s@%s%s", c, v.Mod.Version, strings.TrimPrefix(pkg, c)))
				}
			}
		}
		for _, v := range m.file.Replace {
			if v.New.Path == pkg {
				return v.Old.Path
			} else {
				c := pkg
				valid := false
				for c != "." && c != "/" {
					if v.New.Path == c {
						valid = true
						break
					}
					c = filepath.Join(c, "../")
				}
				if valid {
					return fmt.Sprintf("%s%s", v.Old.Path, strings.TrimPrefix(pkg, c))
				}
			}
		}
	}
	
	path := fs.Join(m.Gopath(), "src", pkg)
	if fs.Exists(path) {
		return path
	}
	
	path = fs.Join(m.Gopath(), "src/vendor", pkg)
	if fs.Exists(path) {
		return path
	}
	
	path = fs.Join(m.root, strings.TrimPrefix(strings.TrimPrefix(pkg, m.file.Module.Mod.Path), "/"))
	if fs.Exists(path) {
		return path
	}
	
	return ""
}

func (m Mod) GetPackageRealName(path string) (name string, err error) {
	
	files, err := m.GetPackageFiles(path)
	if err != nil {
		return
	}
	set := token.NewFileSet()
	
	var f *ast.File
	for _, v := range files {
		if strings.HasSuffix(v, "_test.go") {
			continue
		}
		
		f, err = parser.ParseFile(set, v, nil, parser.AllErrors|parser.ParseComments)
		if err != nil {
			return
		}
		if f.Name == nil {
			err = fmt.Errorf("package name is nil")
			return
		}
		name = f.Name.String()
		
		// ignore main namespace
		if name != "main" {
			return
		}
	}
	
	return
}

func (m Mod) GetPackageFiles(path string) (files []string, err error) {
	
	//var path string
	//if strings.HasPrefix(pkg, m.ModuleName()) {
	//	path = fs.Join(m.root, strings.TrimPrefix(strings.TrimPrefix(pkg, m.ModuleName()), "/"))
	//} else {
	//
	//}
	
	if !fs.Exists(path) {
		return
	}
	
	files, err = fs.Files(path, ".go")
	if err != nil {
		return
	}
	
	return
}
