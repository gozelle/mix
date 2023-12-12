package tests

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/gozelle/mix"
	"github.com/gozelle/mix/client"
	example_api "github.com/gozelle/mix/example/api"
	examples_impl "github.com/gozelle/mix/example/impl"
	"github.com/gozelle/testify/require"
)

func freePort() (port uint64, err error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			err = fmt.Errorf("failed to listen on a port: %v", err)
			return
		}
	}
	err = l.Close()
	if err != nil {
		return
	}
	_, s, ok := strings.Cut(l.Addr().String(), ":")
	if !ok {
		err = fmt.Errorf("port not found")
		return
	}
	port, err = strconv.ParseUint(s, 10, 64)
	if err != nil {
		err = fmt.Errorf("parse port error: %s", err)
		return
	}
	return
}

func TestServer(t *testing.T) {

	h := &examples_impl.ExampleImpl{}
	s := mix.NewServer()

	s.RegisterRPC(s.Group("/rpc/v1"), "example", h)
	s.RegisterAPI(s.Group("/api/v1"), "example", h)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, l.Close())
	}()

	go func() {
		require.NoError(t, s.RunListener(l))
	}()
	c := &example_api.FullAPIStruct{}
	closer, err := client.NewClient(
		fmt.Sprintf("http://%s/rpc/v1", l.Addr().String()), "example",
		client.WithOut(
			&c.Internal,
			&c.BasicAPIStruct.Internal,
			&c.AdvanceAPIStruct.Internal,
			&c.StructAPIStruct.Internal,
			&c.StructPointerAPIStruct.Internal,
		),
	)
	require.NoError(t, err)
	defer func() {
		closer()
	}()
	ctx := context.Background()
	msg, err := c.Ping(ctx)
	require.NoError(t, err)
	require.Equal(t, "Hello from mix!", msg.Msg)

	r := resty.New()
	api := fmt.Sprintf("http://%s/api/v1/example", l.Addr().String())
	_, err = r.R().SetBody(map[string]interface{}{"Name": "test"}).Post(api + "/AddAdvance")
	require.NoError(t, err)

	_, err = r.R().SetBody([]int{1, 2, 3, 4, 5}).Post(api + "/AddInt")
	require.NoError(t, err)

}
