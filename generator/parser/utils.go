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

func parseType(f *File, field string, t ast.Expr) (r *Type) {
	
	r = &Type{Field: field}
	switch e := t.(type) {
	case *ast.Ident:
		r.Name = e.Name
		if !isReserved(r.Name) {
			log.Infof("填充自定义类型: %s", r.Name)
			def := f.pkg.getDef(r.Name)
			if def == nil {
				panic(fmt.Errorf("can't fond type: '%s' in package: %s", r.Name, f.path))
			}
			def.Used = true
			r.Def = def.Type
		}
	case *ast.InterfaceType:
		r.Name = "any"
	case *ast.StructType:
		r.Name = "struct"
		for i, fd := range e.Fields.List {
			var fn string
			if len(fd.Names) == 0 {
				fn = fmt.Sprintf("field%d", i)
			} else {
				fn = fd.Names[0].Name
			}
			st := parseType(f, fn, fd.Type)
			if fd.Tag != nil {
				st.Tags = fd.Tag.Value
			}
			r.StructFields = append(r.StructFields, st)
		}
	case *ast.SliceExpr:
		// ignore range
		r.Name = "[]"
		r.Elem = parseType(f, "", e.X)
	case *ast.ArrayType:
		// ignore len
		r.Name = "[]"
		r.Elem = parseType(f, "", e.Elt)
	case *ast.SelectorExpr:
		pkgName := e.X.(*ast.Ident).String()
		typeName := e.Sel.Name
		r.Name = fmt.Sprintf("%s.%s", pkgName, typeName)
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
			r.Def = &Type{Name: "string", Field: field}
		} else {
			r.Def = parseType(def.File, field, def.Expr)
			f.pkg.AddExternalNalDef(def)
		}
	
	case *ast.StarExpr:
		r.Pointer = true
		r.Elem = parseType(f, "", e.X)
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
