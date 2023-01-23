package parser

import (
	"github.com/gozelle/mix/generator/clients/jsonrpc-client"
	parser2 "github.com/gozelle/mix/generator/parser"
	"github.com/gozelle/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	
	mod, err := parser2.PrepareMod()
	require.NoError(t, err)
	
	_parser, err := parser2.NewParser(mod, filepath.Join(wd, "../rpc"))
	require.NoError(t, err)
	
	i, err := _parser.CombineInterface("TestRpc")
	require.NoError(t, err)
	
	g := jsonrpc_client.Generator{}
	files, err := g.Generate(i)
	require.NoError(t, err)
	
	t.Log(files[0].Content)
	
}
