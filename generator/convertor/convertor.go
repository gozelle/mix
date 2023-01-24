package convertor

import (
	"fmt"
	"github.com/gozelle/mix/generator/langs/golang"
	"github.com/gozelle/mix/generator/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func ToGolangInterface(i *parser.Interface) *golang.Interface {
	
	r := &golang.Interface{
		//Package: pkg,
		Name: i.Name,
	}
	
	for _, v := range i.Packages {
		r.Packages = append(r.Packages, &golang.Package{
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

func parseRenderMethod(m *parser.Method) *golang.Method {
	
	request := &golang.Def{
		Name: fmt.Sprintf("%sRequest", m.Name),
	}
	replay := &golang.Def{
		Name: fmt.Sprintf("%sReplay", m.Name),
	}
	
	var params []string
	merge := true
	for _, v := range m.Params {
		if !v.Type.IsContext() && merge {
			if string(v.Type) == request.Name {
				request.Type = v.Type
				merge = false
			} else {
				request.Properties = append(request.Properties, convertMethodParam(v)...)
			}
			
		}
		params = append(params, fmt.Sprintf("%s %s", strings.Join(v.Names, ","), v.Type))
	}
	
	var results []string
	merge = true
	for _, v := range m.Results {
		if !v.Type.IsError() && merge {
			if string(v.Type) == request.Name {
				replay.Type = v.Type
				merge = false
			} else {
				replay.Properties = append(replay.Properties, convertMethodParam(v)...)
			}
		}
		results = append(results, fmt.Sprintf("%s %s", strings.Join(v.Names, ","), v.Type))
	}
	
	r := &golang.Method{
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

func convertMethodParam(p *parser.Param) []*golang.Def {
	r := make([]*golang.Def, 0)
	for _, v := range p.Names {
		r = append(r, &golang.Def{
			Name: Title(v),
			Type: p.Type,
		})
	}
	return r
}

func Title(v string) string {
	return cases.Title(language.English).String(v)
}

func parseRenderType(t *parser.Def) *golang.Def {
	
	r := &golang.Def{
		Name: t.Name,
		Type: t.Type,
	}
	
	if t.Type == "struct" {
		
		for _, v := range t.Fields {
			r.Properties = append(r.Properties, &golang.Def{
				Name: v.Name,
				Type: v.Type,
				Tags: v.Tags,
			})
		}
	}
	
	return r
}
