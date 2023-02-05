package sentry

import (
	"fmt"
	"github.com/gozelle/mix/sentry/menu"
	"net/http"
)

type Config struct {
	store   InitStore
	routes  Routes
	actions Actions
	menu    *menu.Menu
}

type Option func(c *Config)

func WithRoutes(routes Routes) Option {
	return func(c *Config) {
		c.routes = routes
	}
}

func WithStore(s InitStore) Option {
	return func(c *Config) {
		c.store = s
	}
}

func WithActions(actions Actions) Option {
	return func(c *Config) {
		c.actions = actions
	}
}

func WithMenu(m *menu.Menu) Option {
	return func(c *Config) {
		c.menu = m
	}
}

type Sentry struct {
}

// Auth 验证用户权限
func (s *Sentry) Auth(httpMethod, httpPath string, userRoutes []*Route) bool {
	m := map[string]map[string]bool{}
	for _, v := range userRoutes {
		if m[v.HttpMethod] == nil {
			m[v.HttpMethod] = map[string]bool{}
		}
		m[v.HttpMethod][v.HttpPath] = true
	}
	v, ok := m[httpMethod]
	if !ok {
		return false
	}
	_, ok = v[httpPath]

	return ok
}

// AuthRequest 通过请求验证用户权限
func (s *Sentry) AuthRequest(req *http.Request, userRoutes []*Route) bool {
	return s.Auth(req.Method, req.URL.Path, userRoutes)
}

// Register 注册系统权限配置
func (s *Sentry) Register(options ...Option) (err error) {

	c := &Config{}
	for _, v := range options {
		v(c)
	}
	if c.store == nil {
		err = fmt.Errorf("sentry store is nil")
		return
	}
	return
}

func (s *Sentry) registerRoutes(store InitStore, routes Routes) (err error) {
	return
}

func (s *Sentry) registerActions(store InitStore, actions Actions) (err error) {
	return
}

func (s *Sentry) registerMenu(store InitStore, menu *menu.Menu) (err error) {
	return
}
