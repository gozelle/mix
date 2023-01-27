package client

import (
	"context"
	"fmt"
	"github.com/gozelle/jsonrpc"
)

type Config struct {
	addr      string
	namespace string
	user      string
	password  string
	token     string
	out       []interface{}
}

type Option func(*Config)

func WithBasicAuth(user, password string) Option {
	return func(config *Config) {
		config.user = user
		config.password = password
	}
}

func WithToken(token string) Option {
	return func(config *Config) {
		config.token = token
	}
}

func WithOut(pointer ...interface{}) Option {
	return func(config *Config) {
		config.out = append(config.out, pointer...)
	}
}

func NewClient(addr, namespace string, opts ...Option) (jsonrpc.ClientCloser, error) {
	conf := &Config{
		addr:      addr,
		namespace: namespace,
	}
	for _, v := range opts {
		v(conf)
	}

	if len(conf.out) == 0 {
		return nil, fmt.Errorf("please assign out use WithOut function")
	}

	return jsonrpc.NewMergeClient(context.Background(), conf.addr, conf.namespace, conf.out, nil)
}
