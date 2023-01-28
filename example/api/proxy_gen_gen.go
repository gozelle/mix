package example_api

import (
	"context"
	"io"
	"time"

	example_types "github.com/gozelle/mix/example/types"
	"github.com/shopspring/decimal"
	"golang.org/x/xerrors"
)

var ErrNotSupported = xerrors.New("method not supported")

type AdvanceAPIStruct struct {
	Internal struct {
		AddAdvance func(p0 context.Context, p1 *example_types.Advance) error ``

		Download func(p0 context.Context, p1 string) (io.Reader, error) ``

		GetAdvance func(p0 context.Context, p1 string) (example_types.Advance, error) ``

		Stream func(p0 context.Context) (<-chan any, error) ``
	}
}

type AdvanceAPIStub struct {
}

type BasicAPIStruct struct {
	Internal struct {
		AddBool func(p0 context.Context, p1 bool) error ``

		AddBoolPointer func(p0 context.Context, p1 *bool) error ``

		AddDecimal func(p0 context.Context, p1 decimal.Decimal) error ``

		AddDecimalPinter func(p0 context.Context, p1 *decimal.Decimal) error ``

		AddFloat func(p0 context.Context, p1 float32, p2 float64) error ``

		AddFloatPointer func(p0 context.Context, p1 *float32, p2 *float64) error ``

		AddInt func(p0 context.Context, p1 int, p2 int8, p3 int16, p4 int32, p5 int64) error ``

		AddIntPointer func(p0 context.Context, p1 *int, p2 *int8, p3 *int16, p4 *int32, p5 *int64) error ``

		AddMap func(p0 context.Context, p1 map[string]string) error ``

		AddString func(p0 context.Context, p1 string) error ``

		AddStringPointer func(p0 context.Context, p1 *string) error ``

		AddTime func(p0 context.Context, p1 time.Time) error ``

		AddTimePointer func(p0 context.Context, p1 *time.Time) error ``

		AddUIntPointer func(p0 context.Context, p1 *uint, p2 *uint8, p3 *uint16, p4 *uint32, p5 *int64) error ``

		AddUint func(p0 context.Context, p1 uint, p2 uint8, p3 uint16, p4 uint32, p5 uint64) error ``
	}
}

type BasicAPIStub struct {
}

type FullAPIStruct struct {
	BasicAPIStruct

	AdvanceAPIStruct

	StructAPIStruct

	StructPointerAPIStruct

	Internal struct {
		Ping func(p0 context.Context) (Message, error) ``
	}
}

type FullAPIStub struct {
	BasicAPIStub

	AdvanceAPIStub

	StructAPIStub

	StructPointerAPIStub
}

type StructAPIStruct struct {
	Internal struct {
		AddFull func(p0 context.Context, p1 example_types.Full) error ``

		GetFull func(p0 context.Context) (example_types.Full, error) ``
	}
}

type StructAPIStub struct {
}

type StructPointerAPIStruct struct {
	Internal struct {
		AddFullPointer func(p0 context.Context, p1 *example_types.FullPointer) error ``

		GetFullPointer func(p0 context.Context) (*example_types.FullPointer, error) ``
	}
}

type StructPointerAPIStub struct {
}

func (s *AdvanceAPIStruct) AddAdvance(p0 context.Context, p1 *example_types.Advance) error {
	if s.Internal.AddAdvance == nil {
		return ErrNotSupported
	}
	return s.Internal.AddAdvance(p0, p1)
}

func (s *AdvanceAPIStub) AddAdvance(p0 context.Context, p1 *example_types.Advance) error {
	return ErrNotSupported
}

func (s *AdvanceAPIStruct) Download(p0 context.Context, p1 string) (io.Reader, error) {
	if s.Internal.Download == nil {
		return *new(io.Reader), ErrNotSupported
	}
	return s.Internal.Download(p0, p1)
}

func (s *AdvanceAPIStub) Download(p0 context.Context, p1 string) (io.Reader, error) {
	return *new(io.Reader), ErrNotSupported
}

