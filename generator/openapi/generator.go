package openapi

import (
	"fmt"
	"github.com/gozelle/logging"
	"github.com/gozelle/mix/generator"
	"github.com/gozelle/mix/generator/convertor"
	"github.com/gozelle/mix/generator/langs/golang"
	"github.com/gozelle/mix/generator/parser"
	"github.com/gozelle/openapi/openapi3"
	"github.com/gozelle/pointer"
	"github.com/gozelle/spew"
)

var log = logging.Logger("openapi")

var _ generator.Generator = (*Generator)(nil)

type Generator struct {
}

func (g Generator) Generate(i *parser.Interface) (files []*generator.File, err error) {
	//TODO implement me
	panic("implement me")
}

func (g Generator) TOOpenapiV3(i *parser.Interface) *DocumentV3 {
	
	d := &DocumentV3{}
	d.OpenAPI = "3.0.3"
	d.Info = &openapi3.Info{
		Title:          "",
		Description:    "",
		TermsOfService: "",
		Contact:        nil,
		License:        nil,
		Version:        "",
	}
	
	r := convertor.ToGolangInterface(i)
	
	for _, v := range r.Methods {
		g.convertMethods(d, v)
	}
	
	if d.Components == nil {
		d.Components = &openapi3.Components{}
	}
	if d.Components.Schemas == nil {
		d.Components.Schemas = map[string]*openapi3.SchemaRef{}
	}
	for _, v := range r.Defs {
		d.Components.Schemas[v.Name] = g.makeSchemaRef(d, v)
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
		log.Infof("Request: %s", m.Request.Name)
		if m.Request.Name == "IntStructRequest" {
			spew.Json(m.Request)
		}
		item.Post.RequestBody = &openapi3.RequestBodyRef{
			Ref: g.makeMethodParameterRef(d, m.Request),
			//Value: &openapi3.RequestBody{
			//	Extensions:  nil,
			//	Description: "",
			//	Required:    false,
			//	Content: map[string]*openapi3.MediaType{
			//		application_json: {
			//			Schema: &openapi3.SchemaRef{
			//				Ref:,
			//			},
			//		},
			//	},
			//},
		}
	}
	if item.Post.Responses == nil {
		item.Post.Responses = map[string]*openapi3.ResponseRef{}
	}
	if m.Replay != nil {
		item.Post.Responses["200"] = &openapi3.ResponseRef{
			Ref: g.makeMethodReplyRef(d, m.Replay),
		}
	} else {
		item.Post.Responses["200"] = &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: pointer.ToString("success"),
			},
		}
	}
	
	d.Paths[fmt.Sprintf("/%s", m.Name)] = item
}

func (g Generator) makeMethodParameterRef(d *DocumentV3, def *golang.Def) (ref string) {
	if d.Components == nil {
		d.Components = &openapi3.Components{}
	}
	if d.Components.RequestBodies == nil {
		d.Components.RequestBodies = map[string]*openapi3.RequestBodyRef{}
	}
	
	if def.Concat {
		d.Components.RequestBodies[def.Name] = &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required: false,
				Content:  g.makeContent(d, def),
			},
		}
	} else {
		d.Components.RequestBodies[def.Name] = &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required: false,
				Content: openapi3.Content{
					application_json: &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Ref: fmt.Sprintf("#/components/schemas/%s", def.Name),
						},
					},
				},
			},
		}
	}
	ref = fmt.Sprintf("#/components/requestBodies/%s", def.Name)
	return
}

func (g Generator) makeMethodReplyRef(d *DocumentV3, def *golang.Def) (ref string) {
	if d.Components == nil {
		d.Components = &openapi3.Components{}
	}
	if d.Components.Responses == nil {
		d.Components.Responses = map[string]*openapi3.ResponseRef{}
	}
	if def.Concat {
		d.Components.Responses[def.Name] = &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: pointer.ToString(""),
				Headers:     nil,
				Content:     g.makeContent(d, def),
			},
		}
	} else {
		d.Components.Responses[def.Name] = &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: pointer.ToString(""),
				Headers:     nil,
				Content: openapi3.Content{
					application_json: &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Ref: fmt.Sprintf("#/components/schemas/%s", def.Name),
						},
					},
				},
			},
		}
	}
	ref = fmt.Sprintf("#/components/responses/%s", def.Name)
	
	return
}

func (g Generator) makeContent(d *DocumentV3, def *golang.Def) openapi3.Content {
	
	log.Infof("content golang def: %s", def.Name)
	spew.Json(def)
	
	var c openapi3.Content = map[string]*openapi3.MediaType{}
	c[application_json] = &openapi3.MediaType{
		Extensions: nil,
		Schema:     g.makeSchemaRef(d, def),
	}
	return c
}

func (g Generator) makeSchemaRef(d *DocumentV3, def *golang.Def) (s *openapi3.SchemaRef) {
	
	if def == nil { // TODO should remove
		return
	}
	
	s = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: g.convertType(def.Type),
		},
	}
	
	if def.Type == "struct" {
		s.Value.Properties = map[string]*openapi3.SchemaRef{}
		for _, v := range def.StructFields {
			name := v.Name
			if v.Json != "" {
				name = v.Json
			}
			s.Value.Properties[name] = g.makeSchemaRef(d, v)
		}
	} else if def.Type == "[]" {
		s.Value.Items = g.makeSchemaRef(d, def.Elem)
	} else if def.Type != "map" && !def.Reserved {
		s.Ref = fmt.Sprintf("#/components/schemas/%s", def.Type)
	}
	
	return
}

//func (g Generator) makeContent(d *DocumentV3, fields []*golang.Field) openapi3.Content {
//	var c openapi3.Content = map[string]*openapi3.MediaType{}
//
//	for _, filed := range fields {
//		if filed.Type.IsStruct() {
//			c[filed.Name] = &openapi3.MediaType{
//				Schema: &openapi3.SchemaRef{
//					Ref:   g.makeFieldSchema(d, filed),
//					Value: &openapi3.Schema{},
//				},
//			}
//		} else {
//			c[filed.Name] = g.makeMediaType(filed)
//		}
//	}
//
//	return c
//}

//func (g Generator) makeMediaType(field *golang.Field) *openapi3.MediaType {
//	return &openapi3.MediaType{
//		Schema: &openapi3.SchemaRef{
//			Ref: "",
//			Value: &openapi3.Schema{
//				Type:        g.convertType(field.Type),
//				Title:       "",
//				Format:      "",
//				Description: "",
//			},
//		},
//		Example:  nil,
//		Examples: nil,
//		Encoding: nil,
//	}
//}

func (g Generator) convertType(t string) string {
	
	if t == "[]" {
		return "array"
	}
	
	switch t {
	case golang.String:
		return "string"
	case golang.Int, golang.Int8, golang.Int16, golang.Int32, golang.Int64,
		golang.Uint, golang.Uint8, golang.Uint16, golang.Uint32, golang.Uint64:
		return "integer"
	case golang.Float32, golang.Float64:
		return "number"
	case golang.Bool:
		return "boolean"
	}
	
	return "object"
}
