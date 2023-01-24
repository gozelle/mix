package parser

import "strings"

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

type Type struct {
	pkg      string
	t        string
	toString bool // own String() method
}

func (t Type) Type() string {
	return t.t
}

func (t Type) IsString() bool {
	return t.t == "string" || t.toString
}

func (t Type) IsStruct() bool {
	return t.t == "struct"
}

func (t Type) IsArray() bool {
	return strings.HasPrefix(t.t, "[]")
}

func (t Type) IsContext() bool {
	//TODO
	return t.t == "context.Context"
}

func (t Type) IsError() bool {
	return t.t == "error"
}

type Param struct {
	Names []string
	Type  Type
}

type Def struct {
	Name         string
	Type         Type
	Tags         string
	StructFields []*Def
	Elem         *Def
}