func (s *AdvanceAPIStruct) GetAdvance(p0 context.Context, p1 string) (example_types.Advance, error) {
	if s.Internal.GetAdvance == nil {
		return *new(example_types.Advance), ErrNotSupported
	}
	return s.Internal.GetAdvance(p0, p1)
}

func (s *AdvanceAPIStub) GetAdvance(p0 context.Context, p1 string) (example_types.Advance, error) {
	return *new(example_types.Advance), ErrNotSupported
}

func (s *AdvanceAPIStruct) Stream(p0 context.Context) (<-chan any, error) {
	if s.Internal.Stream == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.Stream(p0)
}

func (s *AdvanceAPIStub) Stream(p0 context.Context) (<-chan any, error) {
	return nil, ErrNotSupported
}

func (s *BasicAPIStruct) AddBool(p0 context.Context, p1 bool) error {
	if s.Internal.AddBool == nil {
		return ErrNotSupported
	}
	return s.Internal.AddBool(p0, p1)
}

func (s *BasicAPIStub) AddBool(p0 context.Context, p1 bool) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddBoolPointer(p0 context.Context, p1 *bool) error {
	if s.Internal.AddBoolPointer == nil {
		return ErrNotSupported
	}
	return s.Internal.AddBoolPointer(p0, p1)
}

func (s *BasicAPIStub) AddBoolPointer(p0 context.Context, p1 *bool) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddDecimal(p0 context.Context, p1 decimal.Decimal) error {
	if s.Internal.AddDecimal == nil {
		return ErrNotSupported
	}
	return s.Internal.AddDecimal(p0, p1)
}

func (s *BasicAPIStub) AddDecimal(p0 context.Context, p1 decimal.Decimal) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddDecimalPinter(p0 context.Context, p1 *decimal.Decimal) error {
	if s.Internal.AddDecimalPinter == nil {
		return ErrNotSupported
	}
	return s.Internal.AddDecimalPinter(p0, p1)
}

func (s *BasicAPIStub) AddDecimalPinter(p0 context.Context, p1 *decimal.Decimal) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddFloat(p0 context.Context, p1 float32, p2 float64) error {
	if s.Internal.AddFloat == nil {
		return ErrNotSupported
	}
	return s.Internal.AddFloat(p0, p1, p2)
}

func (s *BasicAPIStub) AddFloat(p0 context.Context, p1 float32, p2 float64) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddFloatPointer(p0 context.Context, p1 *float32, p2 *float64) error {
	if s.Internal.AddFloatPointer == nil {
		return ErrNotSupported
	}
	return s.Internal.AddFloatPointer(p0, p1, p2)
}

func (s *BasicAPIStub) AddFloatPointer(p0 context.Context, p1 *float32, p2 *float64) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddInt(p0 context.Context, p1 int, p2 int8, p3 int16, p4 int32, p5 int64) error {
	if s.Internal.AddInt == nil {
		return ErrNotSupported
	}
	return s.Internal.AddInt(p0, p1, p2, p3, p4, p5)
}

func (s *BasicAPIStub) AddInt(p0 context.Context, p1 int, p2 int8, p3 int16, p4 int32, p5 int64) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddIntPointer(p0 context.Context, p1 *int, p2 *int8, p3 *int16, p4 *int32, p5 *int64) error {
	if s.Internal.AddIntPointer == nil {
		return ErrNotSupported
	}
	return s.Internal.AddIntPointer(p0, p1, p2, p3, p4, p5)
}

func (s *BasicAPIStub) AddIntPointer(p0 context.Context, p1 *int, p2 *int8, p3 *int16, p4 *int32, p5 *int64) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddMap(p0 context.Context, p1 map[string]string) error {
	if s.Internal.AddMap == nil {
		return ErrNotSupported
	}
	return s.Internal.AddMap(p0, p1)
}

func (s *BasicAPIStub) AddMap(p0 context.Context, p1 map[string]string) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddString(p0 context.Context, p1 string) error {
	if s.Internal.AddString == nil {
		return ErrNotSupported
	}
	return s.Internal.AddString(p0, p1)
}

