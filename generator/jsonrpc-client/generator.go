package jsonrpc_client

import (
	"github.com/fatih/structs"
	"github.com/flosch/pongo2/v6"
	"github.com/gozelle/mix/generator"
	"github.com/gozelle/mix/generator/golang"
	"github.com/gozelle/mix/parser"
)

type Generator struct {
}

func (m Generator) Generate(i *parser.Interface) (files []*generator.File, err error) {
	f, err := renderMethod(i)
	if err != nil {
		return
	}
	files = append(files, f)
	return
}

func renderMethod(i *parser.Interface) (file *generator.File, err error) {
	tpl, err := pongo2.FromString(internalTpl)
	if err != nil {
		panic(err)
	}
	d := golang.PrepareRenderInterface("rpc", i)
	m := structs.Map(d)
	out, err := tpl.Execute(m)
	if err != nil {
		panic(err)
	}
	file = &generator.File{
		Name:    "service.go",
		Content: out,
	}
	return
}

func renderType() {

}

const internalTpl = `
package {{ Package }}

import (
{% for pkg in Packages %}
	{{ pkg.Alias }} "{{ pkg.Path }}"
{% endfor %}
)


type {{ Name }} struct{
{% for method in Methods %}
    {{ method.Name }} func({{ method.Params }}) {{method.Results}}
{% endfor %}
}

{% for type in Types %}
	// {{ type.Name }}
	{% if type.Type == "struct" %}
		type {{ type.Name }} {{ type.Type }} {
			{% for filed in type.Fields %}
				{{ field.Name }} {{ field.Type }}
			{% endfor %}
		}
	{% else %}
		type {{ type.Name }} {{ type.Type }}
	{% endif %}
{% endfor %}
`

const importTpl = `

`
