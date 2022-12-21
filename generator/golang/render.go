package golang

import (
	"fmt"
	"github.com/gozelle/mix/parser"
	"strings"
)

type RenderInterface struct {
	Package  string
	Name     string
	Methods  []*RenderMethod
	Types    []*RenderType
	Packages []*RenderPackage
}

type RenderMethod struct {
	Name    string
	Params  string
	Results string
}

type RenderType struct {
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
	
	for _, v := range i.Types {
		r.Types = append(r.Types, parseRenderType(v))
	}
	
	for _, v := range i.Methods {
		r.Methods = append(r.Methods, parseRenderMethod(v))
	}
	
	return r
}

func parseRenderMethod(m *parser.Method) *RenderMethod {
	var params []string
	for _, v := range m.Params {
		params = append(params, fmt.Sprintf("%s %s", strings.Join(v.Names, ","), v.Type))
	}
	
	var results []string
	for _, v := range m.Results {
		results = append(results, fmt.Sprintf("%s %s", strings.Join(v.Names, ","), v.Type))
	}
	
	r := &RenderMethod{
		Name:    m.Name,
		Params:  strings.Join(params, ","),
		Results: strings.Join(results, ","),
	}
	if len(results) > 0 {
		r.Results = fmt.Sprintf("(%s)", r.Results)
	}
	
	return r
}

func parseRenderType(t *parser.Type) *RenderType {
	
	r := &RenderType{
		Name: t.Name,
		Type: t.Type,
	}
	
	if t.Type == "struct" {
		
		var fields []*parser.Type
		if t.Pointer {
			fields = t.Elem.Fields
		} else {
			fields = t.Fields
		}
		for _, v := range fields {
			r.Fields = append(r.Fields, &RenderField{
				Name: v.Name,
				Type: v.Type,
				Tags: v.Tags,
			})
		}
	}
	
	return r
}
