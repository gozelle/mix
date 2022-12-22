package gen_jsonrpc_client

import (
	"github.com/fatih/structs"
	"github.com/flosch/pongo2/v6"
	"github.com/gozelle/mix/generator"
	gen_golang "github.com/gozelle/mix/generator/go"
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
	d := gen_golang.PrepareRenderInterface("rpc", i)
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

type {{ Name }}API struct{
{% for method in Methods %}
    {{ method.Name }} func(ctx context.Context,request *{{ method.Request.Name }}) (replay *{{method.Replay.Name}}, err error)
{% endfor %}
}

{% for def in Methods %}
	type {{ def.Request.Name }} struct {
		{% for field in def.Request.Fields %}
			 {{ field.Name }} {{ field.Type }}
		{% endfor %}
	}
	type {{ def.Replay.Name }} struct {
		{% for field in def.Replay.Fields %}
			 {{ field.Name }} {{ field.Type }}
		{% endfor %}
	}
{% endfor %}
`

const importTpl = `
{% for def in Defs %}
	{% if def.Type == "struct" %}
		type {{ def.Name }} struct{
			{% for field in def.Fields %}
				 {{ field.Name }} {{ field.Type }}
			{% endfor %}
		}
	{% else %}
		type {{ def.Name }} {{ def.Type }}
	{% endif %}
{% endfor %}
`
