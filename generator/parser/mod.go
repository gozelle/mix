package parser

import (
	"fmt"
	"github.com/gozelle/fs"
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
		if f, err = os.Open(filepath.Join(parent, "go.mod")); err == nil {
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
		parent = filepath.Join(parent, "../")
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

func (m *Mod) AddPackage(pkg *Package) {
	if m.packages == nil {
		m.packages = map[string]*Package{}
	}
	if pkg.Path != "" {
		m.packages[pkg.Path] = pkg
	}
}

func (m Mod) GetPackage(path string) *Package {
	if m.packages == nil {
		return nil
	}
	return m.packages[path]
}

// GetPackagePath 获取包的真实路径
// 1. 首先判断是否被本地替换
// 2. 然后判断是否是直接依赖的包
// 3. 最后判断是否为系统包
func (m Mod) GetPackagePath(pkg string) string {
	if m.file == nil {
		return ""
	}
	for _, v := range m.file.Replace {
		if v.New.Path == pkg {
			return v.Old.Path
		}
	}
	for _, v := range m.file.Require {
		if !v.Indirect && v.Mod.Path == pkg {
			return filepath.Join(m.Gopath(), "pkg/mod", fmt.Sprintf("%s@%s", pkg, v.Mod.Version))
		}
	}
	
	src := filepath.Join(m.Gopath(), "src", pkg)
	if err := fs.Exists(src); err == nil {
		return src
	}
	
	src = filepath.Join(m.Gopath(), "src/vendor", pkg)
	if err := fs.Exists(src); err == nil {
		return src
	}
	
	return filepath.Join(m.root, strings.TrimPrefix(pkg, m.file.Module.Mod.Path))
}

func (m Mod) GetPackageRealName(pkg string) (name string, err error) {
	
	files, err := m.GetPackageFiles(pkg)
	if err != nil {
		return
	}
	set := token.NewFileSet()
	
	file := ""
	for _, v := range files {
		if !strings.HasSuffix(v, "_test.go") {
			file = v
		}
	}
	
	f, err := parser.ParseFile(set, file, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return
	}
	if f.Name == nil {
		err = fmt.Errorf("package name is nil")
		return
	}
	name = f.Name.String()
	
	return
}

func (m Mod) GetPackageFiles(pkg string) (files []string, err error) {
	
	var path string
	if strings.HasPrefix(pkg, m.ModuleName()) {
		path = filepath.Join(m.root, strings.TrimPrefix(pkg, m.ModuleName()))
		
	} else {
		path = m.GetPackagePath(pkg)
	}
	
	if path == "" {
		err = fmt.Errorf("can't resolve pkg path: %s", pkg)
		return
	}
	
	files, err = fs.Files(path, ".go")
	if err != nil {
		return
	}
	
	return
}
