package parser

import (
	"github.com/gozelle/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	
	file := filepath.Join(wd, "../test/rpc.go")
	
	err = ParseFile(file)
	require.NoError(t, err)
}
