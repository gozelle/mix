package openapi

import (
	"fmt"
	"github.com/gozelle/logging"
	"github.com/gozelle/mix/generator/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var log = logging.Logger("convertor")

func ConvertAPI(i *parser.Interface) *API {
	
	r := &API{
		Name: i.Name,
	}
	
	//for _, v := range i.Imports {
	//	r.Imports = append(r.Imports, &Import{
	//		Alias: v.Alias,
	//		Path:  v.Path,
	//	})
	//}
	
	for _, v := range i.Defs {
		if v.Type.Type != parser.TStruct {
			continue
		}
		d := convertRenderDef(v)
		r.Defs = append(r.Defs, d)
		//fmt.Printf("===============Start Def: %s  ===============\n", v.Name)
		//dd, _ := json.MarshalIndent(d, "", "\t")
		//fmt.Println(string(dd))
		//fmt.Printf("===============End Def: %s ===============\n", v.Name)
	}
	
	for _, v := range i.Methods {
		r.Methods = append(r.Methods, convertRenderMethod(v))
	}
	
	return r
}

func convertRenderMethod(m *parser.Method) *Method {
	
	r := &Method{
		Name:    m.Name,
		Request: convertRenderMethodRequest(m),
		Replay:  convertRenderMethodReply(m),
	}
	return r
}

func convertRenderMethodRequest(m *parser.Method) *Def {
	request := &Def{
		Field: fmt.Sprintf("%sRequest", m.Name),
		Type:  parser.TStruct,
	}
	params := m.ExportParams()
	if len(params) == 1 && params[0].Type.NoPointer().Def != nil && params[0].Type.NoPointer().Def.Type.RealType().IsStruct() {
		request.Use = convertRenderDef(params[0].Type.NoPointer().Def)
	} else if len(params) > 0 {
		for _, v := range params {
			request.StructFields = append(request.StructFields, convertRenderMethodParam(v)...)
		}
	} else {
		request = nil
	}
	
	return request
}

func convertRenderMethodReply(m *parser.Method) *Def {
	replay := &Def{
		Field: fmt.Sprintf("%sReplay", m.Name),
		Type:  parser.TStruct,
	}
	
	results := m.ExportResults()
	
	if len(results) > 0 && results[0].Type.NoPointer().Def != nil && results[0].Type.NoPointer().Def.Type.RealType().IsStruct() {
		replay.Use = convertRenderDef(results[0].Type.NoPointer().Def)
	} else if len(results) > 0 {
		for _, v := range results {
			replay.StructFields = append(replay.StructFields, convertRenderMethodParam(v)...)
		}
	} else {
		replay = nil
	}
	return replay
}

func convertRenderMethodParam(p *parser.Param) []*Def {
	r := make([]*Def, 0)
	//log.Infof("convertRenderMethodParam: %v", p.Names)
	//spew.Json(p)
	for _, v := range p.Names {
		d := convertRenderType(p.Type)
		d.Field = Title(v)
		d.Json = v
		r = append(r, d)
	}
	return r
}

func Title(v string) string {
	return cases.Title(language.English).String(v)
}

func convertRenderDef(d *parser.Def) *Def {
	n := convertRenderType(d.Type)
	n.Name = d.Name
	return n
}

func convertRenderType(t *parser.Type) *Def {
	
	pointer := t.Pointer
	rt := t.NoPointer().RealType()
	
	n := &Def{
		Field:        t.Field,
		Json:         rt.Json(),
		Type:         rt.Type,
		Pointer:      pointer,
		StructFields: nil,
		Elem:         nil,
		Use:          nil,
		Tags:         string(rt.Tags),
	}
	
	switch rt.Type {
	case parser.TStruct:
		for _, v := range rt.StructFields {
			n.StructFields = append(n.StructFields, convertRenderType(v))
		}
	case parser.TSlice, parser.TArray:
		n.Type = parser.TSlice
		n.Elem = convertRenderType(rt.Elem)
	}
	
	// 处理引用关系
	if rt.Def != nil && rt.Def.Name == rt.Type {
		n.Use = &Def{Name: rt.Def.Name}
	}
	
	return n
}