package generator

import (
	"github.com/gozelle/mix/generator/openapi"
)

var (
	GenGoFileSuffix = ".mix.go"
)

type File struct {
	Name    string
	Content string
}

type Generator interface {
	Generate(i *openapi.API) (files []*File, err error)
}
