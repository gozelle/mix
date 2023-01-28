package typescript_axios

import (
	"github.com/gozelle/fs"
	"github.com/gozelle/spew"
	"github.com/gozelle/testify/require"
	"testing"
)

func TestGenerate(t *testing.T) {
	file, err := fs.Lookup("./generator/tests/feature/openapi.json")
	require.NoError(t, err)
	
	files, err := Generate(file)
	require.NoError(t, err)
	
	spew.Json(files)
}
