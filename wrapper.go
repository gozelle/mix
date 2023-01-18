package mix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gozelle/gin"
	"io"
	"io/ioutil"
	"net/http"
	"time"
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
			"id":      time.Now().UnixNano(), // TODO use request header X-Request-Id value
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
			params = bytes.Join([][]byte{{91}, d, {93}}, []byte{})
		} else {
			params = []byte{91, 93}
		}
		i = bytes.Replace(i, []byte(fmt.Sprintf("\"%s\"", placeholder)), params, 1)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(i))
		
		err = ctx.Error(fmt.Errorf("some error"))
		if err != nil {
			log.Errorf(" error: %s", err)
			return
		}
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func HandleReader(wrapper func(ctx *gin.Context) io.Reader) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := wrapper(ctx)
		d, err := ioutil.ReadAll(r)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		}
		_, _ = ctx.Writer.Write(d)
	}
}
