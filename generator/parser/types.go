package parser

import (
	"encoding/json"
	"go/ast"
	"strings"
)

type Import struct {
	Alias   string
	Path    string
	Package *Package
}

type Method struct {
	Name    string
	Params  []*Param
	Results []*Param
}

// func(a string) => Param(a Type{Name: string})
// func(fil Fil)  => Param(fil Type{Name: Fil, Def: Type{Name: decimal.Decimal, Def;Type{Name: string}}}
type Type struct {
	Name         string  `json:"name,omitempty"`
	ToString     bool    `json:"toString,omitempty"` // own String() method
	Pointer      bool    `json:"pointer,omitempty"`
	StructFields []*Type `json:"structFields,omitempty"`
	Elem         *Type   `json:"elem,omitempty"`
	Tags         string  `json:"tags,omitempty"`
	Def          *Type   `json:"def,omitempty"`
	Field        string  `json:"field,omitempty"`
}

func (t Type) Fork() *Type {
	n := &Type{
		Name:         t.Name,
		ToString:     t.ToString,
		Pointer:      t.Pointer,
		StructFields: nil,
		Elem:         nil,
		Tags:         t.Tags,
		Def:          nil,
		Field:        t.Field,
	}
	for _, v := range t.StructFields {
		n.StructFields = append(n.StructFields, v.Fork())
	}
	if t.Elem != nil {
		n.Elem = t.Elem.Fork()
	}
	if t.Def != nil {
		n.Def = t.Def.Fork()
	}
	return n
}

func (t *Type) RealType() *Type {
	if t.Def == nil {
		return t
	}
	return t.Def.RealType()
}

func (t Type) String() string {
	d, _ := json.Marshal(t)
	return string(d)
}

func (t Type) IsString() bool {
	return t.Name == "string" || t.ToString
}

func (t Type) IsStruct() bool {
	return t.Name == "struct"
}

func (t Type) IsArray() bool {
	return strings.HasPrefix(t.Name, "[]")
}

func (t Type) IsContext() bool {
	//TODO
	return t.Name == "context.Context"
}

func (t Type) IsError() bool {
	return t.Name == "error"
}

type Param struct {
	Names []string
	Type  *Type
	Def   *Def
}

type Def struct {
	Name     string   `json:"name"`
	Used     bool     `json:"used,omitempty"`
	Expr     ast.Expr `json:"-"`
	File     *File    `json:"-"`
	Type     *Type    `json:"type,omitempty"`
	ToString bool     `json:"toString,omitempty"`
}

func (d Def) ShallowFork() *Def {
	n := &Def{
		Name:     d.Name,
		Used:     d.Used,
		Expr:     d.Expr,
		File:     d.File,
		Type:     d.Type,
		ToString: d.ToString,
	}
	return n
}

func (d Def) String() string {
	a, _ := json.Marshal(d)
	return string(a)
}
