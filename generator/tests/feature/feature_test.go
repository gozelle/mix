package tests_feature

import (
	"encoding/json"
	"github.com/gozelle/fs"
	"github.com/gozelle/mix/generator/openapi"
	"github.com/gozelle/mix/generator/parser"
	"github.com/gozelle/mix/generator/render"
	"github.com/gozelle/spew"
	"github.com/gozelle/testify/require"
	"testing"
)

func TestFeature(t *testing.T) {
	mod, err := parser.PrepareMod()
	require.NoError(t, err)
	
	dir, err := fs.Lookup("./generator/tests/feature")
	require.NoError(t, err)
	
	pkg, err := parser.Parse(mod, dir)
	require.NoError(t, err)
	
	defJson, err := json.Marshal(pkg.GetDef("Feature"))
	require.NoError(t, err)
	
	t.Log(string(defJson))
	
	//c, err := fs.Read("./parser_basic_def.json")
	//require.NoError(t, err)
	//
	//// 比较 Parser 类型定义
	//err = fastjson.EqualsBytes(c, defJson)
	//require.NoError(t, err)
	//
	parserInterface := pkg.GetInterface("FeatureAPI")
	require.True(t, parserInterface != nil)
	
	renderInterface := render.Convert(parserInterface)
	d, _ := json.MarshalIndent(renderInterface, "", "\t")
	t.Log(string(d))
	
	g := openapi.Generator{}
	
	v3 := g.TOOpenapiV3(renderInterface)
	
	d, err = v3.MarshalJSON()
	require.NoError(t, err)
	
	//
	t.Log(string(d))
	//err = fs.Write("./render_interface.json", d)
	//require.NoError(t, err)
	
	//t.Log(string(d))
}

type A struct {
	Name  string
	Child *A
}

func TestJson(t *testing.T) {
	a := A{
		Name: "a",
		Child: &A{
			Name: "b",
		},
	}
	spew.Json(a)
}