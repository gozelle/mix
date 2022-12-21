package gen_golang

import (
	"fmt"
	"github.com/gozelle/mix/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

type RenderInterface struct {
	Package  string
	Name     string
	Methods  []*RenderMethod
	Defs     []*RenderDef
	Packages []*RenderPackage
}

type RenderMethod struct {
	Name    string
	Request *RenderDef
	Replay  *RenderDef
	Params  string
	Results string
}

type RenderDef struct {
	Name   string
	Type   string
	Fields []*RenderField
}

type RenderField struct {
	Name string
	Type string
	Tags string
}

type RenderPackage struct {
	Alias string
	Path  string
}

func PrepareRenderInterface(pkg string, i *parser.Interface) *RenderInterface {
	
	r := &RenderInterface{
		Package: pkg,
		Name:    i.Name,
	}
	
	for _, v := range i.Packages {
		r.Packages = append(r.Packages, &RenderPackage{
			Alias: v.Alias,
			Path:  v.Path,
		})
	}
	
	for _, v := range i.Defs {
		r.Defs = append(r.Defs, parseRenderType(v))
	}
	
	for _, v := range i.Methods {
		r.Methods = append(r.Methods, parseRenderMethod(v))
	}
	
	return r
}

func parseRenderMethod(m *parser.Method) *RenderMethod {
	
	request := &RenderDef{
		Name: fmt.Sprintf("%sRequest", m.Name),
	}
	replay := &RenderDef{
		Name: fmt.Sprintf("%sReplay", m.Name),
	}
	
	var params []string
	merge := true
	for _, v := range m.Params {
		if !v.Type.IsContext() && merge {
			if string(v.Type) == request.Name {
				request.Type = string(v.Type)
				merge = false
			} else {
				request.Fields = append(request.Fields, convertMethodParam(v)...)
			}
			
		}
		params = append(params, fmt.Sprintf("%s %s", strings.Join(v.Names, ","), v.Type))
	}
	
	var results []string
	merge = true
	for _, v := range m.Results {
		if !v.Type.IsError() && merge {
			if string(v.Type) == request.Name {
				replay.Type = string(v.Type)
				merge = false
			} else {
				replay.Fields = append(replay.Fields, convertMethodParam(v)...)
			}
		}
		results = append(results, fmt.Sprintf("%s %s", strings.Join(v.Names, ","), v.Type))
	}
	
	r := &RenderMethod{
		Name:    m.Name,
		Request: request,
		Replay:  replay,
		Params:  strings.Join(params, ","),
		Results: strings.Join(results, ","),
	}
	if len(results) > 0 {
		r.Results = fmt.Sprintf("(%s)", r.Results)
	}
	
	return r
}

func convertMethodParam(p *parser.Param) []*RenderField {
	r := make([]*RenderField, 0)
	for _, v := range p.Names {
		r = append(r, &RenderField{
			Name: Title(v),
			Type: string(p.Type),
		})
	}
	return r
}

func Title(v string) string {
	return cases.Title(language.English).String(v)
}

func parseRenderType(t *parser.Def) *RenderDef {
	
	r := &RenderDef{
		Name: t.Name,
		Type: string(t.Type),
	}
	
	if t.Type == "struct" {
		
		for _, v := range t.Fields {
			r.Fields = append(r.Fields, &RenderField{
				Name: v.Name,
				Type: string(v.Type),
				Tags: v.Tags,
			})
		}
	}
	
	return r
}
