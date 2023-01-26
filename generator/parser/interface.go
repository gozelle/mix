package parser

import (
	"fmt"
	"go/ast"
	"reflect"
)

type Interface struct {
	Name          string
	Methods       []*Method
	Defs          []*Def
	defs          map[string]bool
	Imports       []*Import
	imports       map[string]bool
	includes      []*Interface
	interfaceType *ast.InterfaceType
	file          *File
}

func (i *Interface) load(mod *Mod, pkg *Package, file *File) (err error) {
	for _, m := range i.interfaceType.Methods.List {
		switch mt := m.Type.(type) {
		case *ast.Ident:
		// TODO parse include
		case *ast.FuncType:
			i.Methods = append(i.Methods, i.parseMethod(mod, pkg, file, m.Names[0].Name, mt))
		}
	}
	for _, v := range pkg.Defs {
		i.addDef(v)
	}
	return
}

func (i *Interface) parseMethod(mod *Mod, pkg *Package, file *File, name string, t *ast.FuncType) (r *Method) {
	r = &Method{Name: name}
	
	if t.Params != nil {
		for _, f := range t.Params.List {
			r.Params = append(r.Params, i.parseParam(mod, pkg, file, parseNames(f.Names), f.Type))
		}
	}
	
	if t.Results != nil {
		for _, f := range t.Results.List {
			r.Results = append(r.Results, i.parseParam(mod, pkg, file, parseNames(f.Names), f.Type))
		}
	}
	return
}

func (i *Interface) addDef(t *Def) {
	if i.defs == nil {
		i.defs = map[string]bool{}
	}
	if _, ok := i.defs[t.Name]; ok {
		return
	}
	i.defs[t.Name] = true
	i.Defs = append(i.Defs, t)
}

func (i *Interface) parseParam(mod *Mod, pkg *Package, file *File, names []string, t ast.Expr) (r *Param) {
	
	r = &Param{
		Names: names,
	}
	
	//defer func() {
	//	log.Infof("param: %v Type: %s", r.Names, r.Type)
	//}()
	
	switch e := t.(type) {
	case *ast.Ident:
		r.Type = parseType(i.file, "", e)
	case *ast.SelectorExpr:
		
		//pkgName := e.X.(*ast.Ident).String()
		//typ := e.Sel.Name
		r.Type = parseType(i.file, "", e) // TODO
		//r.Type = &Type{Name: fmt.Sprintf("%s.%s", pkgName, typ)}
		//// 去 Parser 上下文中寻找对应的 Package 进来
		//imt := file.getImport(pkgName)
		//
		//if imt != nil {
		//	i.addImport(imt)
		//	//r.Pkg = imt.Package
		//} else {
		//	//dd, _ := json.MarshalIndent(file.Imports, "", "\t")
		//	//fmt.Println(string(dd))
		//	panic(fmt.Errorf("can't found import package: %s in %s", pkgName, i.file.path))
		//}
	case *ast.SliceExpr:
		// ignore range
		//r.Type = &Type{Name: "[]" + i.parseParam(mod, pkg, file, names, e.X).Type.Name}
		r.Type = parseType(file, "", e.X)
	case *ast.ArrayType:
		// ignore len
		//r.Type = &Type{Name: "[]" + i.parseParam(mod, pkg, file, names, e.Elt).Type.Name}
		r.Type = parseType(file, "", e.Elt)
	case *ast.StarExpr:
		//r.Type = &Type{Name: "*" + i.parseParam(mod, pkg, file, names, e.X).Type.Name}
		r.Type = parseType(file, "", e.X)
	case *ast.FuncType:
		// TODO
	case *ast.MapType:
		r.Type = parseType(file, "", e.Value)
	case *ast.ChanType:
		// TODO
	default:
		panic(fmt.Errorf("unknown type: %s", reflect.TypeOf(e)))
	}
	
	return
}

func (i *Interface) addImport(imt *Import) {
	if i.imports == nil {
		i.imports = map[string]bool{}
	}
	if _, ok := i.imports[imt.Path]; ok {
		return
	}
	i.imports[imt.Path] = true
	i.Imports = append(i.Imports, imt)
}

func isReserved(t string) bool {
	switch t {
	case "int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"float32",
		"float64",
		"string",
		"bool",
		"byte",
		"rune",
		"uintptr",
		"any",
		"error":
		return true
	}
	return false
}
