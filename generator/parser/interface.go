package parser

import (
	"fmt"
	"go/ast"
)

type Interface struct {
	Name          string
	Methods       []*Method
	Defs          []*Def
	defs          map[string]bool
	Imports       []*Import
	imports       map[string]bool
	includes      []*Interface
	interfaceType *ast.InterfaceType
	file          *File
}

func (i *Interface) load(mod *Mod, pkg *Package, file *File) (err error) {
	for _, m := range i.interfaceType.Methods.List {
		switch mt := m.Type.(type) {
		case *ast.Ident:
		// TODO parse include
		case *ast.FuncType:
			i.Methods = append(i.Methods, i.parseMethod(file, m.Names[0].Name, mt))
		}
	}
	for _, v := range pkg.Defs {
		i.addDef(v)
	}
	return
}

func (i *Interface) parseMethod(file *File, name string, t *ast.FuncType) (r *Method) {
	r = &Method{Name: name}
	
	if t.Params != nil {
		for index, f := range t.Params.List {
			names := parseNames(f.Names)
			if len(names) == 0 {
				names = append(names, fmt.Sprintf("p%d", index))
			}
			r.Params = append(r.Params, i.parseParam(file, names, f.Type))
		}
	}
	
	if t.Results != nil {
		for index, f := range t.Results.List {
			names := parseNames(f.Names)
			if len(names) == 0 {
				names = append(names, fmt.Sprintf("r%d", index))
			}
			r.Results = append(r.Results, i.parseParam(file, names, f.Type))
		}
	}
	return
}

func (i *Interface) addDef(t *Def) {
	if i.defs == nil {
		i.defs = map[string]bool{}
	}
	if _, ok := i.defs[t.Name]; ok {
		return
	}
	i.defs[t.Name] = true
	i.Defs = append(i.Defs, t)
}

func (i *Interface) parseParam(f *File, names []string, t ast.Expr) (r *Param) {
	
	r = &Param{
		Names: names,
	}
	
	defer func() {
		log.Infof("param: %v Type: %s", r.Names, r.Type)
	}()
	
	r.Type = parseType(f, "", t)
	
	return
}

func (i *Interface) addImport(imt *Import) {
	if i.imports == nil {
		i.imports = map[string]bool{}
	}
	if _, ok := i.imports[imt.Path]; ok {
		return
	}
	i.imports[imt.Path] = true
	i.Imports = append(i.Imports, imt)
}

func isReserved(t string) bool {
	switch t {
	case "int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"float32",
		"float64",
		"string",
		"bool",
		"byte",
		"rune",
		"uintptr",
		"any",
		"error":
		return true
	}
	return false
}
