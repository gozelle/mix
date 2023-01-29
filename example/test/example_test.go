package test

import (
	"encoding/json"
	"github.com/gozelle/fs"
	"github.com/gozelle/mix/generator/openapi"
	"github.com/gozelle/mix/generator/parser"
	typescript_axios "github.com/gozelle/mix/generator/sdks/typescript-axios"
	"github.com/gozelle/testify/require"
	"os"
	"os/exec"
	"testing"
)

func TestExample(t *testing.T) {
	testGenClient(t)
	testGenOpenapi(t)
	testGenSDK(t)
}

func testGenClient(t *testing.T) {
	examplePath, err := fs.Lookup("example")
	require.NoError(t, err)
	
	c := fs.Join(os.Getenv("GOROOT"), "bin/go")
	if !fs.Exists(c) {
		t.Fatal("can't find go cmd via $GOROOT")
	}
	_ = fs.MakeDir(fs.Join(examplePath, "dist"))
	cmd := exec.Command(
		c,
		"run",
		fs.Join(examplePath, "../cmd/mix/mix.go"),
		"generate",
		"client",
		"--path",
		fs.Join(examplePath, "api"),
		"--pkg",
		"example_api",
		"--outpkg",
		"example_api",
		"--outfile",
		fs.Join(examplePath, "api", "proxy_gen.go"),
	)
	cmd.Env = os.Environ()
	t.Log("exec", cmd.String())
	d, err := cmd.CombinedOutput()
	require.NoError(t, err)
	t.Log("exec Result", string(d))
}

func testGenOpenapi(t *testing.T) {
	examplePath, err := fs.Lookup("example")
	require.NoError(t, err)
	
	c := fs.Join(os.Getenv("GOROOT"), "bin/go")
	if !fs.Exists(c) {
		t.Fatal("can't find go cmd via $GOROOT")
	}
	_ = fs.MakeDir(fs.Join(examplePath, "dist"))
	cmd := exec.Command(
		c,
		"run",
		fs.Join(examplePath, "../cmd/mix/mix.go"),
		"generate",
		"openapi",
		"--path",
		fs.Join(examplePath, "api"),
		"--interface",
		"FullAPI",
		"--outfile",
		fs.Join(examplePath, "dist", "openapi.json"),
	)
	cmd.Env = os.Environ()
	t.Log("exec", cmd.String())
	d, err := cmd.CombinedOutput()
	require.NoError(t, err)
	t.Log("exec Result", string(d))
}

func testGenSDK(t *testing.T) {
	examplePath, err := fs.Lookup("example")
	require.NoError(t, err)
	
	c := fs.Join(os.Getenv("GOROOT"), "bin/go")
	if !fs.Exists(c) {
		t.Fatal("can't find go cmd via $GOROOT")
	}
	_ = fs.MakeDir(fs.Join(examplePath, "dist"))
	cmd := exec.Command(
		c,
		"run",
		fs.Join(examplePath, "../cmd/mix/mix.go"),
		"generate",
		"sdk",
		"--openapi",
		fs.Join(examplePath, "dist", "openapi.json"),
		"--sdk",
		"axios",
		"--outdir",
		fs.Join(examplePath, "dist", "sdk"),
	)
	cmd.Env = os.Environ()
	t.Log("exec", cmd.String())
	d, err := cmd.CombinedOutput()
	require.NoError(t, err)
	t.Log("exec Result", string(d))
}

func TestHandleGenOpenapi(t *testing.T) {
	mod, err := parser.PrepareMod()
	require.NoError(t, err)
	
	dir, err := fs.Lookup("./example/api")
	require.NoError(t, err)
	
	pkg, err := parser.Parse(mod, dir)
	require.NoError(t, err)
	
	api := pkg.GetInterface("FullAPI")
	require.True(t, api != nil)
	d, err := json.Marshal(api)
	require.NoError(t, err)
	t.Log(string(d))
	
	//c, err := fs.Read("./parser_basic_def.json")
	//require.NoError(t, err)
	//
	//// 比较 Parser 类型定义
	//err = fastjson.EqualsBytes(c, defJson)
	//require.NoError(t, err)
	//
	
	r := openapi.ConvertAPI(api)
	d, _ = json.Marshal(r)
	t.Log(string(d))
	
	doc := &openapi.DocumentV3{}
	
	openapi.ConvertOpenapi(doc, r)
	
	d, err = doc.MarshalJSON()
	require.NoError(t, err)
	t.Log(string(d))
	
	a, err := typescript_axios.Convert(doc)
	require.NoError(t, err)
	
	d, err = json.Marshal(a)
	require.NoError(t, err)
	t.Log(string(d))
}
