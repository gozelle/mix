package typescript_axios

import (
	"github.com/gozelle/fs"
	"github.com/gozelle/testify/require"
	"testing"
)

func TestGenerate(t *testing.T) {
	file, err := fs.Lookup("./generator/tests/feature/openapi.json")
	require.NoError(t, err)
	
	files, err := Generate(file, "")
	require.NoError(t, err)
	
	for _, v := range files {
		t.Log("็ๆๆไปถ:", v.Name)
		t.Log(v.Content)
	}
}
