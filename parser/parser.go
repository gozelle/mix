package parser

import (
	"fmt"
	"github.com/gozelle/fs"
)

type Method struct {
	Name    string
	Params  []*Type
	Results []*Type
}

type Type struct {
	Names    []string
	Type     string
	Tags     string
	Reserved bool
	Pointer  bool
	Length   int
	Slice    bool
	Ellipsis bool
	Fields   []*Type // for struct
	Elem     *Type   // for pointer and slice or array
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

func (p *Parser) findRootType(name string) *Type {
	if p.Root == nil {
		return nil
	}
	if p.Root.Types == nil {
		return nil
	}
	return p.Root.Types[name]
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
