package parser

import (
	jsonrpc_client "github.com/gozelle/mix/maker/jsonrpc"
	"github.com/gozelle/mix/parser"
	"github.com/gozelle/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	
	_parser, err := parser.NewParser(filepath.Join(wd, "../rpc"), "TestRpc")
	require.NoError(t, err)
	
	p, err := _parser.Parse()
	require.NoError(t, err)
	
	err = p.Make("TestRpc", jsonrpc_client.Maker{})
	require.NoError(t, err)
	
	//d, err := json.MarshalIndent(p, "", "\t")
	//require.NoError(t, err)
	//
	//fmt.Println(string(d))
}
