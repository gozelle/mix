package parser

import (
	"fmt"
	"github.com/gozelle/logging"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

var log = logging.Logger("parser")

type Package struct {
	Name       string
	Alias      string
	Path       string
	Imports    map[string]*Package
	Interfaces map[string]*Interface
	Defs       map[string]*Def
	loaded     bool
}

func (p *Package) getDef(name string) *Def {
	if p.Defs == nil {
		return nil
	}
	return p.Defs[name]
}

func (p *Package) getImport(name string) *Package {
	if p.Imports == nil {
		return nil
	}
	return p.Imports[name]
}

func (p *Package) addImport(item *Package) *Package {
	if p.Imports == nil {
		p.Imports = map[string]*Package{}
	}
	name := item.Name
	if item.Alias != "" && item.Alias != "nil" {
		name = p.Alias
	}
	if v, ok := p.Imports[name]; ok {
		return v
	}
	p.Imports[name] = item
	return item
}

func (p *Package) addType(name string, item *Def) *Def {
	if p.Defs == nil {
		p.Defs = map[string]*Def{}
	}
	if v, ok := p.Defs[name]; ok {
		return v
	}
	p.Defs[name] = item
	return item
}

func (p *Package) addInterface(name string, item *Interface) *Interface {
	if p.Interfaces == nil {
		p.Interfaces = map[string]*Interface{}
	}
	if v, ok := p.Interfaces[name]; ok {
		return v
	}
	p.Interfaces[name] = item
	return item
}

func (p *Package) getInterface(name string) *Interface {
	if p.Interfaces == nil {
		return nil
	}
	v, ok := p.Interfaces[name]
	if ok {
		return v
	}
	return nil
}

func (p *Package) loadFiles(mod *Mod, files []string) (err error) {
	for _, v := range files {
		if !strings.HasSuffix(v, "_test.go") {
			err = p.parseFile(mod, v)
			if err != nil {
				return
			}
		}
	}
	p.loaded = true
	return
}

func (p *Package) parseFile(mod *Mod, file string) (err error) {
	
	set := token.NewFileSet()
	f, err := parser.ParseFile(set, file, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return
	}
	
	p.Name = f.Name.String() // package name
	
	for _, i := range f.Imports {
		p.addImport(p.parseImport(mod, i))
	}
	// parse file use Visit
	ast.Walk(p, f)
	
	return
}

func (p *Package) parseAlias(name string) string {
	if name == "<nil>" {
		return ""
	}
	return name
}

func (p *Package) parseImport(mod *Mod, i *ast.ImportSpec) *Package {
	
	r := &Package{
		Alias: p.parseAlias(i.Name.String()),
		Path:  strings.Trim(i.Path.Value, "\""),
	}
	var err error
	r.Name, err = mod.GetPackageRealName(r.Path)
	if err != nil {
		panic(fmt.Errorf("get %s package name error: %s", r.Path, err))
	}
	
	return r
}

func (p *Package) Visit(node ast.Node) ast.Visitor {
	s, ok := node.(*ast.TypeSpec)
	if !ok {
		return p
	}
	
	switch t := s.Type.(type) {
	
	case *ast.InterfaceType:
		p.addInterface(s.Name.String(), &Interface{Name: s.Name.String(), interfaceType: t})
	case *ast.FuncType:
		return p
	case *ast.StructType, *ast.Ident:
		p.addType(s.Name.String(), parseType(s.Name.String(), s.Type))
	case *ast.MapType:
		// TODO
	case *ast.SliceExpr:
		panic(s.Name.String())
	case *ast.ArrayType:
		panic(s.Name.String())
	case *ast.SelectorExpr:
		p.addType(s.Name.String(), parseType(s.Name.String(), s.Type))
	default:
		panic(fmt.Errorf("unsupport parse type: %s", reflect.TypeOf(t)))
	}
	
	return p
}

func (p *Package) getPackage(mod *Mod, name string) *Package {
	if p.Imports == nil {
		p.Imports = map[string]*Package{}
	}
	v, ok := p.Imports[name]
	if !ok {
		return nil
	}
	if !v.loaded {
		files, err := mod.OpenPackage(v.Path)
		if err != nil {
			panic(fmt.Errorf("open package '%s' error: %s", v.Path, err))
		}
		err = v.loadFiles(mod, files)
		if err != nil {
			panic(fmt.Errorf("load package '%s' error: %s", v.Path, err))
		}
	}
	
	return v
}
