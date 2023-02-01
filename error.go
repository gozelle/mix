package mix

import (
	"fmt"
	"github.com/gozelle/jsonrpc"
)

type Warn = jsonrpc.Error

func NewWarn(code int, message string, detail any) *Warn {
	return &Warn{Code: code, Message: message, Detail: detail}
}

func Warnf(format string, a ...any) *Warn {
	return &Warn{
		Message: fmt.Sprintf(format, a...),
	}
}

func Codef(code int, format string, a ...any) *Warn {
	return &Warn{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}
