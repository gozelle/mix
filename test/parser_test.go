package parser

import (
	jsonrpc_client "github.com/gozelle/mix/generator/jsonrpc"
	"github.com/gozelle/mix/parser"
	"github.com/gozelle/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	
	_parser, err := parser.NewParser()
	require.NoError(t, err)
	
	p, err := _parser.LoadPackage(filepath.Join(wd, "../rpc"))
	require.NoError(t, err)
	
	files, err := p.Generate("TestRpc", jsonrpc_client.Maker{})
	require.NoError(t, err)
	
	t.Log(files[0].Name, files[0].Content)
	//d, err := json.MarshalIndent(p, "", "\t")
	//require.NoError(t, err)
	//
	//fmt.Println(string(d))
}
