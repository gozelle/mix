package mix

import (
	"github.com/gozelle/gin"
	"github.com/gozelle/jsonrpc"
	"github.com/gozelle/logging"
	"github.com/gozelle/zap"
	"time"
)

var log = logging.Logger("[gin]")

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// SkipPaths is an url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

func Logger() gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

func LoggerWithConfig(conf LoggerConfig) gin.HandlerFunc {
	
	notlogged := conf.SkipPaths
	
	var skip map[string]struct{}
	
	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)
		
		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}
	
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		
		// Process request
		c.Next()
		
		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := gin.LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}
			param.Latency = time.Since(start)
			
			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.BodySize = c.Writer.Size()
			
			if raw != "" {
				path = path + "?" + raw
			}
			
			param.Path = path
			if param.Latency > time.Minute {
				param.Latency = param.Latency.Truncate(time.Second)
			}
			
			fields := []interface{}{
				zap.String("method", param.Method),
				zap.String("path", param.Path),
				zap.String("ip", param.ClientIP),
				zap.String("latency", param.Latency.String()),
				zap.Int("size", param.BodySize),
			}
			
			if h := c.Request.Header.Get(jsonrpc.X_RPC_Handler); h != "" {
				fields = append(fields, zap.String("handler", h))
			}
			
			if m := c.Writer.Header().Get(jsonrpc.X_RPC_ERROR); m != "" {
				fields = append(fields, zap.String("message", m))
			}
			
			if 200 <= param.StatusCode && param.StatusCode < 300 {
				log.With(fields...).Infof("%d", param.StatusCode)
			} else if 400 <= param.StatusCode && param.StatusCode < 500 {
				log.With(fields...).Warnf("%d", param.StatusCode)
			} else {
				log.With(fields...).Errorf("%d", param.StatusCode)
			}
		}
	}
}
