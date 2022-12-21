package parser

import (
	"fmt"
	"go/ast"
	"reflect"
)

func parseNames(idents []*ast.Ident) []string {
	names := make([]string, 0)
	for _, i := range idents {
		names = append(names, i.Name)
	}
	return names
}

func parseType(name string, t ast.Expr) (r *Type) {
	
	// TODO Cache package.Type
	r = &Type{
		Name: name,
	}
	
	switch e := t.(type) {
	case *ast.Ident:
		r.Type = e.Name
		r.Reserved = true
	case *ast.StructType:
		r.Type = "struct"
		for _, f := range e.Fields.List {
			st := parseType(parseNames(f.Names)[0], f.Type)
			if f.Tag != nil {
				st.Tags = f.Tag.Value
			}
			r.Fields = append(r.Fields, st)
		}
	case *ast.SelectorExpr:
		// TODO 处理包引用类型
		r.Type = fmt.Sprintf("%s.%s", e.X.(*ast.Ident), e.Sel.Name)
	case *ast.StarExpr:
		r.Pointer = true
		r.Elem = parseType(name, e.X)
	default:
		// TODO 报错，未知的类型
		fmt.Println(e)
		
		panic(fmt.Errorf("unknown type: %s", reflect.TypeOf(e)))
	}
	
	return
}
