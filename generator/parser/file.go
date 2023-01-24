package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

type File struct {
	mod     *Mod
	pkg     *Package
	Imports map[string]*Import
}

func (f *File) getImport(name string) *Import {
	if f.Imports == nil {
		return nil
	}
	return f.Imports[name]
}

func (f *File) addImport(item *Import) *Import {
	if f.Imports == nil {
		f.Imports = map[string]*Import{}
	}
	
	if v, ok := f.Imports[item.Alias]; ok {
		return v
	}
	f.Imports[item.Alias] = item
	return item
}

func (f *File) Visit(node ast.Node) ast.Visitor {
	s, ok := node.(*ast.TypeSpec)
	if !ok {
		return f
	}
	
	switch t := s.Type.(type) {
	
	case *ast.InterfaceType:
		f.pkg.addInterface(s.Name.String(), &Interface{Name: s.Name.String(), interfaceType: t})
	case *ast.FuncType:
		return f
	case *ast.StructType, *ast.Ident:
		f.pkg.addType(s.Name.String(), parseType(s.Name.String(), s.Type))
	//case *ast.MapType:
	//	// TODO
	//case *ast.SliceExpr:
	//	panic(s.Name.String())
	//case *ast.ArrayType:
	//	panic(s.Name.String())
	case *ast.SelectorExpr:
		f.pkg.addType(s.Name.String(), parseType(s.Name.String(), s.Type))
	default:
		panic(fmt.Errorf("unsupport parse type: %s", reflect.TypeOf(t)))
	}
	
	return f
}

func (f *File) load(mod *Mod, file string) (err error) {
	
	set := token.NewFileSet()
	af, err := parser.ParseFile(set, file, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return
	}
	
	if f.pkg.Name == "" {
		f.pkg.Name = af.Name.String() // package name
	}
	
	for _, i := range af.Imports {
		f.addImport(f.parseImport(mod, i))
	}
	// parse file use Visit
	ast.Walk(f, af)
	
	return
}

func (f *File) parseAlias(name string) string {
	if name == "<nil>" {
		return ""
	}
	return name
}

func (f *File) parseImport(mod *Mod, i *ast.ImportSpec) *Import {
	
	r := &Import{
		Alias:   f.parseAlias(i.Name.String()),
		Path:    strings.Trim(i.Path.Value, "\""),
		Package: &Package{},
	}
	if r.Alias == "" {
		var err error
		r.Alias, err = mod.GetPackageRealName(r.Path)
		if err != nil {
			panic(fmt.Errorf("get %s package name error: %s", r.Path, err))
		}
	}
	
	files, err := mod.GetPackageFiles(r.Path)
	if err != nil {
		panic(fmt.Errorf("get package: %s fiels error: %s", r.Path, err))
	}
	
	err = r.Package.loadFiles(mod, files)
	if err != nil {
		panic(fmt.Errorf("load package: %s fiels error: %s", r.Path, err))
	}
	
	return r
}
