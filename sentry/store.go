package sentry

import (
	"context"
	"github.com/gozelle/mix/sentry/menu"
)

type InitStore interface {
	SysRoutes(ctx context.Context) (routes []*Route, err error)                   // 获取系统中已有的路由
	SysActions(ctx context.Context) (actions []*Action, err error)                // 获取系统中以后的操作
	SysMenu(ctx context.Context) (m *menu.Menu, err error)                        // 获取系统中已有的菜单
	UserRoutes(ctx context.Context) (routes []*Route, err error)                  // 已分配的所有用户（角色）的路由列表
	UserActions(ctx context.Context) (actions []*Action, err error)               // 已分配的所有用户（角色）的操作列表
	UserMenu(ctx context.Context) (m *menu.Menu, err error)                       // 已分配的用户（角色）的菜单列表
	Save(ctx context.Context, routes Routes, actions Actions, m *menu.Menu) error // 保存当前系统权限配置，并将原来的用户权限标记备份
}

type AdminStore interface {
}
