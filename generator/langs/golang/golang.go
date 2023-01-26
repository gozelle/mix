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

type Param struct {
	Names []string
	Type  string
}

type Def struct {
	Name         string `json:"name,omitempty"`
	Json         string `json:"json,omitempty"`
	Type         string `json:"type"`
	Pointer      bool   `json:"pointer,omitempty"`
	Reserved     bool   `json:"reserved,omitempty"`
	StructFields []*Def `json:"struct_fields,omitempty"`
	Elem         *Def   `json:"elem,omitempty"`
	Tags         string `json:"tags,omitempty"`
	Concat       bool   `json:"contact"`
}

type Package struct {
	Alias string
	Path  string
}
