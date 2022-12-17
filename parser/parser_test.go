package parser

import (
	"encoding/json"
	"fmt"
	"github.com/gozelle/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	
	parser, err := NewParser(filepath.Join(wd, "../test"), "TestRpc")
	require.NoError(t, err)
	
	p, err := parser.Parse()
	require.NoError(t, err)
	
	d, err := json.MarshalIndent(p, "", "\t")
	require.NoError(t, err)
	
	fmt.Println(string(d))
}
