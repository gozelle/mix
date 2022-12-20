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
	
	for _, v := range i.Methods {
		r.Methods = append(r.Methods, parseRenderMethod(v))
	}
	
	return r
}

func parseRenderMethod(m *parser.Method) *RenderMethod {
	var params []string
	for _, v := range m.Params {
		params = append(params, fmt.Sprintf("%s %s", v.Name, parseFiledType(v)))
	}
	
	var results []string
	for _, v := range m.Results {
		results = append(results, fmt.Sprintf("%s %s", v.Name, parseFiledType(v)))
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

func parseFiledType(f *parser.Field) string {
	if f.Type.Pointer {
		return fmt.Sprintf("*%s", f.Type.Type)
	}
	return f.Type.Type
}
