# Mix

## 快速上手

```go
package main

import (
	"context"
	"github.com/gozelle/mix"
)

type IHandler interface {
	Ping(ctx context.Context, message string) (reply string, err error)
}

var _ IHandler = (*Handler)(nil)

type Handler struct {
}

func (h Handler) Ping(ctx context.Context, message string) (reply string, err error) {
	reply = message
	return
}

func main() {
	h := &Handler{}
	s := mix.NewServer()
	
	mix.RegisterRPC(s.Group("/rpc/v1"), "foo", h)
	mix.RegisterAPI(s.Group("/api/v1"), "foo", h)
	
	s.Run(":1332")
}
```

请求示例

API:
```
curl --location --request POST '127.0.0.1:11111/api/v1/Ping' \
--header 'Content-Type: text/plain' \
--data-raw '"Hello world!"'
```

RPC: 
```
curl --location --request POST '127.0.0.1:1332/rpc/v1/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "method": "foo.Ping",
    "params": ["Hello world!"]
}'
```
