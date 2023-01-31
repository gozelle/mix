package examples_impl

import (
	"context"
	"fmt"
	example_api "github.com/gozelle/mix/example/api"
	example_types "github.com/gozelle/mix/example/types"
)

var _ example_api.FullAPI = (*ExampleImpl)(nil)

type ExampleImpl struct {
	example_api.FullAPIStub
}

func (e ExampleImpl) Ping(ctx context.Context) (msg example_api.Message, err error) {
	msg.Msg = "Hello from mix!"
	return
}

func (e ExampleImpl) AddAdvance(ctx context.Context, advance *example_types.Advance) error {
	fmt.Println(advance.Name)
	return nil
}

func (e ExampleImpl) AddInt(ctx context.Context, p1 int, p2 int8, p3 int16, p4 int32, p5 int64) (err error) {
	fmt.Println(p1, p2, p3, p4, p5)
	return
}
