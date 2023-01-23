package generator

import (
	"github.com/gozelle/mix/generator/parser"
)

type File struct {
	Name    string
	Content string
}

type Generator interface {
	Generate(i *parser.Interface) (files []*File, err error)
}
