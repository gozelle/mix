package parser

import (
	"fmt"
	"github.com/gozelle/fs"
	"strings"
)

type Method struct {
	Name    string
	Params  []*Param
	Results []*Param
}

type Type string

func (t Type) IsStruct() bool {
	return string(t) == "struct"
}

func (t Type) IsArray() bool {
	return strings.HasPrefix(string(t), "[]")
}

func (t Type) IsContext() bool {
	return string(t) == "context.Context"
}

func (t Type) IsError() bool {
	return string(t) == "error"
}

type Param struct {
	Names []string
	Type  Type
}

type Def struct {
	Name   string
	Type   Type
	Tags   string
	Fields []*Def
}

func NewParser(dir string) (parser *Parser, err error) {
	
	mod, err := FindModFile()
	if err != nil {
		return
	}
	
	parser = &Parser{
		mod: mod,
	}
	
	parser.Root, err = parser.loadPackage(dir)
	if err != nil {
		return
	}
	
	return
}

type Parser struct {
	mod      *Mod
	Root     *Package
	Packages []*Package
}

func (p *Parser) findRootType(name string) *Def {
	if p.Root == nil {
		return nil
	}
	if p.Root.Defs == nil {
		return nil
	}
	return p.Root.Defs[name]
}

func (p *Parser) findRootImport(name string) *Package {
	if p.Root == nil {
		return nil
	}
	if p.Root.Imports == nil {
		return nil
	}
	
	return p.Root.Imports[name]
}

func (p *Parser) CombineInterface(name string) (*Interface, error) {
	
	if p.Root == nil {
		return nil, fmt.Errorf("root package is nil")
	}
	if p.Root.Interfaces == nil {
		return nil, fmt.Errorf("root package not contains interface")
	}
	i, ok := p.Root.Interfaces[name]
	if !ok {
		return nil, fmt.Errorf("interface: %s not found", name)
	}
	
	err := i.load(p)
	if err != nil {
		return nil, err
	}
	
	return i, nil
}

func (p *Parser) loadPackage(dir string) (*Package, error) {
	ok, err := fs.IsDir(dir)
	if err != nil || !ok {
		return nil, fmt.Errorf("only accept dir")
	}
	files, err := fs.Files(dir, ".go")
	if err != nil {
		return nil, err
	}
	pkg := &Package{
		mod: p.mod,
	}
	err = pkg.loadPackage(files)
	if err != nil {
		return nil, err
	}
	
	return pkg, nil
}
