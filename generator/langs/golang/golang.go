package golang

const (
	Int     = "int"
	Int8    = "int8"
	Int16   = "int16"
	Int32   = "int32"
	Int64   = "int64"
	Uint    = "uint"
	Uint8   = "uint8"
	Uint16  = "uint16"
	Uint32  = "uint32"
	Uint64  = "uint64"
	Float32 = "float32"
	Float64 = "float64"
	Time    = "time"
	String  = "string"
	Struct  = "struct"
	Slice   = "slice"
	Array   = "array"
	Bool    = "bool"
)

type Interface struct {
	Package  string
	Name     string
	Methods  []*Method
	Defs     []*Def
	Packages []*Package
}

type Method struct {
	Name    string
	Request *Def
	Replay  *Def
	Params  string
	Results string
}

type Def struct {
	Name         string
	Type         string
	StructFields []*Def
	Elem         *Def
}

type Package struct {
	Alias string
	Path  string
}
