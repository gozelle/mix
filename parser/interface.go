package parser

import (
	"fmt"
	"go/ast"
	"reflect"
)

type Interface struct {
	Name     string
	Methods  []*Method
	Types    []*Type
	types    map[string]bool
	Packages []*Package
	packages map[string]bool
	includes []*Interface
	t        *ast.InterfaceType
}

func (i *Interface) load(p *Parser) (err error) {
	for _, m := range i.t.Methods.List {
		switch mt := m.Type.(type) {
		case *ast.Ident:
		// TODO parse include
		case *ast.FuncType:
			i.Methods = append(i.Methods, i.parseMethod(p, m.Names[0].Name, mt))
		}
	}
	return
}

func (i *Interface) parseMethod(p *Parser, name string, t *ast.FuncType) (r *Method) {
	r = &Method{Name: name}
	
	for _, f := range t.Params.List {
		r.Params = append(r.Params, i.parseParam(p, parseNames(f.Names), f.Type))
	}
	
	for _, f := range t.Results.List {
		r.Results = append(r.Results, i.parseParam(p, parseNames(f.Names), f.Type))
	}
	
	return
}

func (i *Interface) addType(t *Type) {
	if i.types == nil {
		i.types = map[string]bool{}
	}
	if _, ok := i.types[t.Name]; ok {
		return
	}
	i.types[t.Name] = true
	i.Types = append(i.Types, t)
}

func (i *Interface) parseParam(p *Parser, names []string, t ast.Expr) (r *Param) {
	
	r = &Param{
		Names: names,
	}
	
	switch e := t.(type) {
	case *ast.Ident:
		r.Type = e.Name
		// 不是基础数据类型，去当前包里，寻找到该类型定义，放入 Interface 上下文中
		if !isReserved(r.Type) {
			rt := p.findRootType(r.Type)
			if rt == nil {
				panic(fmt.Errorf("can't fond type: '%s' in root package", r.Type))
			} else {
				i.addType(rt)
			}
		}
		// TODO 进一步判读是基础类型还是复合类型
	case *ast.SelectorExpr:
		// TODO 处理包引用类型
		r.Type = fmt.Sprintf("%s.%s", e.X.(*ast.Ident), e.Sel.Name)
		i.addPackage(p.findRootImport(e.X.(*ast.Ident).String()))
		// 去 Parser 上下文中寻找对应的 Package 进来
	case *ast.SliceExpr:
		r.Type = "[]" + i.parseParam(p, names, e.X).Type
	case *ast.StarExpr:
		r.Type = "*" + i.parseParam(p, names, e.X).Type
	default:
		// TODO 报错，未知的类型
		panic(fmt.Errorf("unknown type: %s", reflect.TypeOf(e)))
	}
	
	return
}

func (i *Interface) addPackage(pkg *Package) {
	if i.packages == nil {
		i.packages = map[string]bool{}
	}
	if _, ok := i.packages[pkg.Path]; ok {
		return
	}
	i.packages[pkg.Path] = true
	i.Packages = append(i.Packages, pkg)
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
		"error":
		return true
	}
	return false
}
