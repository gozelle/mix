package mix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gozelle/gin"
	"github.com/gozelle/jsonrpc"
	"io"
	"io/ioutil"
	"net/http"
)

const placeholder = "$params$"
const method = "method"
const module = "module"

// 包装 API
// ns == ""  =>  gin.Post("/api/v0/:method")
// ns != ""  =>  gin.Post("/api/v0/:namespace/:method")
func wrapAPI(ns string, h http.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := ctx.Request.Body
		d, err := ioutil.ReadAll(body)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("read request body error:%s", err))
		}
		var m string
		if ns == "" {
			m = fmt.Sprintf(".%s", ctx.Param(method))
		} else {
			m = fmt.Sprintf("%s.%s", ctx.Param(module), ctx.Param(method))
		}
		
		r := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  m,
			"params":  placeholder,
		}
		i, err := json.Marshal(r)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("prepare params error:%s", err))
			return
		}
		var params []byte
		if len(d) > 0 {
			//params = bytes.Join([][]byte{{91}, d, {93}}, []byte{})
			params = d
		} else {
			params = []byte{91, 93}
		}
		i = bytes.Replace(i, []byte(fmt.Sprintf("\"%s\"", placeholder)), params, 1)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(i))
		
		if err != nil {
			log.Errorf(" error: %s", err)
			return
		}
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func WrapHandler(wrap func(ctx *gin.Context) (any, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r, err := wrap(ctx)
		if err != nil {
			if e, ok := err.(*Error); ok {
				ctx.Header(jsonrpc.X_RPC_ERROR, e.Message)
				ctx.JSON(http.StatusBadRequest, &jsonrpc.Response{
					ID:    ctx.Writer.Header().Get(jsonrpc.X_RPC_ID),
					Error: e,
				})
			} else {
				HandleServerError(ctx, err)
			}
			return
		}
		if bs, ok1 := r.([]byte); ok1 {
			_, err = ctx.Writer.Write(bs)
			if err != nil {
				HandleServerError(ctx, err)
				return
			}
		} else if reader, ok2 := r.(io.Reader); ok2 {
			var d []byte
			d, err = ioutil.ReadAll(reader)
			if err != nil {
				HandleServerError(ctx, err)
				return
			}
			_, err = ctx.Writer.Write(d)
			if err != nil {
				HandleServerError(ctx, err)
				return
			}
		} else {
			ctx.JSON(200, &jsonrpc.Response{
				ID:     ctx.Writer.Header().Get(jsonrpc.X_RPC_ID),
				Result: r,
			})
		}
	}
}

func HandleServerError(ctx *gin.Context, err error) {
	ctx.Header(jsonrpc.X_RPC_ERROR, err.Error())
	ctx.JSON(http.StatusInternalServerError, &jsonrpc.Response{
		ID: ctx.Writer.Header().Get(jsonrpc.X_RPC_ID),
		Error: &jsonrpc.Error{
			Code:    500,
			Message: err.Error(),
		},
	})
}
