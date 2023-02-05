package sentry

import "fmt"

type Route struct {
	Method      string
	HttpMethod  string
	HttpPath    string
	Description string
}

type Routes map[string]*Route

func (a Routes) List() (list []*Route, err error) {
	for k, v := range a {
		if v.Method != "" && v.Method != k {
			err = fmt.Errorf("route %s not equal method: %s", k, v.Method)
			return
		}
		if v.Method == "" {
			v.Method = k
		}
		list = append(list, v)
	}
	return
}
