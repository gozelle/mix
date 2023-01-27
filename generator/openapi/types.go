package openapi

import "github.com/gozelle/openapi/openapi3"

type DocumentV3 struct {
	openapi3.T
}

const (
	ApplicationJson = "application/json"
)
const (
	String = "string"
	Number = "number"
	Any    = "any"
	Object = "object"
	Array  = "array"
)
