package parser

import (
	"github.com/gozelle/testify/require"
	"testing"
)

func TestFindMod(t *testing.T) {
	mod, err := PrepareMod()
	require.NoError(t, err)
	
	t.Log(mod.root, mod.file.Module)
	t.Log(mod.Gopath())
	
	files, err := mod.GetPackageFiles("github.com/gozelle/mix/origin")
	require.NoError(t, err)
	t.Log(files)
	
	files, err = mod.GetPackageFiles("github.com/gozelle/fs")
	require.NoError(t, err)
	t.Log(files)
}
