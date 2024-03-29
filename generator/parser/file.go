package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

type File struct {
	path     string
	mod      *Mod
	pkg      *Package
	Imports  map[string]*Import
	comments ast.CommentMap
	ast      *ast.File
	set      *token.FileSet
}

func (f File) Errorf(pos token.Pos, format string, a ...any) error {
	p := f.set.Position(pos)
	return fmt.Errorf("parse file: %s error: %s", p, fmt.Sprintf(format, a...))
}

func (f *File) getImport(name string) *Import {
	if f.Imports == nil {
		return nil
	}
	return f.Imports[name]
}

func (f *File) addImport(item *Import) {
	if f.Imports == nil {
		f.Imports = map[string]*Import{}
	}
	
	f.Imports[item.Alias] = item
}

func (f *File) markToStringDef(node ast.Node) {
	d, ok := node.(*ast.FuncDecl)
	if !ok {
		return
	}
	if d.Recv == nil {
		return
	}
	if d.Recv.NumFields() == 0 {
		return
	}
	
	if d.Type.Results == nil {
		return
	}
	if d.Type.Results.NumFields() != 1 {
		return
	}
	s, ok := d.Type.Results.List[0].Type.(*ast.Ident)
	if !ok {
		return
	}
	if s.Name != "string" {
		return
	}
	
	i, ok := d.Recv.List[0].Type.(*ast.Ident)
	if !ok {
		return
	}
	
	f.pkg.markStringer(i.String())
}

func (f *File) Visit(node ast.Node) ast.Visitor {
	
	f.markToStringDef(node)
	
	s, ok := node.(*ast.TypeSpec)
	if !ok {
		return f
	}
	switch t := s.Type.(type) {
	case *ast.InterfaceType:
		f.pkg.addInterface(s.Name.String(), &Interface{Name: s.Name.String(), interfaceType: t, file: f})
	}
	
	switch t := s.Type.(type) {
	
	case *ast.InterfaceType,
		*ast.FuncType,
		*ast.StructType,
		*ast.Ident,
		*ast.MapType,
		*ast.SliceExpr,
		*ast.ArrayType,
		*ast.StarExpr,
		*ast.SelectorExpr,
		*ast.ChanType,
		*ast.IndexExpr,
		*ast.IndexListExpr:
		if !s.Name.IsExported() {
			return f
		}
		_, isStruct := t.(*ast.StructType)
		f.pkg.addDef(s.Name.String(), &Def{Name: s.Name.String(), File: f, Expr: s.Type, IsStrut: isStruct})
	
	default:
		panic(fmt.Errorf("unsupport parse type: %s at: %s", reflect.TypeOf(t), f.set.Position(s.Pos())))
	}
	
	return f
}

func (f *File) parse(file string) (err error) {
	set := token.NewFileSet()
	af, err := parser.ParseFile(set, file, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return
	}
	f.ast = af
	f.set = set
	f.comments = ast.NewCommentMap(set, af, af.Comments)
	
	for _, i := range af.Imports {
		v := f.parseImport(i)
		if v != nil {
			f.addImport(v)
		}
	}
	
	// parse file use Visit
	ast.Walk(f, af)
	
	return
}

func (f *File) parseAlias(name string) string {
	if name == "<nil>" {
		return ""
	}
	return name
}

func (f *File) parseImport(i *ast.ImportSpec) *Import {
	
	r := &Import{
		Alias: f.parseAlias(i.Name.String()),
		Path:  strings.Trim(i.Path.Value, "\""),
	}
	
	if strings.Contains(r.Path, "/internal") {
		return nil
	}
	
	r.Package = &Package{Path: r.Path}
	
	if r.Path == "C" {
		return r
	}
	
	var err error
	
	// 忽略 go.mod 不依赖的包
	realPath := f.mod.GetPackagePath(r.Path)
	if realPath == "" {
		return nil
	}
	
	if r.Alias == "" {
		r.Alias, err = f.mod.GetPackageRealName(realPath)
		if err != nil {
			panic(fmt.Errorf("%s: get %s package name error: %s", f.path, r.Path, err))
		}
	}
	
	if f.mod.packages != nil {
		if v, ok := f.mod.packages[realPath]; ok {
			r.Package = v
			return r
		}
	}
	defer func() {
		f.mod.cachePackage(realPath, r.Package)
	}()
	
	err = r.Package.Parse(f.mod, realPath)
	if err != nil {
		// panic(fmt.Errorf("load package: %s files from: %s error: %s", r.Path, realPath, err))
		return nil
	}
	return r
}
