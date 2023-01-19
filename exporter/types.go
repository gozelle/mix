package exporter

type Filed struct {
	Name   string
	Type   any
	Elem   string
	Fields []*Filed
}

type Elem struct {
	Type   string
	Fields []*Filed
}

type Type struct {
	Type any
}

const (
	String = "string"
	Number = "number"
	Any    = "any"
	Object = "object"
	Array  = "array"
)
