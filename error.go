package mix

import (
	"fmt"
	"github.com/gozelle/jsonrpc"
)

type Error = jsonrpc.Error

func Errorf(format string, a ...any) *Error {
	return &Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func Codef(code int, format string, a ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}
