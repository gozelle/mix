package generator

import (
	"github.com/gozelle/mix/generator/render"
)

var (
	GenGoFileSuffix = ".mix.go"
)

type File struct {
	Name    string
	Content string
}

type Generator interface {
	Generate(i *render.Interface) (files []*File, err error)
}
