package typescript_axios

import (
	"fmt"
	"github.com/gozelle/mix/generator"
	"github.com/gozelle/mix/generator/openapi"
	"github.com/gozelle/openapi/openapi3"
	"github.com/gozelle/spew"
	"path/filepath"
)

type API struct {
	Name    string
	Methods []*Method
	Types   []*Type
}

type Method struct {
	Method   string
	Path     string
	Name     string
	Request  string
	Response string
	Comment  string
}

type Type struct {
	Name   string
	Fields []*Filed
}

type Filed struct {
	Name string
	Type string
}

func Generate(file string) (files []*generator.File, err error) {
	
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
			api.Types = append(api.Types, convertSchemas(doc.Components.Schemas)...)
		}
		if doc.Components.RequestBodies != nil {
			api.Types = append(api.Types, convertRequestBodies(doc.Components.RequestBodies)...)
		}
		if doc.Components.Responses != nil {
			api.Types = append(api.Types, convertResponses(doc.Components.Responses)...)
		}
	}
	
	return
}

func convertMethod(doc *openapi.DocumentV3, path string, item *openapi3.PathItem) (m *Method, err error) {
	var op *openapi3.Operation
	var method string
	if item.Post != nil {
		op = item.Post
		method = "POST"
	} else if item.Get != nil {
		op = item.Get
		method = "GET"
	} else if item.Put != nil {
		op = item.Put
		method = "PUT"
	} else if item.Delete != nil {
		op = item.Delete
		method = "DELETE"
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

func convertSchemas(schemas openapi3.Schemas) (types []*Type) {
	
	for name, v := range schemas {
		types = append(types, convertSchema(name, v))
	}
	
	return
}

func convertRequestBodies(reqs openapi3.RequestBodies) (types []*Type) {
	
	for name, v := range reqs {
		if v.Value != nil &&
			v.Value.Content != nil &&
			v.Value.Content[openapi.ApplicationJson] != nil &&
			v.Value.Content[openapi.ApplicationJson].Schema != nil {
			types = append(types, convertSchema(name, v.Value.Content[openapi.ApplicationJson].Schema))
		}
	}
	
	return
}

func convertResponses(reps openapi3.Responses) (types []*Type) {
	for name, v := range reps {
		if v.Value != nil &&
			v.Value.Content != nil &&
			v.Value.Content[openapi.ApplicationJson] != nil &&
			v.Value.Content[openapi.ApplicationJson].Schema != nil {
			types = append(types, convertSchema(name, v.Value.Content[openapi.ApplicationJson].Schema))
		}
	}
	return
}

func convertSchema(name string, schema *openapi3.SchemaRef) (t *Type) {
	
	return
}

const typesTpl = `

`
