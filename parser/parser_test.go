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
	
	file := filepath.Join(wd, "../test/rpc.go")
	
	v := ParseFile(file)
	d, err := json.MarshalIndent(v, "", "\t")
	require.NoError(t, err)
	
	fmt.Println(string(d))
}
