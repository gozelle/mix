package sentry

import "fmt"

type Action struct {
	Slug        string
	Routes      []*Route
	Description string
}

type Actions map[string]*Action

func (a Actions) List() (list []*Action, err error) {
	for k, v := range a {
		if v.Slug != "" && v.Slug != k {
			err = fmt.Errorf("action %s not equal slug: %s", k, v.Slug)
			return
		}
		if v.Slug == "" {
			v.Slug = k
		}
		list = append(list, v)
	}
	return
}
