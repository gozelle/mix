package render

type Interface struct {
	Package string    `json:"Package,omitempty"`
	Name    string    `json:"Name,omitempty"`
	Methods []*Method `json:"Methods,omitempty"`
	Defs    []*Def    `json:"Defs,omitempty"`
	Imports []*Import `json:"Imports,omitempty"`
}

type Method struct {
	Name    string `json:"Name,omitempty"`
	Request *Def   `json:"Request,omitempty"`
	Replay  *Def   `json:"Replay,omitempty"`
	Params  string `json:"Params,omitempty"`
	Results string `json:"Results,omitempty"`
}

type Param struct {
	Names []string `json:"Names,omitempty"`
	Type  string   `json:"Type,omitempty"`
}

type Def struct {
	Name    string `json:"Name,omitempty"` // 定义时的名称
	Field   string `json:"Field,omitempty"`
	Json    string `json:"Json,omitempty"`
	Type    string `json:"Type"`
	Pointer bool   `json:"Pointer,omitempty"`
	//Reserved     bool   `json:"Reserved,omitempty"`
	StructFields []*Def `json:"StructFields,omitempty"`
	Elem         *Def   `json:"Elem,omitempty"`
	Use          *Def   `json:"Use,omitempty"`
	Tags         string `json:"Tags,omitempty"`
}

type Import struct {
	Alias string `json:"Alias,omitempty"`
	Path  string `json:"Path,omitempty"`
}
