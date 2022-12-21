package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

type Package struct {
	Name       string
	Alias      string
	Path       string
	Imports    map[string]*Package
	Interfaces map[string]*Interface
	Defs       map[string]*Def
	mod        *Mod
	loaded     bool
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

func (p *Package) loadPackage(files []string) (err error) {
	for _, v := range files {
		if !strings.HasSuffix(v, "_test.go") {
			err = p.parseFile(v)
			if err != nil {
				return
			}
		}
	}
	p.loaded = true
	return
}

func (p *Package) parseFile(file string) (err error) {
	//defer func() {
	//	if e := recover(); e != nil {
	//		err = fmt.Errorf("%v", e)
	//		return
	//	}
	//}()
	
	set := token.NewFileSet()
	f, err := parser.ParseFile(set, file, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return
	}
	
	p.Name = f.Name.String()
	for _, i := range f.Imports {
		p.addImport(p.parseImport(i))
	}
	ast.Walk(p, f)
	
	//for _, v := range p.Interfaces {
	//	p.parseInterface(v)
	//}
	
	return
}

func (p *Package) parseAlias(name string) string {
	if name == "<nil>" {
		return ""
	}
	return name
}

func (p *Package) parseImport(i *ast.ImportSpec) *Package {
	
	r := &Package{
		Alias: p.parseAlias(i.Name.String()),
		Path:  strings.Trim(i.Path.Value, "\""),
		mod:   p.mod,
	}
	var err error
	r.Name, err = r.mod.GetPackageRealName(r.Path)
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
		p.addInterface(s.Name.String(), &Interface{Name: s.Name.String(), t: t})
	case *ast.FuncType:
		return p
	case *ast.StructType, *ast.Ident:
		p.addType(s.Name.String(), parseType(s.Name.String(), s.Type))
	case *ast.MapType:
		// TODO
	case *ast.SliceExpr:
		// TODO
	case *ast.ArrayType:
	// TODO
	case *ast.SelectorExpr:
		p.addType(s.Name.String(), parseType(s.Name.String(), s.Type))
	default:
		panic(fmt.Errorf("unsupport parse type: %s", reflect.TypeOf(t)))
	}
	
	return p
}

func (p *Package) getPackage(name string) *Package {
	if p.Imports == nil {
		p.Imports = map[string]*Package{}
	}
	v, ok := p.Imports[name]
	if !ok {
		return nil
	}
	if !v.loaded {
		files, err := p.mod.OpenPackage(v.Path)
		if err != nil {
			panic(fmt.Errorf("open package '%s' error: %s", v.Path, err))
		}
		err = v.loadPackage(files)
		if err != nil {
			panic(fmt.Errorf("load package '%s' error: %s", v.Path, err))
		}
	}
	
	return v
}
