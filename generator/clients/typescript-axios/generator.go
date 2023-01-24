package typescript_axios

import (
	"github.com/gozelle/mix/generator"
	"github.com/gozelle/mix/generator/parser"
)

var _ generator.Generator = (*Generator)(nil)

type Generator struct {
}

func (g Generator) Generate(i *parser.Interface) (files []*generator.File, err error) {
	//TODO implement me
	panic("implement me")
}
