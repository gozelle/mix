package sentry

import (
	"github.com/gozelle/mix/sentry/menu"
)

type Role struct {
	Slug    string
	Name    string
	System  bool
	Actions []*Action
	Routes  []*Route
	Menu    []*menu.Menu
}

type Roles struct {
	roles []*Role
}

//func (r *Roles) Routes() (list []*Route) {
//	m := map[string]*Route{}
//	for _, v := range r.roles {
//		for _, vv := range v.Routes {
//			key := fmt.Sprintf("%s@%s", vv.HttpMethod, vv.HttpPath)
//			m[key] = vv
//		}
//	}
//	for _, v := range m {
//		list = append(list, v)
//	}
//	return
//}
