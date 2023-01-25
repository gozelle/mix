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

func (i *Interface) addType(t *Def) {
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
	
	switch e := t.(type) {
	case *ast.Ident:
		r.Type = Type{t: e.Name}
		if !isReserved(r.Type.Type()) {
			// 不是基础类型，寻找到该类型定义，放入 Interface 上下文中
			rt := pkg.getDef(r.Type.Type())
			if rt == nil {
				panic(fmt.Errorf("can't fond type: '%s' in package: %s", r.Type.Type(), i.file.path))
			} else {
				i.addType(rt)
			}
		}
	case *ast.SelectorExpr:
		pkgName := e.X.(*ast.Ident).String()
		typ := e.Sel.Name
		r.Type = Type{t: fmt.Sprintf("%s.%s", pkgName, typ)}
		// 去 Parser 上下文中寻找对应的 Package 进来
		imt := file.getImport(pkgName)
		if imt != nil {
			i.addImport(imt)
		} else {
			// TODO
		}
	case *ast.SliceExpr:
		// ignore range
		r.Type = Type{t: "[]" + i.parseParam(mod, pkg, file, names, e.X).Type.t}
	case *ast.ArrayType:
		// ignore len
		r.Type = Type{t: "[]" + i.parseParam(mod, pkg, file, names, e.Elt).Type.t}
	case *ast.StarExpr:
		r.Type = Type{t: "*" + i.parseParam(mod, pkg, file, names, e.X).Type.t}
	case *ast.FuncType:
		// TODO
	case *ast.MapType:
		// TODO
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
