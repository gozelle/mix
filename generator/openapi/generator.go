package openapi

import (
	"fmt"
	"github.com/gozelle/mix/generator"
	"github.com/gozelle/mix/generator/convertor"
	"github.com/gozelle/mix/generator/langs/golang"
	"github.com/gozelle/mix/generator/parser"
	
	"github.com/gozelle/openapi/openapi3"
)

var _ generator.Generator = (*Generator)(nil)

type Generator struct {
}

func (g Generator) Generate(i *parser.Interface) (files []*generator.File, err error) {
	//TODO implement me
	panic("implement me")
}

func (g Generator) TOOpenapiV3(i *parser.Interface) *DocumentV3 {
	
	d := &DocumentV3{}
	
	r := convertor.ToGolangInterface(i)
	
	for _, v := range r.Methods {
		g.convertMethods(d, v)
	}
	
	for _, v := range r.Defs {
		g.makeDefSchema(d, v)
	}
	
	return d
}

func (g Generator) convertMethods(d *DocumentV3, m *golang.Method) {
	if d.Paths == nil {
		d.Paths = map[string]*openapi3.PathItem{}
	}
	item := &openapi3.PathItem{
		Summary:     "",
		Description: "",
		Post: &openapi3.Operation{
			Tags:        nil,
			Summary:     "",
			Description: "",
			OperationID: m.Name,
			Parameters:  nil,
			RequestBody: nil,
			Responses:   nil,
		},
	}
	if m.Request != nil {
		g.convertMethodParameter(d, m.Request)
	}
	if m.Replay != nil {
		g.convertMethodReply(d, m.Replay)
	}
	
	d.Paths[fmt.Sprintf("/%s", m.Name)] = item
}

func (g Generator) convertMethodParameter(d *DocumentV3, def *golang.Def) {
	if d.Components == nil {
		d.Components = &openapi3.Components{}
	}
	if d.Components.RequestBodies == nil {
		d.Components.RequestBodies = map[string]*openapi3.RequestBodyRef{}
	}
	d.Components.RequestBodies[def.Name] = &openapi3.RequestBodyRef{
		Value: &openapi3.RequestBody{
			Required: false,
			Content:  g.makeContent(d, def.Fields),
		},
	}
}

func (g Generator) convertMethodReply(d *DocumentV3, def *golang.Def) {
	if d.Components == nil {
		d.Components = &openapi3.Components{}
	}
	if d.Components.Responses == nil {
		d.Components.Responses = map[string]*openapi3.ResponseRef{}
	}
	d.Components.Responses[def.Name] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: nil,
			Headers:     nil,
			Content:     g.makeContent(d, def.Fields),
		},
	}
}

func (g Generator) makeDefSchema(d *DocumentV3, def *golang.Def) {

}

func (g Generator) makeFieldSchema(d *DocumentV3, field *golang.Field) (ref string) {
	
	return
}

func (g Generator) makeContent(d *DocumentV3, fields []*golang.Field) openapi3.Content {
	var c openapi3.Content = map[string]*openapi3.MediaType{}
	
	for _, filed := range fields {
		if filed.Type.IsStruct() {
			c[filed.Name] = &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{
					Ref:   g.makeFieldSchema(d, filed),
					Value: &openapi3.Schema{},
				},
			}
		} else {
			c[filed.Name] = g.makeMediaType(filed)
		}
	}
	
	return c
}

func (g Generator) makeMediaType(field *golang.Field) *openapi3.MediaType {
	return &openapi3.MediaType{
		Schema: &openapi3.SchemaRef{
			Ref: "",
			Value: &openapi3.Schema{
				Type:        g.convertType(field.Type),
				Title:       "",
				Format:      "",
				Description: "",
			},
		},
		Example:  nil,
		Examples: nil,
		Encoding: nil,
	}
}

func (g Generator) convertType(t parser.Type) string {
	
	return "any"
}
