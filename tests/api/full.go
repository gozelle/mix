package example_api

import (
	"context"
	example_types "github.com/gozelle/mix/example/types"
	"github.com/shopspring/decimal"
	"io"
	"time"
)

type Message struct {
	Msg string // TODO test remove
}

type FullAPI interface {
	Ping(ctx context.Context) (msg Message, err error)
	BasicAPI
	AdvanceAPI
	StructAPI
	StructPointerAPI
}

type StructAPI interface {
	AddFull(ctx context.Context, full example_types.Full) error
	GetFull(ctx context.Context) (full example_types.Full, err error)
}

type StructPointerAPI interface {
	AddFullPointer(ctx context.Context, full *example_types.FullPointer) error
	GetFullPointer(ctx context.Context) (full *example_types.FullPointer, err error)
}

type BasicAPI interface {
	AddInt(ctx context.Context, p1 int, p2 int8, p3 int16, p4 int32, p5 int64) error
	AddUint(ctx context.Context, p1 uint, p2 uint8, p3 uint16, p4 uint32, p5 uint64) error
	AddString(ctx context.Context, p1 string) error
	AddFloat(ctx context.Context, p1 float32, p2 float64) error
	AddBool(ctx context.Context, p1 bool) error
	AddTime(ctx context.Context, p1 time.Time) error
	AddDecimal(ctx context.Context, p1 decimal.Decimal) error
	AddMap(ctx context.Context, m1 map[string]string) error
	
	AddIntArray(ctx context.Context, p1 []int, p2 []int8, p3 []int16, p4 []int32, p5 []int64) error
	
	GetInt(ctx context.Context) (p1 int, err error)
	GetUint(ctx context.Context) (p1 uint, err error)
	GetString(ctx context.Context) (p1 string, err error)
	GetFloat(ctx context.Context) (p1 float32, err error)
	GetBool(ctx context.Context) (p1 bool, err error)
	GetTime(ctx context.Context) (p1 time.Time, err error)
	GetDecimal(ctx context.Context) (p1 decimal.Decimal, err error)
	GetMap(ctx context.Context) (m1 map[string]string, err error)
	
	GetIntArray(ctx context.Context) (int []int, err error)
	
	//AddIntPointer(context.Context, *int, *int8, *int16, *int32, *int64) error
	//AddUIntPointer(context.Context, *uint, *uint8, *uint16, *uint32, *int64) error
	//AddStringPointer(context.Context, *string) error
	//AddFloatPointer(context.Context, *float32, *float64) error
	//AddBoolPointer(context.Context, *bool) error
	//AddTimePointer(context.Context, *time.Time) error
	//AddDecimalPinter(context.Context, *decimal.Decimal) error
	
	//GetIntPointer(context.Context, *int, *int8, *int16, *int32, *int64) (*int, *int8, *int16, *int32, *int64, error)
	//GetUIntPointer(context.Context) (*uint, *uint8, *uint16, *uint32, *int64, error)
	//GetStringPointer(context.Context) (*string, error)
	//GetFloatPointer(context.Context) (*float32, *float64, error)
	//GetBoolPointer(context.Context) (*bool, error)
	//GetTimePointer(context.Context) (*time.Time, error)
	//GetDecimalPinter(context.Context) (*decimal.Decimal, error)
}

type AdvanceAPI interface {
	Download(ctx context.Context, filename string) (io.Reader, error)
	Stream(ctx context.Context) (<-chan any, error)
	AddAdvance(ctx context.Context, advance *example_types.Advance) error
	GetAdvance(ctx context.Context, name string) (example_types.Advance, error)
}
