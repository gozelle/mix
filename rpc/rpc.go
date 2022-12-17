package rpc

import (
	"context"
	"github.com/shopspring/decimal"
	"time"
)

type TestRpc interface {
	Register(ctx context.Context, username, password string) (token string, err error)
	GetUser(ctx context.Context, id int64) (user *User, err error)
}

type ID int64

type Fil decimal.Decimal

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
