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

func parseType(f *File, t ast.Expr) (r *Type) {
	
	switch e := t.(type) {
	case *ast.Ident:
		r = &Type{Name: e.Name}
		if !isReserved(r.Name) {
			//r.Type.Pkg = f.pkg
			log.Infof("填充自定义类型: %s", r.Name)
			def := f.pkg.getDef(r.Name)
			if def == nil {
				panic(fmt.Errorf("can't fond type: '%s' in package: %s", r.Name, f.path))
			}
			def.Used = true
		}
		// 不是基础类型，寻找到该类型定义，放入 Interface 上下文中
		//	rt := pkg.getDef(r.Type.Name)
		//	if rt == nil {
		//		panic(fmt.Errorf("can't fond type: '%s' in package: %s", r.Type.Name, i.file.path))
		//	} else {
		//		i.addDef(rt)
		//		r.Def = rt
		//	}
	case *ast.InterfaceType:
		r = &Type{Name: "any"}
	case *ast.StructType:
		r = &Type{Name: "struct"}
		
		for _, field := range e.Fields.List {
			if len(field.Names) == 0 {
				// TODO
				continue
			}
			st := parseType(f, field.Type)
			if field.Tag != nil {
				st.Tags = field.Tag.Value
			}
			
			r.StructFields = append(r.StructFields, st)
		}
	case *ast.SliceExpr:
		// ignore range
		r = &Type{Name: "[]"}
		r.Elem = parseType(f, e.X)
	case *ast.ArrayType:
		// ignore len
		r = &Type{Name: "[]"}
		r.Elem = parseType(f, e.Elt)
	case *ast.SelectorExpr:
		pkgName := e.X.(*ast.Ident).String()
		typeName := e.Sel.Name
		r = &Type{Name: fmt.Sprintf("%s.%s", pkgName, typeName)}
		imt := f.getImport(pkgName)
		if imt == nil {
			log.Infof("import file path: %s", f.path)
			for k, v := range f.Imports {
				log.Infof("name: %s, import path:%s", k, v.Path)
			}
			panic(fmt.Errorf("cant' get import: %s in: %s", pkgName, f.path))
		}
		if imt.Package == nil {
			panic(fmt.Errorf("import: %s Package is nil in: %s", pkgName, f.path))
		}
		
		def := imt.Package.getDef(typeName)
		if def == nil {
			panic(fmt.Errorf("package: %s type %s def is nil in: %s", pkgName, typeName, f.path))
		}
		def.Used = true
		if def.ToString {
			r.Def = &Type{Name: "string"}
		} else {
			// TODO 思考子包结构如何映射到 Package 中
			r.Def = parseType(def.File, def.Expr)
		}
	
	case *ast.StarExpr:
		r = &Type{Pointer: true, Elem: parseType(f, e.X)}
	case *ast.MapType:
		// TODO
	case *ast.FuncType:
		// TODO
	case *ast.ChanType:
		// TODO
	default:
		panic(fmt.Errorf("unknown type: %s", reflect.TypeOf(e)))
	}
	
	return
}
