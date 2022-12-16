package parser

import (
	"fmt"
	"github.com/gozelle/fs"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

type Interface struct {
	Name     string       `json:"Name,omitempty"`
	Methods  []*Method    `json:"Methods,omitempty"`
	Includes []*Interface `json:"Includes,omitempty"`
	t        *ast.InterfaceType
}

type Method struct {
	Name   string   `json:"Name,omitempty"`
	Params []*Field `json:"Params,omitempty"`
	Return []*Field `json:"Return,omitempty"`
}

type Field struct {
	Name string `json:"Name,omitempty"`
	Type *Type  `json:"Type,omitempty"`
}

type Type struct {
	Name     string  `json:"Name,omitempty"`
	Type     string  `json:"Type,omitempty"`
	Tags     string  `json:"Tags,omitempty"`
	Reserved bool    `json:"Reserved"`
	Children []*Type `json:"Children,omitempty"` // only for struct
}

type Package struct {
	Name       string       `json:"Name,omitempty"`
	Alias      string       `json:"Alias,omitempty"`
	Path       string       `json:"Path,omitempty"`
	Remote     bool         `json:"Remote"`
	Imports    []*Package   `json:"Imports,omitempty"`
	Interfaces []*Interface `json:"Interfaces,omitempty"`
	Types      []*Type      `json:"Types,omitempty"`
}

func ParseDir(dir string) (p *Package) {
	files, err := fs.Files(dir, ".go")
	if err != nil {
		return
	}
	_ = files
	return
}

func ParseFile(file string) (p *Package) {
	
	set := token.NewFileSet()
	f, err := parser.ParseFile(set, file, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return
	}
	
	p = &Package{
		Name: f.Name.String(),
	}
	for _, i := range f.Imports {
		p.Imports = append(p.Imports, &Package{
			Alias: strings.Trim(i.Name.String(), "<>"),
			Path:  strings.Trim(i.Path.Value, "\""),
		})
	}
	ast.Walk(p, f)
	
	for _, v := range p.Interfaces {
		p.parseInterface(v)
	}
	
	return
}

func (p *Package) GetTypePackage(t *Type) *Package {
	return nil
}

func (p *Package) Visit(node ast.Node) ast.Visitor {
	s, ok := node.(*ast.TypeSpec)
	if !ok {
		return p
	}
	
	switch t := s.Type.(type) {
	
	case *ast.InterfaceType:
		p.Interfaces = append(p.Interfaces, &Interface{Name: s.Name.String(), t: t})
	case *ast.FuncType:
		return p
	case *ast.StructType, *ast.Ident:
		p.Types = append(p.Types, p.parseType(s.Name.String(), s.Type))
	case *ast.MapType:
		// TODO
	case *ast.SliceExpr:
		// TODO
	case *ast.ArrayType:
	// TODO
	case *ast.SelectorExpr:
		fmt.Println("selector:", s.Name, t.X.(*ast.Ident).Name, t.Sel.Name)
	default:
		
		fmt.Println("package:", t, reflect.TypeOf(t).String())
	}
	
	return p
}

func (p *Package) parseInterface(r *Interface) {
	for _, m := range r.t.Methods.List {
		switch mt := m.Type.(type) {
		case *ast.Ident:
		// TODO parse include
		case *ast.FuncType:
			r.Methods = append(r.Methods, p.parseMethod(m.Names[0].Name, mt))
		}
	}
	return
}

func (p *Package) parseMethod(name string, t *ast.FuncType) (r *Method) {
	r = &Method{Name: name}
	
	for _, f := range t.Params.List {
		r.Params = append(r.Params, p.parseField(p.parseNames(f.Names), f)...)
	}
	
	for _, f := range t.Results.List {
		r.Return = append(r.Return, p.parseField(p.parseNames(f.Names), f)...)
	}
	
	return
}

func (p *Package) parseNames(idents []*ast.Ident) []string {
	names := make([]string, 0)
	for _, i := range idents {
		names = append(names, i.Name)
	}
	return names
}

func (p *Package) parseField(names []string, t *ast.Field) (r []*Field) {
	
	for _, n := range names {
		f := &Field{
			Name: n,
			Type: p.parseType("", t.Type), // TODO 复用, 定义基本类型
		}
		r = append(r, f)
	}
	
	return
}

func (p *Package) parseType(name string, t ast.Expr) (r *Type) {
	
	// TODO Cache package.Type
	r = &Type{
		Name: name,
	}
	
	switch e := t.(type) {
	case *ast.Ident:
		r.Type = e.Name
		r.Reserved = true
	case *ast.StructType:
		r.Type = "struct"
		for _, f := range e.Fields.List {
			st := p.parseType(p.parseNames(f.Names)[0], f.Type)
			if f.Tag != nil {
				st.Tags = f.Tag.Value
			}
			r.Children = append(r.Children, st)
		}
	
	case *ast.SelectorExpr:
		// TODO 处理包引用类型
		r.Type = fmt.Sprintf("%s.%s", e.X.(*ast.Ident), e.Sel.Name)
	default:
		// TODO 报错，未知的类型
		fmt.Println(e)
	}
	
	return
}
