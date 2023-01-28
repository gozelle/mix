package openapi

import (
	"fmt"
	"github.com/gozelle/fs"
	"github.com/gozelle/mix/generator/parser"
	"github.com/gozelle/openapi/openapi3"
	"github.com/gozelle/pointer"
	"github.com/invopop/yaml"
	"strings"
)

type GenerateCmd struct {
	Path      string
	Interface string
}

func Parse(tplFile, path string, api string) (doc *DocumentV3, err error) {
	
	doc, err = Load(tplFile)
	if err != nil {
		err = fmt.Errorf("load template error: %s", err)
		return
	}
	
	mod, err := parser.PrepareMod()
	if err != nil {
		err = fmt.Errorf("prepare mod error: %s", err)
		return
	}
	
	pkg, err := parser.Parse(mod, path)
	if err != nil {
		err = fmt.Errorf("parse package error: %s", err)
		return
	}
	
	i := pkg.GetInterface(api)
	if i == nil {
		err = fmt.Errorf("api interface: %s not found", api)
		return
	}
	
	r := ConvertAPI(i)
	
	ConvertOpenapi(doc, r)
	
	return
}

func Load(file string) (doc *DocumentV3, err error) {
	doc = &DocumentV3{}
	if file != "" {
		var c []byte
		c, err = fs.Read(file)
		if err != nil {
			err = fmt.Errorf("read openapi file error: %s", err)
			return
		}
		if strings.HasSuffix(file, ".json") {
			err = doc.UnmarshalJSON(c)
			if err != nil {
				err = fmt.Errorf("unmarshal openapi file error: %s", err)
				return
			}
		} else if strings.HasSuffix(file, ".yaml") {
			var j []byte
			j, err = yaml.YAMLToJSON(c)
			if err != nil {
				err = fmt.Errorf("convert yaml to json error: %s", err)
				return
			}
			
			err = doc.UnmarshalJSON(j)
			if err != nil {
				err = fmt.Errorf("unmarshal openapi file error: %s", err)
				return
			}
		} else {
			err = fmt.Errorf("unsupport openapi file: %s suffix, accept json or yaml", file)
			return
		}
	}
	
	return
}

func ConvertOpenapi(doc *DocumentV3, r *API) {
	
	doc.OpenAPI = "3.0.3"
	
	if doc.Info == nil {
		doc.Info = &openapi3.Info{}
	}
	
	for _, v := range r.Methods {
		convertOpenapiMethods(doc, v)
	}
	
	if doc.Components == nil {
		doc.Components = &openapi3.Components{}
	}
	if doc.Components.Schemas == nil {
		doc.Components.Schemas = map[string]*openapi3.SchemaRef{}
	}
	
	for _, v := range r.Defs {
		doc.Components.Schemas[v.Name] = makeOpenapiSchemaRef(doc, v)
	}
	
	return
}

func convertOpenapiMethods(d *DocumentV3, m *Method) {
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
		item.Post.RequestBody = &openapi3.RequestBodyRef{
			Ref: makeOpenapiMethodParameterRef(d, m.Request),
		}
	}
	if item.Post.Responses == nil {
		item.Post.Responses = map[string]*openapi3.ResponseRef{}
	}
	// 过滤响应 io.Reader 和 chan
	if m.Replay != nil && len(m.Replay.StructFields) > 0 &&
		m.Replay.StructFields[0].Type != parser.TAny &&
		m.Replay.StructFields[0].Type != parser.TChan {
		item.Post.Responses["200"] = &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: pointer.ToString("success"),
				Content: openapi3.Content{
					ApplicationJson: &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Ref: makeOpenapiMethodReplyRef(d, m.Replay),
						},
					},
				},
			},
		}
	} else if m.Replay.Use != nil {
		item.Post.Responses["200"] = &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: pointer.ToString("success"),
				Content: openapi3.Content{
					ApplicationJson: &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Ref: fmt.Sprintf("#/components/schemas/%s", m.Replay.Use.Name),
						},
					},
				},
			},
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

func makeOpenapiMethodParameterRef(d *DocumentV3, def *Def) (ref string) {
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
				Content:  makeOpenapiContent(d, def),
			},
		}
	} else {
		d.Components.RequestBodies[def.Field] = &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required: false,
				Content: openapi3.Content{
					ApplicationJson: &openapi3.MediaType{
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

func makeOpenapiMethodReplyRef(d *DocumentV3, def *Def) (ref string) {
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
				Content:     makeOpenapiContent(d, def),
			},
		}
	} else {
		d.Components.Responses[def.Field] = &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: pointer.ToString(""),
				Headers:     nil,
				Content: openapi3.Content{
					ApplicationJson: &openapi3.MediaType{
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

func makeOpenapiContent(d *DocumentV3, def *Def) openapi3.Content {
	var c openapi3.Content = map[string]*openapi3.MediaType{}
	c[ApplicationJson] = &openapi3.MediaType{
		Extensions: nil,
		Schema:     makeOpenapiSchemaRef(d, def),
	}
	return c
}

func makeOpenapiSchemaRef(d *DocumentV3, def *Def) (s *openapi3.SchemaRef) {
	
	s = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: convertOpenapiType(def.Type),
		},
	}
	
	if def.Type == parser.TStruct {
		s.Value.Properties = map[string]*openapi3.SchemaRef{}
		for _, v := range def.StructFields {
			name := v.Field
			if v.Json != "" {
				name = v.Json
			}
			s.Value.Properties[name] = makeOpenapiSchemaRef(d, v)
		}
	} else if def.Type == parser.TSlice {
		s.Value.Items = makeOpenapiSchemaRef(d, def.Elem)
	} else if def.Use != nil {
		s.Ref = fmt.Sprintf("#/components/schemas/%s", def.Use.Name)
	}
	
	return
}

func convertOpenapiType(t string) string {
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
	case parser.TStruct, parser.TMap:
		return "object"
	}
	return ""
}
