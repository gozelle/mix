package parser

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

type Type struct {
	Type         string            `json:"Type,omitempty"`
	ToString     bool              `json:"ToString,omitempty"` // own String() method
	Pointer      bool              `json:"Pointer,omitempty"`
	StructFields []*Type           `json:"StructFields,omitempty"`
	Elem         *Type             `json:"Elem,omitempty"`
	Tags         reflect.StructTag `json:"Tags,omitempty"`
	Def          *Def              `json:"Def,omitempty"`  // 所使用的的类型声明
	Real         *Type             `json:"Real,omitempty"` // 所使用的真实类型
	Field        string            `json:"Field,omitempty"`
	pkg          *Package
}

func (t Type) Fork() *Type {
	n := &Type{
		Type:         t.Type,
		ToString:     t.ToString,
		Pointer:      t.Pointer,
		StructFields: nil,
		Elem:         nil,
		Tags:         t.Tags,
		Def:          t.Def,
		Real:         nil,
		Field:        t.Field,
		pkg:          t.pkg,
	}
	for _, v := range t.StructFields {
		n.StructFields = append(n.StructFields, v.Fork())
	}
	if t.Elem != nil {
		n.Elem = t.Elem.Fork()
	}
	if t.Real != nil {
		n.Real = t.Real.Fork()
	}
	return n
}

func (t *Type) NoPointer() *Type {
	if t.Pointer {
		return t.Real.NoPointer()
	}
	return t
}

func (t *Type) RealType() *Type {
	if t.Real == nil {
		return t
	}
	return t.Real.RealType()
}

func (t Type) String() string {
	d, _ := json.Marshal(t)
	return string(d)
}

func (t Type) IsString() bool {
	return t.Type == "string" || t.ToString
}

func (t Type) IsStruct() bool {
	return t.Type == "struct"
}

func (t Type) IsArray() bool {
	return strings.HasPrefix(t.Type, "[]")
}

func (t Type) IsContext() bool {
	//TODO
	return t.Type == "context.Context"
}

func (t Type) IsError() bool {
	return t.Type == "error"
}

func (t Type) Json() string {
	items := strings.Split(t.Tags.Get("json"), ",")
	if len(items) > 0 {
		return items[0]
	}
	return ""
}

func handleTypeDef(pkg *Package, field string, r *Type, name string) *Def {
	def := pkg.GetDef(name)
	if def == nil {
		return nil
	}
	def.Used = true
	if def.Type == nil {
		def.Type = parseType(def.File, "", def.Expr)
	}
	if def.ToString {
		r.Real = &Type{Type: TString, Field: field}
	} else {
		r.Def = def
		r.Real = def.Type
	}
	return def
}

func parseType(f *File, field string, t ast.Expr) (r *Type) {
	
	r = &Type{Field: field}
	switch e := t.(type) {
	case *ast.Ident:
		r.Type = e.Name
		if !isReserved(r.Type) {
			def := handleTypeDef(f.pkg, field, r, r.Type)
			if def == nil {
				panic(fmt.Errorf("can't fond type: '%s' in package: %s", r.Type, f.path))
			}
		}
	case *ast.InterfaceType:
		r.Type = TAny
	case *ast.StructType:
		r.Type = TStruct
		for i, fd := range e.Fields.List {
			var fn string
			if len(fd.Names) == 0 {
				fn = fmt.Sprintf("field%d", i)
			} else {
				fn = fd.Names[0].Name
			}
			st := parseType(f, fn, fd.Type)
			if fd.Tag != nil {
				st.Tags = reflect.StructTag(strings.Trim(fd.Tag.Value, "`"))
			}
			r.StructFields = append(r.StructFields, st)
		}
	case *ast.SliceExpr:
		// ignore range
		r.Type = TSlice
		r.Elem = parseType(f, "", e.X)
	case *ast.ArrayType:
		// ignore len
		r.Type = TArray
		r.Elem = parseType(f, "", e.Elt)
	case *ast.SelectorExpr:
		pkgName := e.X.(*ast.Ident).String()
		typeName := e.Sel.Name
		r.Type = fmt.Sprintf("%s.%s", pkgName, typeName)
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
		
		def := handleTypeDef(imt.Package, field, r, typeName)
		if def == nil {
			panic(fmt.Errorf("package: %s type %s def is nil in: %s", pkgName, typeName, f.path))
		}
		if !def.ToString {
			f.pkg.AddExternalNalDef(def)
		}
	case *ast.StarExpr:
		r.Pointer = true
		r.Real = parseType(f, "", e.X)
	case *ast.MapType:
		r.Type = TMap
	case *ast.FuncType:
		r.Type = TFunc
	case *ast.ChanType:
		r.Type = TChan
	default:
		panic(fmt.Errorf("unknown type: %s", reflect.TypeOf(e)))
	}
	
	return
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
		"map",
		"any",
		"error":
		return true
	}
	return false
}
