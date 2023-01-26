package parser

import "fmt"

func Parse(mod *Mod, dir string) (pkg *Package, err error) {
	pkg = &Package{}
	err = pkg.load(mod, dir)
	if err != nil {
		return
	}
	
	for _, v := range pkg.Interfaces {
		err = v.load(mod, pkg, v.file)
		if err != nil {
			panic(fmt.Errorf("load interface: %s error: %s", v.Name, err))
		}
	}
	
	for _, v := range pkg.Defs {
		if !v.Used {
			continue
		}
		v.Type = parseType(v.File, "", v.Expr)
	}
	
	for name, v := range pkg.Defs {
		if v.Used {
			log.Infof("def: %s  => %s", name, v.String())
		}
	}
	
	return
}

//
//func NewParser(mod *Mod, dir string) (parser *Parser, err error) {
//	parser = &Parser{
//		Root: &Package{},
//	}
//	err = parser.Root.load(mod, dir)
//	if err != nil {
//		return
//	}
//	return
//}
//
//type Parser struct {
//	Root *Package
//}
//
//func (p *Parser) CombineInterface(name string) (*Interface, error) {
//
//	if p.Root == nil {
//		return nil, fmt.Errorf("root package is nil")
//	}
//	if p.Root.Interfaces == nil {
//		return nil, fmt.Errorf("root package not contains interface")
//	}
//	i, ok := p.Root.Interfaces[name]
//	if !ok {
//		return nil, fmt.Errorf("interface: %s not found", name)
//	}
//
//	//err := i.load(p.Root)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	return i, nil
//}
