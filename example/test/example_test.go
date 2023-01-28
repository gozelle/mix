package test

import (
	"github.com/gozelle/fs"
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
		"dist",
		"--outfile",
		fs.Join(examplePath, "dist", "proxy_gen.go"),
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
