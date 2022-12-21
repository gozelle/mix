package parser

import (
	"github.com/gozelle/mix/generator/jsonrpc-client"
	"github.com/gozelle/mix/parser"
	"github.com/gozelle/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	
	_parser, err := parser.NewParser(filepath.Join(wd, "../rpc"))
	require.NoError(t, err)
	
	i, err := _parser.CombineInterface("TestRpc")
	require.NoError(t, err)
	
	g := gen_jsonrpc_client.Generator{}
	files, err := g.Generate(i)
	require.NoError(t, err)
	
	t.Log(files[0].Content)
	
}
