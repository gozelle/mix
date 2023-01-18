package mix

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gozelle/gin"
	"github.com/gozelle/testify/require"
	"io"
	"testing"
)

type ITestHandler interface {
	Ping(ctx context.Context, msg string) (reply string, err error)
	Query(ctx context.Context, page, limit int) (reply string, err error)
	Download(ctx context.Context, file string) io.Reader
	Error(ctx context.Context) error
	Code(ctx context.Context) error
	Upload(ctx context.Context, file string, size int64, data []byte) (err error)
}

var _ ITestHandler = (*TestHandler)(nil)

type TestHandler struct {
}

func (t TestHandler) Upload(ctx context.Context, file string, size int64, data []byte) (err error) {
	//TODO implement me
	panic("implement me")
}

func (t TestHandler) Error(ctx context.Context) error {
	return fmt.Errorf("some error")
}

func (t TestHandler) Code(ctx context.Context) error {
	return &Error{
		Code:    1000,
		Message: "自定义错误",
		Detail:  "一些详情",
	}
}

func (t TestHandler) Query(ctx context.Context, page, limit int) (reply string, err error) {
	reply = fmt.Sprintf("page:%d limit:%d", page, limit)
	return
}

func (t TestHandler) Ping(ctx context.Context, msg string) (reply string, err error) {
	reply = fmt.Sprintf("received: %s", msg)
	return
}

func (t TestHandler) Download(ctx context.Context, file string) io.Reader {
	buf := &bytes.Buffer{}
	buf.WriteString("<h1>Hello world</h1>")
	return buf
}

func TestServer(t *testing.T) {
	h := &TestHandler{}
	
	server := NewServer()
	group := server.Group("/api/v1")
	
	RegisterRPC(server.Group("/rpc/v1"), "", h)
	RegisterAPI(group, "", h)
	
	group.GET("/download", WrapHandler(func(ctx *gin.Context) (data any, err error) {
		ctx.Header("Content-Type", "text/html; charset=UTF-8")
		return h.Download(ctx, "ok"), nil
	}))
	
	require.NoError(t, server.Run(":11111"))
}
