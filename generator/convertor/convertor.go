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
	
	for _, v := range i.Imports {
		r.Packages = append(r.Packages, &golang.Package{
			Alias: v.Alias,
			Path:  v.Path,
		})
	}
	
	// TODO
	//for _, v := range i.Defs {
	//	r.Defs = append(r.Defs, v)
	//}
	
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
			if v.Type.Name == request.Name {
				request.Type = v.Type.Name // TODO
				merge = false
			} else {
				request.StructFields = append(request.StructFields, convertMethodParam(v)...)
			}
			
		}
		params = append(params, fmt.Sprintf("%s %s", strings.Join(v.Names, ","), v.Type))
	}
	
	var results []string
	merge = true
	for _, v := range m.Results {
		if !v.Type.IsError() && merge {
			if v.Type.Name == request.Name {
				replay.Type = v.Type.Name
				merge = false
			} else {
				replay.StructFields = append(replay.StructFields, convertMethodParam(v)...)
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

//func convertMethod

func convertMethodParam(p *parser.Param) []*golang.Def {
	r := make([]*golang.Def, 0)
	for _, v := range p.Names {
		r = append(r, &golang.Def{
			Name: Title(v),
			Type: p.Type.Name,
		})
	}
	return r
}

func Title(v string) string {
	return cases.Title(language.English).String(v)
}
