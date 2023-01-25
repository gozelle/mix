package openapi

import (
	"github.com/gozelle/fs"
	"github.com/gozelle/mix/generator/parser"
	"github.com/gozelle/testify/require"
	"testing"
)

func TestParser(t *testing.T) {
	
	mod, err := parser.PrepareMod()
	require.NoError(t, err)
	
	path, err := fs.Lookup("./example")
	require.NoError(t, err)
	
	p, err := parser.Parse(mod, path)
	require.NoError(t, err)
	
	i := p.GetInterface("TestRpc")
	
	g := Generator{}
	
	v3 := g.TOOpenapiV3(i)
	
	d, err := v3.MarshalJSON()
	require.NoError(t, err)
	
	t.Log(string(d))
	
	//t.Log(r.String())
	//files, err := g.Generate(i)
	//require.NoError(t, err)
	//t.Log(files[0].Content)
}
