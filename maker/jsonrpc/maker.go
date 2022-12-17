package jsonrpc

import (
	"fmt"
	"github.com/gozelle/mix/parser"
)

type Maker struct {
}

func (m Maker) Make(methods []*parser.Method, typer *parser.Typer, packager *parser.Packager) (files []*parser.File, err error) {
	for _, v := range methods {
		fmt.Println(v.Name)
	}
	return
}

func renderMethod() {

}

func renderType() {

}
