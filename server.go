package mix

import (
	"fmt"
	"github.com/gozelle/gin"
	"github.com/gozelle/jsonrpc"
	"github.com/gozelle/mix/middlewares/cors"
	"path/filepath"
	"strings"
)

func NewServer() *Server {
	return &Server{}
}

type Server struct {
	*gin.Engine
}

func (s *Server) init() {
	if s.Engine == nil {
		s.Engine = gin.New()
		s.Engine.Use(Logger(), gin.Recovery(), cors.Cors())
	}
}

func (s *Server) RegisterRPC(path, namespace string, handler any, middlewares ...gin.HandlerFunc) {
	s.init()
	
	rpcServer := jsonrpc.NewServer()
	rpcServer.Register(namespace, handler)
	s.Engine.POST(path, append([]gin.HandlerFunc{gin.WrapH(rpcServer)}, middlewares...)...)
}

func (s *Server) RegisterAPI(path, namespace string, handler any, middlewares ...gin.HandlerFunc) {
	s.init()
	
	if strings.TrimSpace(namespace) == "" {
		path = fmt.Sprintf("%s/:%s", filepath.Clean(path), method)
	} else {
		path = fmt.Sprintf("%s/:%s/:%s", filepath.Clean(path), module, method)
	}
	
	rpcServer := jsonrpc.NewServer()
	rpcServer.Register(namespace, handler)
	
	s.Engine.POST(path, append([]gin.HandlerFunc{wrapAPI(namespace, rpcServer)}, middlewares...)...)
}
