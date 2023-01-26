package parser

import (
	"encoding/json"
	"go/ast"
)

type Def struct {
	Name     string   `json:"Name"`
	Used     bool     `json:"Used,omitempty"`
	Expr     ast.Expr `json:"-"`
	File     *File    `json:"-"`
	Type     *Type    `json:"Type,omitempty"`
	ToString bool     `json:"ToString,omitempty"`
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
