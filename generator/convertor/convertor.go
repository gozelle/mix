package convertor

import (
	"encoding/json"
	"fmt"
	"github.com/gozelle/logging"
	"github.com/gozelle/mix/generator/langs/golang"
	"github.com/gozelle/mix/generator/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

var log = logging.Logger("convertor")

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
	
	for _, v := range i.Defs {
		if !v.Used || v.Type.Name != "struct" {
			continue
		}
		d := convertDef(v)
		r.Defs = append(r.Defs, d)
		fmt.Printf("===============Start Def: %s  ===============\n", v.Name)
		dd, _ := json.MarshalIndent(d, "", "\t")
		fmt.Println(string(dd))
		fmt.Printf("===============End Def: %s ===============\n", v.Name)
	}
	
	for _, v := range i.Methods {
		r.Methods = append(r.Methods, parseRenderMethod(v))
	}
	
	return r
}

func convertDef(d *parser.Def) *golang.Def {
	
	rt := d.Type.RealType()
	
	n := &golang.Def{
		Name: d.Name,
		Type: rt.Name,
	}
	
	if n.Type == "struct" {
		for _, v := range rt.StructFields {
			n.StructFields = append(n.StructFields, convertType(v.Field, v))
		}
	} else if n.Type == "[]" {
		n.Elem = convertType(rt.Elem.Field, rt.Elem)
	}
	
	return n
}

func convertType(name string, t *parser.Type) *golang.Def {
	n := &golang.Def{
		Name: name,
		Type: t.RealType().Name,
		Tags: t.Tags,
	}
	if n.Type == "[]" {
		n.Elem = convertType(t.Elem.Field, t.Elem)
	}
	
	return n
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
