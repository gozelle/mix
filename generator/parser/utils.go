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

func parseType(name string, t ast.Expr) (r *Def) {
	
	r = &Def{Name: name}
	
	switch e := t.(type) {
	case *ast.Ident:
		r.Type = Type{t: e.Name}
	case *ast.StructType:
		r.Type = Type{t: "struct"}
		for _, f := range e.Fields.List {
			st := parseType(parseNames(f.Names)[0], f.Type)
			if f.Tag != nil {
				st.Tags = f.Tag.Value
			}
			r.StructFields = append(r.StructFields, st)
		}
	case *ast.SliceExpr:
		// ignore range
		r.Type = Type{t: "[]"}
		r.Elem = &Def{
			Type: parseType(name, e.X).Type,
		}
	case *ast.ArrayType:
		// ignore len
		r.Type = Type{t: "[]"}
		r.Elem = &Def{
			Type: parseType(name, e.Elt).Type,
		}
	case *ast.SelectorExpr:
		r.Type = Type{t: fmt.Sprintf("%s.%s", e.X.(*ast.Ident), e.Sel.Name)}
	case *ast.StarExpr:
		r.Type = Type{t: fmt.Sprintf("*%s", parseType(name, e.X).Type)}
	default:
		panic(fmt.Errorf("unknown type: %s", reflect.TypeOf(e)))
	}
	
	return
}
