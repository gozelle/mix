package test

import (
	"context"
	"time"
)

type TestRpc interface {
	Register(ctx context.Context, username, password string) (token string, err error)
	GetUser(ctx context.Context, id int64) (user *User, err error)
}

type User struct {
	ID        int64
	Name      string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
