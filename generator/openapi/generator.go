package openapi

import (
	"fmt"
	"github.com/gozelle/logging"
	"github.com/gozelle/mix/generator"
	"github.com/gozelle/mix/generator/parser"
	"github.com/gozelle/mix/generator/render"
	"github.com/gozelle/openapi/openapi3"
	"github.com/gozelle/pointer"
)

var log = logging.Logger("openapi")

var _ generator.Generator = (*Generator)(nil)

type Generator struct {
}

func (g Generator) Generate(r *render.Interface) (files []*generator.File, err error) {
	//TODO implement me
	panic("implement me")
}

func (g Generator) TOOpenapiV3(r *render.Interface) *DocumentV3 {
	
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

func (g Generator) convertMethods(d *DocumentV3, m *render.Method) {
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
		//log.Infof("Request: %s", m.Request.Field)
		//if m.Request.Field == "IntStructRequest" {
		//	spew.Json(m.Request)
		//}
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

func (g Generator) makeMethodParameterRef(d *DocumentV3, def *render.Def) (ref string) {
	if d.Components == nil {
		d.Components = &openapi3.Components{}
	}
	if d.Components.RequestBodies == nil {
		d.Components.RequestBodies = map[string]*openapi3.RequestBodyRef{}
	}
	
	if def.Use == nil {
		d.Components.RequestBodies[def.Field] = &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required: false,
				Content:  g.makeContent(d, def),
			},
		}
	} else {
		d.Components.RequestBodies[def.Field] = &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required: false,
				Content: openapi3.Content{
					application_json: &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Ref: fmt.Sprintf("#/components/schemas/%s", def.Use.Name),
						},
					},
				},
			},
		}
	}
	ref = fmt.Sprintf("#/components/requestBodies/%s", def.Field)
	return
}

func (g Generator) makeMethodReplyRef(d *DocumentV3, def *render.Def) (ref string) {
	if d.Components == nil {
		d.Components = &openapi3.Components{}
	}
	if d.Components.Responses == nil {
		d.Components.Responses = map[string]*openapi3.ResponseRef{}
	}
	if def.Use == nil {
		d.Components.Responses[def.Field] = &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: pointer.ToString(""),
				Headers:     nil,
				Content:     g.makeContent(d, def),
			},
		}
	} else {
		d.Components.Responses[def.Field] = &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: pointer.ToString(""),
				Headers:     nil,
				Content: openapi3.Content{
					application_json: &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Ref: fmt.Sprintf("#/components/schemas/%s", def.Use.Name),
						},
					},
				},
			},
		}
	}
	ref = fmt.Sprintf("#/components/responses/%s", def.Field)
	
	return
}

func (g Generator) makeContent(d *DocumentV3, def *render.Def) openapi3.Content {
	
	//log.Infof("content golang def: %s", def.Field)
	//spew.Json(def)
	
	var c openapi3.Content = map[string]*openapi3.MediaType{}
	c[application_json] = &openapi3.MediaType{
		Extensions: nil,
		Schema:     g.makeSchemaRef(d, def),
	}
	return c
}

func (g Generator) makeSchemaRef(d *DocumentV3, def *render.Def) (s *openapi3.SchemaRef) {
	
	s = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: g.convertType(def.Type),
		},
	}
	
	if def.Type == parser.TStruct {
		s.Value.Properties = map[string]*openapi3.SchemaRef{}
		for _, v := range def.StructFields {
			name := v.Field
			if v.Json != "" {
				name = v.Json
			}
			s.Value.Properties[name] = g.makeSchemaRef(d, v)
		}
	} else if def.Type == parser.TSlice {
		s.Value.Items = g.makeSchemaRef(d, def.Elem)
	} else if def.Use != nil {
		s.Ref = fmt.Sprintf("#/components/schemas/%s", def.Use.Name)
	}
	
	return
}

func (g Generator) convertType(t string) string {
	switch t {
	case parser.TSlice, parser.TArray:
		return "array"
	case parser.TString:
		return "string"
	case parser.TInt, parser.TInt8, parser.TInt16, parser.TInt32, parser.TInt64,
		parser.TUint, parser.TUint8, parser.TUint16, parser.TUint32, parser.TUint64:
		return "integer"
	case parser.TFloat32, parser.TFloat64:
		return "number"
	case parser.TBool:
		return "boolean"
	}
	return "object"
}
