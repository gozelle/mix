package example

import (
	"context"

	"github.com/gozelle/mix/example/example_types"
	"github.com/shopspring/decimal"
	"golang.org/x/xerrors"
)

var ErrNotSupported = xerrors.New("method not supported")

type TestRpcStruct struct {
	Internal struct {
		GetInfo func(p0 context.Context, p1 Info) error ``

		GetOrderByID func(p0 context.Context, p1 ID) ([]*ExampleUser, error) ``

		GetOrderPrice func(p0 context.Context, p1 int64) (Fil, error) ``

		GetOrderPrice2 func(p0 context.Context, p1 int64) (decimal.Decimal, error) ``

		GetUser func(p0 context.Context, p1 int64) (*ExampleUser, error) ``

		Register func(p0 context.Context, p1 string, p2 string) (string, error) ``

		RegisterByEmail func(p0 context.Context, p1 RegisterByEmailRequest) error ``

		SaveApple func(p0 context.Context, p1 example_types.Apple) error ``
	}
}

type TestRpcStub struct {
}

func (s *TestRpcStruct) GetInfo(p0 context.Context, p1 Info) error {
	if s.Internal.GetInfo == nil {
		return ErrNotSupported
	}
	return s.Internal.GetInfo(p0, p1)
}

func (s *TestRpcStub) GetInfo(p0 context.Context, p1 Info) error {
	return ErrNotSupported
}

func (s *TestRpcStruct) GetOrderByID(p0 context.Context, p1 ID) ([]*ExampleUser, error) {
	if s.Internal.GetOrderByID == nil {
		return *new([]*ExampleUser), ErrNotSupported
	}
	return s.Internal.GetOrderByID(p0, p1)
}

func (s *TestRpcStub) GetOrderByID(p0 context.Context, p1 ID) ([]*ExampleUser, error) {
	return *new([]*ExampleUser), ErrNotSupported
}

func (s *TestRpcStruct) GetOrderPrice(p0 context.Context, p1 int64) (Fil, error) {
	if s.Internal.GetOrderPrice == nil {
		return *new(Fil), ErrNotSupported
	}
	return s.Internal.GetOrderPrice(p0, p1)
}

func (s *TestRpcStub) GetOrderPrice(p0 context.Context, p1 int64) (Fil, error) {
	return *new(Fil), ErrNotSupported
}

func (s *TestRpcStruct) GetOrderPrice2(p0 context.Context, p1 int64) (decimal.Decimal, error) {
	if s.Internal.GetOrderPrice2 == nil {
		return *new(decimal.Decimal), ErrNotSupported
	}
	return s.Internal.GetOrderPrice2(p0, p1)
}

func (s *TestRpcStub) GetOrderPrice2(p0 context.Context, p1 int64) (decimal.Decimal, error) {
	return *new(decimal.Decimal), ErrNotSupported
}

func (s *TestRpcStruct) GetUser(p0 context.Context, p1 int64) (*ExampleUser, error) {
	if s.Internal.GetUser == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.GetUser(p0, p1)
}

func (s *TestRpcStub) GetUser(p0 context.Context, p1 int64) (*ExampleUser, error) {
	return nil, ErrNotSupported
}

func (s *TestRpcStruct) Register(p0 context.Context, p1 string, p2 string) (string, error) {
	if s.Internal.Register == nil {
		return "", ErrNotSupported
	}
	return s.Internal.Register(p0, p1, p2)
}

func (s *TestRpcStub) Register(p0 context.Context, p1 string, p2 string) (string, error) {
	return "", ErrNotSupported
}

func (s *TestRpcStruct) RegisterByEmail(p0 context.Context, p1 RegisterByEmailRequest) error {
	if s.Internal.RegisterByEmail == nil {
		return ErrNotSupported
	}
	return s.Internal.RegisterByEmail(p0, p1)
}

func (s *TestRpcStub) RegisterByEmail(p0 context.Context, p1 RegisterByEmailRequest) error {
	return ErrNotSupported
}

func (s *TestRpcStruct) SaveApple(p0 context.Context, p1 example_types.Apple) error {
	if s.Internal.SaveApple == nil {
		return ErrNotSupported
	}
	return s.Internal.SaveApple(p0, p1)
}

func (s *TestRpcStub) SaveApple(p0 context.Context, p1 example_types.Apple) error {
	return ErrNotSupported
}

var _ TestRpc = new(TestRpcStruct)