func (s *BasicAPIStub) AddString(p0 context.Context, p1 string) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddStringPointer(p0 context.Context, p1 *string) error {
	if s.Internal.AddStringPointer == nil {
		return ErrNotSupported
	}
	return s.Internal.AddStringPointer(p0, p1)
}

func (s *BasicAPIStub) AddStringPointer(p0 context.Context, p1 *string) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddTime(p0 context.Context, p1 time.Time) error {
	if s.Internal.AddTime == nil {
		return ErrNotSupported
	}
	return s.Internal.AddTime(p0, p1)
}

func (s *BasicAPIStub) AddTime(p0 context.Context, p1 time.Time) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddTimePointer(p0 context.Context, p1 *time.Time) error {
	if s.Internal.AddTimePointer == nil {
		return ErrNotSupported
	}
	return s.Internal.AddTimePointer(p0, p1)
}

func (s *BasicAPIStub) AddTimePointer(p0 context.Context, p1 *time.Time) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddUIntPointer(p0 context.Context, p1 *uint, p2 *uint8, p3 *uint16, p4 *uint32, p5 *int64) error {
	if s.Internal.AddUIntPointer == nil {
		return ErrNotSupported
	}
	return s.Internal.AddUIntPointer(p0, p1, p2, p3, p4, p5)
}

func (s *BasicAPIStub) AddUIntPointer(p0 context.Context, p1 *uint, p2 *uint8, p3 *uint16, p4 *uint32, p5 *int64) error {
	return ErrNotSupported
}

func (s *BasicAPIStruct) AddUint(p0 context.Context, p1 uint, p2 uint8, p3 uint16, p4 uint32, p5 uint64) error {
	if s.Internal.AddUint == nil {
		return ErrNotSupported
	}
	return s.Internal.AddUint(p0, p1, p2, p3, p4, p5)
}

func (s *BasicAPIStub) AddUint(p0 context.Context, p1 uint, p2 uint8, p3 uint16, p4 uint32, p5 uint64) error {
	return ErrNotSupported
}

func (s *FullAPIStruct) Ping(p0 context.Context) (Message, error) {
	if s.Internal.Ping == nil {
		return *new(Message), ErrNotSupported
	}
	return s.Internal.Ping(p0)
}

func (s *FullAPIStub) Ping(p0 context.Context) (Message, error) {
	return *new(Message), ErrNotSupported
}

func (s *StructAPIStruct) AddFull(p0 context.Context, p1 example_types.Full) error {
	if s.Internal.AddFull == nil {
		return ErrNotSupported
	}
	return s.Internal.AddFull(p0, p1)
}

func (s *StructAPIStub) AddFull(p0 context.Context, p1 example_types.Full) error {
	return ErrNotSupported
}

func (s *StructAPIStruct) GetFull(p0 context.Context) (example_types.Full, error) {
	if s.Internal.GetFull == nil {
		return *new(example_types.Full), ErrNotSupported
	}
	return s.Internal.GetFull(p0)
}

func (s *StructAPIStub) GetFull(p0 context.Context) (example_types.Full, error) {
	return *new(example_types.Full), ErrNotSupported
}

func (s *StructPointerAPIStruct) AddFullPointer(p0 context.Context, p1 *example_types.FullPointer) error {
	if s.Internal.AddFullPointer == nil {
		return ErrNotSupported
	}
	return s.Internal.AddFullPointer(p0, p1)
}

func (s *StructPointerAPIStub) AddFullPointer(p0 context.Context, p1 *example_types.FullPointer) error {
	return ErrNotSupported
}

func (s *StructPointerAPIStruct) GetFullPointer(p0 context.Context) (*example_types.FullPointer, error) {
	if s.Internal.GetFullPointer == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.GetFullPointer(p0)
}

func (s *StructPointerAPIStub) GetFullPointer(p0 context.Context) (*example_types.FullPointer, error) {
	return nil, ErrNotSupported
}

var _ AdvanceAPI = new(AdvanceAPIStruct)
var _ BasicAPI = new(BasicAPIStruct)
var _ FullAPI = new(FullAPIStruct)
var _ StructAPI = new(StructAPIStruct)
var _ StructPointerAPI = new(StructPointerAPIStruct)
