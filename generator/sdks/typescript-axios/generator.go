package typescript_axios

import (
	"fmt"
	"github.com/gozelle/mix/generator/openapi"
	"github.com/gozelle/mix/generator/writter"
	"github.com/gozelle/openapi/openapi3"
	"github.com/gozelle/pongo2"
	"github.com/gozelle/spew"
	"github.com/gozelle/structs"
	"path/filepath"
	"strings"
)

type API struct {
	Name    string    `json:"Name,omitempty"`
	Methods []*Method `json:"Methods,omitempty"`
	Types   []*Type   `json:"Types,omitempty"`
}

type Method struct {
	Method   string `json:"Method,omitempty"`
	Path     string `json:"Path,omitempty"`
	Name     string `json:"Name,omitempty"`
	Request  string `json:"Request,omitempty"`
	Response string `json:"Response,omitempty"`
	Comment  string `json:"Comment,omitempty"`
}

type Type struct {
	Name   string  `json:"Name,omitempty"`
	Type   string  `json:"Type,omitempty"`
	Fields []*Type `json:"Fields,omitempty"`
	Elem   *Type   `json:"Elem,omitempty"`
	Ref    string  `json:"Ref,omitempty"`
}

func Generate(file string) (files []*writter.File, err error) {
	
	doc, err := openapi.Load(file)
	if err != nil {
		err = fmt.Errorf("load openapi file error: %s", err)
		return
	}
	
	api, err := Convert(doc)
	if err != nil {
		err = fmt.Errorf("convert ts type error: %s", err)
		return
	}
	
	fmt.Println("print API:")
	spew.Json(api)
	
	tpl, err := pongo2.FromString(typesTpl)
	if err != nil {
		err = fmt.Errorf("preprea tpl error: %s", err)
		return
	}
	c, err := tpl.Execute(structs.Map(api))
	files = append(files, []*writter.File{
		{Name: "api.ts", Content: c},
		{Name: "base.ts", Content: baseTpl},
	}...)
	
	return
}

func Convert(doc *openapi.DocumentV3) (api *API, err error) {
	
	api = &API{
		Name: "SDK",
	}
	
	for path, item := range doc.Paths {
		var m *Method
		m, err = convertMethod(doc, path, item)
		if err != nil {
			return
		}
		api.Methods = append(api.Methods, m)
	}
	
	if doc.Components != nil {
		if doc.Components.Schemas != nil {
			api.Types = append(api.Types, convertSchemas(doc, doc.Components.Schemas)...)
		}
		if doc.Components.RequestBodies != nil {
			api.Types = append(api.Types, convertRequestBodies(doc, doc.Components.RequestBodies)...)
		}
		if doc.Components.Responses != nil {
			api.Types = append(api.Types, convertResponses(doc, doc.Components.Responses)...)
		}
	}
	
	return
}

func convertMethod(doc *openapi.DocumentV3, path string, item *openapi3.PathItem) (m *Method, err error) {
	var op *openapi3.Operation
	var method string
	if item.Post != nil {
		op = item.Post
		method = "post"
	} else if item.Get != nil {
		op = item.Get
		method = "get"
	} else if item.Put != nil {
		op = item.Put
		method = "put"
	} else if item.Delete != nil {
		op = item.Delete
		method = "delete"
	} else {
		err = fmt.Errorf("handle path: %s error: only accept POST、GET、PUT、DELETE method", path)
		return
	}
	
	m = &Method{
		Method:   method,
		Path:     path,
		Name:     op.OperationID,
		Comment:  op.Description,
		Request:  "",
		Response: "",
	}
	
	if op.RequestBody != nil {
		m.Request = convertMethodRequestBody(op.RequestBody)
	}
	
	if op.Responses != nil {
		m.Response = convertMethodResponses(op.Responses)
	}
	
	return
}

func convertMethodRequestBody(req *openapi3.RequestBodyRef) (t string) {
	if req.Ref != "" {
		t = filepath.Base(req.Ref)
	}
	return
}

func convertMethodResponses(resp openapi3.Responses) (t string) {
	if resp["200"] != nil &&
		resp["200"].Value != nil &&
		resp["200"].Value.Content != nil &&
		resp["200"].Value.Content[openapi.ApplicationJson] != nil &&
		resp["200"].Value.Content[openapi.ApplicationJson].Schema != nil &&
		resp["200"].Value.Content[openapi.ApplicationJson].Schema.Ref != "" {
		t = filepath.Base(resp["200"].Value.Content[openapi.ApplicationJson].Schema.Ref)
	}
	return
}

