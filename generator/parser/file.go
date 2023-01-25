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
	path    string
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
		f.pkg.addInterface(s.Name.String(), &Interface{Name: s.Name.String(), interfaceType: t, file: f})
		f.pkg.addType(s.Name.String(), parseType(s.Name.String(), s.Type))
	case *ast.FuncType:
		return f
	case *ast.StructType, *ast.Ident:
		f.pkg.addType(s.Name.String(), parseType(s.Name.String(), s.Type))
	case *ast.MapType:
		// TODO
		f.pkg.addType(s.Name.String(), parseType(s.Name.String(), s.Type))
	
	case *ast.SliceExpr:
		// TODO
		f.pkg.addType(s.Name.String(), parseType(s.Name.String(), s.Type))
	
	case *ast.ArrayType:
		f.pkg.addType(s.Name.String(), parseType(s.Name.String(), s.Type))
		
		// TODO
	case *ast.StarExpr:
		// TODO
		f.pkg.addType(s.Name.String(), parseType(s.Name.String(), s.Type))
	
	case *ast.SelectorExpr:
		f.pkg.addType(s.Name.String(), parseType(s.Name.String(), s.Type))
	case *ast.ChanType:
		// TODO
	default:
		panic(fmt.Errorf("unsupport parse type: %s", reflect.TypeOf(t)))
	}
	
	return f
}

func (f *File) load(file string) (err error) {
	log.Debugf("load file: %s", file)
	set := token.NewFileSet()
	af, err := parser.ParseFile(set, file, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return
	}
	
	if f.pkg.Name == "" {
		f.pkg.Name = af.Name.String() // package name
	}
	
	for _, i := range af.Imports {
		f.addImport(f.parseImport(i))
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

func (f *File) parseImport(i *ast.ImportSpec) *Import {
	
	r := &Import{
		Alias:   f.parseAlias(i.Name.String()),
		Path:    strings.Trim(i.Path.Value, "\""),
		Package: &Package{},
	}
	
	if r.Path == "C" {
		return r
	}
	
	if r.Alias == "" {
		var err error
		r.Alias, err = f.mod.GetPackageRealName(r.Path)
		if err != nil {
			panic(fmt.Errorf("%s: get %s package name error: %s", f.path, r.Path, err))
		}
	}
	log.Debugf("load import: %s", r.Path)
	realPath := f.mod.GetPackagePath(r.Path)
	//log.Debugf("load import path: %s", realPath)
	err := r.Package.load(f.mod, realPath)
	if err != nil {
		panic(fmt.Errorf("load package: %s files from: %s error: %s", r.Path, realPath, err))
	}
	
	return r
}
