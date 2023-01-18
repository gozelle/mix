package mix

import (
	"fmt"
	"github.com/gozelle/gin"
	"github.com/gozelle/jsonrpc"
	"github.com/gozelle/mix/middlewares/cors"
	"strings"
)

func NewServer() *Server {
	s := &Server{}
	s.Engine = gin.New()
	s.Engine.Use(Logger(), gin.Recovery(), cors.Cors())
	return s
}

type Server struct {
	*gin.Engine
}

func RegisterRPC(router gin.IRouter, namespace string, handler any, middlewares ...gin.HandlerFunc) {
	rpcServer := jsonrpc.NewServer()
	rpcServer.Register(namespace, handler)
	router.POST("", append([]gin.HandlerFunc{gin.WrapH(rpcServer)}, middlewares...)...)
}

func RegisterAPI(router gin.IRouter, namespace string, handler any, middlewares ...gin.HandlerFunc) {
	
	var path string
	if strings.TrimSpace(namespace) == "" {
		path = fmt.Sprintf("/:%s", method)
	} else {
		path = fmt.Sprintf("/:%s/:%s", module, method)
	}
	
	rpcServer := jsonrpc.NewServer()
	rpcServer.Register(namespace, handler)
	
	router.POST(path, append([]gin.HandlerFunc{wrapAPI(namespace, rpcServer)}, middlewares...)...)
}
