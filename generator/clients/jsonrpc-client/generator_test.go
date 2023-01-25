package jsonrpc_client

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
	files, err := g.Generate(i)
	require.NoError(t, err)
	
	t.Log(files[0].Content)
	
}