func convertSchemas(doc *openapi.DocumentV3, schemas openapi3.Schemas) (types []*Type) {
	
	for name, v := range schemas {
		types = append(types, convertSchemaRef(doc, name, v))
	}
	
	return
}

func convertRequestBodies(doc *openapi.DocumentV3, reqs openapi3.RequestBodies) (types []*Type) {
	
	for name, v := range reqs {
		if v.Value != nil &&
			v.Value.Content != nil &&
			v.Value.Content[openapi.ApplicationJson] != nil &&
			v.Value.Content[openapi.ApplicationJson].Schema != nil {
			types = append(types, convertSchemaRef(doc, name, v.Value.Content[openapi.ApplicationJson].Schema))
		}
	}
	
	return
}

func convertResponses(doc *openapi.DocumentV3, reps openapi3.Responses) (types []*Type) {
	for name, v := range reps {
		if v.Value != nil &&
			v.Value.Content != nil &&
			v.Value.Content[openapi.ApplicationJson] != nil &&
			v.Value.Content[openapi.ApplicationJson].Schema != nil {
			types = append(types, convertSchemaRef(doc, name, v.Value.Content[openapi.ApplicationJson].Schema))
		}
	}
	return
}

func convertSchemaRef(doc *openapi.DocumentV3, name string, schema *openapi3.SchemaRef) (t *Type) {
	if schema.Ref != "" {
		if strings.HasPrefix(schema.Ref, "#/components/schemas") {
			if v := getSchemaRef(doc, filepath.Base(schema.Ref)); v != nil {
				schema.Value = v.Value
			}
			
		}
	}
	t = convertSchema(doc, name, schema.Value)
	return
}

func convertSchema(doc *openapi.DocumentV3, name string, value *openapi3.Schema) (field *Type) {
	field = &Type{
		Name: name,
	}
	switch value.Type {
	case openapi.Object:
		if value.Properties != nil {
			field.Type = openapi.Object
			for k, v := range value.Properties {
				if v.Ref != "" {
					field.Ref = filepath.Base(v.Ref)
				} else if v.Value != nil {
					field.Fields = append(field.Fields, convertSchema(doc, k, v.Value))
				}
			}
		} else {
			field.Type = "any"
		}
	case openapi.Array:
		typ := "any"
		if value.Items != nil {
			v := value.Items.Value
			if value.Items.Ref != "" {
				typ = filepath.Base(value.Items.Ref)
			} else if v != nil {
				typ = convertSchema(doc, "", v).Type
			}
		}
		field.Type = fmt.Sprintf("%s[]", typ)
	case openapi.String:
		field.Type = openapi.String
	case openapi.Boolean:
		field.Type = openapi.Boolean
	case openapi.Integer, openapi.Number:
		field.Type = openapi.Number
	}
	return
}

func getSchemaRef(doc *openapi.DocumentV3, name string) *openapi3.SchemaRef {
	if doc.Components != nil &&
		doc.Components.Schemas != nil {
		ref := doc.Components.Schemas[name]
		if ref != nil {
			if ref.Ref != "" {
				return getSchemaRef(doc, filepath.Base(ref.Ref)) // TODO 递归获取引用，可以加上计数避免死循环
			}
			return ref
		}
	}
	return nil
}

const typesTpl = `
import {BaseAPI} from "./base";
import {AxiosRequestConfig} from "axios";

export class API extends BaseAPI {
{%- for method in Methods %}
    public {{method.Name}} ({% if method.Request %}request?: {{method.Request}}, {% endif %}options?: AxiosRequestConfig):Promise<{% if method.Response %}{{method.Response}}{% else %}null{% endif %}> {
        return this.client.{{ method.Method }}('{{ method.Path }}',{% if method.Request %}request{% else %}null{% endif %}, options)
    }
{%- endfor %}
}

{%- for type in Types %}
export interface {{ type.Name }} {
	{%- for field in type.Fields %}
    {{ field.Name }}?: {{ field.Type }};
	{%- endfor %}
}
{%- endfor %}
`
