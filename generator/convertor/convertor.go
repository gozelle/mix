package convertor

import (
	"encoding/json"
	"fmt"
	"github.com/gozelle/logging"
	"github.com/gozelle/mix/generator/langs/golang"
	"github.com/gozelle/mix/generator/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	
	pointer := t.Pointer
	if pointer {
		t = t.NoPointer()
	}
	
	n := &golang.Def{
		Name:     name,
		Pointer:  pointer,
		Type:     t.RealType().Name,
		Reserved: t.Reserved,
		Tags:     t.Tags,
	}
	if n.Type == "[]" {
		n.Elem = convertType(t.Elem.Field, t.Elem)
	}
	
	return n
}

func parseRenderMethod(m *parser.Method) *golang.Method {
	
	request := &golang.Def{
		Name:   fmt.Sprintf("%sRequest", m.Name),
		Type:   golang.Struct,
		Concat: true,
	}
	replay := &golang.Def{
		Name:   fmt.Sprintf("%sReplay", m.Name),
		Type:   golang.Struct,
		Concat: true,
	}
	
	params := m.ExportParams()
	
	if len(params) == 1 && params[0].Type.IsStruct() {
		request = convertType(params[0].Type.Name, params[0].Type)
	} else if len(params) > 0 {
		for _, v := range params {
			request.StructFields = append(request.StructFields, convertMethodParam(v)...)
		}
	} else {
		request = nil
	}
	//log.Infof("request.Name: %s %v", request.Name, len(params) == 1 && params[0].Type.IsStruct())
	//log.Infof("params[0]: %s", params[0].Type.String())
	
	results := m.ExportResults()
	if len(results) > 0 && results[0].Type.IsStruct() {
		replay = convertType(results[0].Type.Name, results[0].Type)
	} else if len(results) > 0 {
		for _, v := range results {
			replay.StructFields = append(replay.StructFields, convertMethodParam(v)...)
		}
	} else {
		replay = nil
	}
	
	r := &golang.Method{
		Name:    m.Name,
		Request: request,
		Replay:  replay,
		//Params:  strings.Join(params, ","),
		//Results: strings.Join(results, ","),
	}
	//if len(results) > 0 {
	//	r.Results = fmt.Sprintf("(%s)", r.Results)
	//}
	
	return r
}

func convertMethodParam(p *parser.Param) []*golang.Def {
	r := make([]*golang.Def, 0)
	//log.Infof("convertMethodParam: %v", p.Names)
	//spew.Json(p)
	for _, v := range p.Names {
		d := convertType(Title(v), p.Type)
		d.Name = Title(v)
		d.Json = v
		r = append(r, d)
	}
	return r
}

func Title(v string) string {
	return cases.Title(language.English).String(v)
}
