package parser

import (
	"encoding/json"
	"github.com/gozelle/testify/require"
	"testing"
)

func TestFindMod(t *testing.T) {
	mod, err := PrepareMod()
	require.NoError(t, err)
	
	t.Log(mod.root, mod.file.Module)
	t.Log(mod.Gopath())
	
	files, err := mod.OpenPackage("github.com/gozelle/mix/origin")
	require.NoError(t, err)
	t.Log(files)
	
	files, err = mod.OpenPackage("github.com/gozelle/fs")
	require.NoError(t, err)
	t.Log(files)
}

type A struct {
}

var _ json.Marshaler
var _ json.Unmarshaler

func (a A) MarshalJSON() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
