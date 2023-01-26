package tests_basic

import (
	"encoding/json"
	"github.com/gozelle/fastjson"
	"github.com/gozelle/fs"
	"github.com/gozelle/mix/generator/parser"
	"github.com/gozelle/mix/generator/render"
	"github.com/gozelle/testify/require"
	"testing"
)

func TestBasic(t *testing.T) {
	mod, err := parser.PrepareMod()
	require.NoError(t, err)
	
	dir, err := fs.Lookup("./generator/tests/basic")
	require.NoError(t, err)
	
	pkg, err := parser.Parse(mod, dir)
	require.NoError(t, err)
	
	defJson, err := json.Marshal(pkg.GetDef("Basic"))
	require.NoError(t, err)
	
	c, err := fs.Read("./parser_basic_def.json")
	require.NoError(t, err)
	
	// 比较 Parser 类型定义
	err = fastjson.EqualsBytes(c, defJson)
	require.NoError(t, err)
	
	parserInterface := pkg.GetInterface("BasicAPI")
	require.True(t, parserInterface != nil)
	
	renderInterface := render.ToGolangInterface(parserInterface)
	d, _ := json.MarshalIndent(renderInterface, "", "\t")
	
	err = fs.Write("./render_interface.json", d)
	require.NoError(t, err)
	
	t.Log(string(d))
}
