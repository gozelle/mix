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

type GenFile struct {
	Name    string
	Content string
}

type Generator interface {
	Generate(i *Interface) (files []*GenFile, err error)
}

type Interface struct {
	Name     string
	Methods  []*Method
	Types    []*Type
	Packages []*Package
	includes []*Interface
	t        *ast.InterfaceType
}

type Method struct {
	Name    string
	Params  []*Type
	Results []*Type
}

type Type struct {
	Name     string
	Type     string
	Tags     string
	Reserved bool
	Pointer  bool
	Length   int
	Slice    bool
	Fields   []*Type // for struct
	Elem     *Type   // for pointer and slice or array
}

type Package struct {
	Name       string
	Alias      string
	Path       string
	Imports    map[string]*Package
	Interfaces map[string]*Interface
	Types      map[string]*Type
	mod        *Mod
	loaded     bool
}

func (p *Package) Generate(target string, maker Generator) ([]*GenFile, error) {
	i := p.getInterface(target)
	if i == nil {
		return nil, fmt.Errorf("interface: '%s' not found", target)
	}

	for _, m := range i.t.Methods.List {
		switch mt := m.Type.(type) {
		case *ast.Ident:
		// TODO parse include
		case *ast.FuncType:
			i.Methods = append(i.Methods, p.parseMethod(m.Names[0].Name, mt))
		}
	}

	return maker.Generate(i)
}

func (p *Package) addImports(name string, item *Package) *Package {
	if p.Imports == nil {
		p.Imports = map[string]*Package{}
	}
	if v, ok := p.Imports[name]; ok {
		return v
	}
	p.Imports[name] = item
	return item
}

func (p *Package) addType(name string, item *Type) *Type {
	if p.Types == nil {
		p.Types = map[string]*Type{}
	}
	if v, ok := p.Types[name]; ok {
		return v
	}
	p.Types[name] = item
	return item
}

func (p *Package) addInterface(name string, item *Interface) *Interface {
	if p.Interfaces == nil {
		p.Interfaces = map[string]*Interface{}
	}
	if v, ok := p.Interfaces[name]; ok {
		return v
	}
	p.Interfaces[name] = item
	return item
}

func (p *Package) getInterface(name string) *Interface {
	if p.Interfaces == nil {
		return nil
	}
	v, ok := p.Interfaces[name]
	if ok {
		return v
	}
	return nil
}

type Parser struct {
	mod *Mod
}

func NewParser() (parser *Parser, err error) {

	mod, err := FindModFile()
	if err != nil {
		return
	}

	parser = &Parser{
		mod: mod,
	}

	return
}

func (r Parser) LoadPackage(dir string) (p *Package, err error) {
	ok, err := fs.IsDir(dir)
	if err != nil || !ok {
		err = fmt.Errorf("only accept dir")
		return
	}
	files, err := fs.Files(dir, ".go")
	if err != nil {
		return
	}
	p = &Package{
		mod: r.mod,
	}
	err = p.loadPackage(files)
	if err != nil {
		return
	}

	return
}

func (p *Package) loadPackage(files []string) (err error) {
	for _, v := range files {
		if !strings.HasSuffix(v, "_test.go") {
			err = p.parseFile(v)
			if err != nil {
				return
			}
		}
	}
	p.loaded = true
	return
}

func (p *Package) parseFile(file string) (err error) {
	//defer func() {
	//	if e := recover(); e != nil {
	//		err = fmt.Errorf("%v", e)
	//		return
	//	}
	//}()

	set := token.NewFileSet()
	f, err := parser.ParseFile(set, file, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return
	}

	p.Name = f.Name.String()
	for _, i := range f.Imports {
		v := p.parseImport(i)
		p.addImports(v.Name, v)
	}
	ast.Walk(p, f)

	//for _, v := range p.Interfaces {
	//	p.parseInterface(v)
	//}

	return
}

func (p *Package) parseImport(i *ast.ImportSpec) *Package {

	r := &Package{
		Alias: strings.Trim(i.Name.String(), "<>"),
		Path:  strings.Trim(i.Path.Value, "\""),
		mod:   p.mod,
	}
	var err error
	r.Name, err = r.mod.GetPackageRealName(r.Path)
	if err != nil {
		panic(fmt.Errorf("get %s package name error: %s", r.Path, err))
	}
	if r.Alias != "" && r.Alias != r.Name {
		r.Name = r.Alias
	}

	return r
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
		p.addInterface(s.Name.String(), &Interface{Name: s.Name.String(), t: t})
	case *ast.FuncType:
		return p
	case *ast.StructType, *ast.Ident:
		p.addType(s.Name.String(), p.parseType(s.Name.String(), s.Type))
	case *ast.MapType:
		// TODO
	case *ast.SliceExpr:
		// TODO
	case *ast.ArrayType:
	// TODO
	case *ast.SelectorExpr:
		fmt.Println("selector:", s.Name, t.X.(*ast.Ident).Name, t.Sel.Name)
		//p.addType(s.Name.String(), p.parseType(s.Name.String(), s.Type))
	default:
		fmt.Println("package:", t, reflect.TypeOf(t).String())
	}

	return p
}

func (p *Package) getPackage(name string) *Package {
	if p.Imports == nil {
		p.Imports = map[string]*Package{}
	}
	v, ok := p.Imports[name]
	if !ok {
		return nil
	}
	if !v.loaded {
		files, err := p.mod.OpenPackage(v.Path)
		if err != nil {
			panic(fmt.Errorf("open package '%s' error: %s", v.Path, err))
		}
		err = v.loadPackage(files)
		if err != nil {
			panic(fmt.Errorf("load package '%s' error: %s", v.Path, err))
		}
	}

	return v
}

func (p *Package) parseMethod(name string, t *ast.FuncType) (r *Method) {
	r = &Method{Name: name}

	for _, f := range t.Params.List {
		r.Params = append(r.Params, p.parseField(p.parseNames(f.Names), f)...)
	}

	for _, f := range t.Results.List {
		r.Results = append(r.Results, p.parseField(p.parseNames(f.Names), f)...)
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

func (p *Package) parseField(names []string, t *ast.Field) (r []*Type) {

	for _, n := range names {
		r = append(r, p.parseType(n, t.Type))
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
			r.Fields = append(r.Fields, st)
		}
	case *ast.SelectorExpr:
		// TODO 处理包引用类型
		r.Type = fmt.Sprintf("%s.%s", e.X.(*ast.Ident), e.Sel.Name)
	case *ast.StarExpr:
		r.Pointer = true
		r.Elem = p.parseType(name, e.X)
	default:
		// TODO 报错，未知的类型
		fmt.Println(e)

		panic(fmt.Errorf("unknown type: %s", reflect.TypeOf(e)))
	}

	return
}
