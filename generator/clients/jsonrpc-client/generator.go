package jsonrpc_client

import (
	"github.com/fatih/structs"
	"github.com/gozelle/mix/generator"
	"github.com/gozelle/mix/generator/parser"
	"github.com/gozelle/mix/generator/render"
	"github.com/gozelle/pongo2"
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
	t, err := pongo2.FromString(tpl)
	if err != nil {
		panic(err)
	}
	d := render.Convert(i)
	d.Package = "rpc"
	m := structs.Map(d)
	out, err := t.Execute(m)
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

const tpl = `
package {{ Package }}

import (
{%- for pkg in Imports %}
	{{ pkg.Alias }} "{{ pkg.Path }}"
{%- endfor %}
)


type {{ Name }} struct{
{%- for method in Methods %}
    {{ method.Name }} func({{ method.Params }}) {{method.Results}}
{%- endfor %}
}

type {{ Name }}API struct{
{%- for method in Methods %}
    {{ method.Name }} func(ctx context.Context,request *{{ method.Request.Name }}) (replay *{{method.Replay.Name}}, err error)
{%- endfor %}
}

{%- for def in Methods %}

type {{ def.Request.Name }} struct {
	{%- for field in def.Request.StructFields %}
	{{ field.Name }} {{ field.Type }}
	{%- endfor %}
}

type {{ def.Replay.Name }} struct {
	{%- for field in def.Replay.StructFields %}
	{{ field.Name }} {{ field.Type }}
	{%- endfor %}
}

{%- endfor %}
`
