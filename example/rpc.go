package example

import (
	"context"
	"github.com/gozelle/mix/example/example_types"
	"github.com/shopspring/decimal"
	time2 "time"
)

type TestRpc interface {
	Register(ctx context.Context, username, password string) (token string, err error) //
	GetUser(ctx context.Context, id int64) (user *ExampleUser, err error)
	GetOrderByID(ctx context.Context, id ID) (user []*ExampleUser, err error)
	GetOrderPrice(ctx context.Context, id int64) (price Fil, err error)
	GetOrderPrice2(ctx context.Context, id int64) (price decimal.Decimal, err error)
	SaveApple(ctx context.Context, apple example_types.Apple) (err error)
	RegisterByEmail(ctx context.Context, r RegisterByEmailRequest) (err error)
}

type ID int64

type Fil decimal.Decimal

type RegisterByEmailRequest struct {
	Email string
}

type ExampleUser struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Price     Fil
	Active    bool
	Tags      []string `json:"tags"`
	CreatedAt time2.Time
	UpdatedAt time2.Time
}

type NoUsed struct {
}
