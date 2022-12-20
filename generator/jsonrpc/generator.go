package jsonrpc

import (
	"github.com/fatih/structs"
	"github.com/flosch/pongo2/v6"
	"github.com/gozelle/mix/generator/golang"
	"github.com/gozelle/mix/parser"
)

type Maker struct {
}

func (m Maker) Generate(i *parser.Interface) (files []*parser.GenFile, err error) {
	f, err := renderMethod(i)
	if err != nil {
		return
	}
	files = append(files, f)
	return
}

func renderMethod(i *parser.Interface) (file *parser.GenFile, err error) {
	tpl, err := pongo2.FromString(serviceTpl)
	if err != nil {
		panic(err)
	}
	d := golang.PrepareRenderInterface("", i)
	m := structs.Map(d)
	out, err := tpl.Execute(m)
	if err != nil {
		panic(err)
	}
	file = &parser.GenFile{
		Name:    "service.go",
		Content: out,
	}
	return
}

func renderType() {

}

const serviceTpl = `
type {{ Name }} struct{
{% for method in Methods %}
    {{ method.Name }} func({{ method.Params }}) {{method.Results}}
{% endfor %}
}
`
