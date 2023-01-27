package typescript_axios

import (
	"github.com/gozelle/mix/generator"
	"github.com/gozelle/mix/generator/render"
)

var _ generator.Generator = (*Generator)(nil)

type Generator struct {
}

func (g Generator) Generate(i *render.Interface) (files []*generator.File, err error) {
	//TODO implement me
	panic("implement me")
}

const tpl = `

`
