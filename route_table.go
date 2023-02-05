package mix

import (
	"net/http"
	"sync"
)

type handler = string
type description = string
type httpMethod = string
type httpPath = string

type route struct {
	method string
	path   string
}

func NewRouteTable(prefix string) *RouteTable {
	return &RouteTable{prefix: prefix}
}

type RouteTable struct {
	lock   sync.Mutex
	prefix string
	routes map[httpMethod]map[httpPath]bool
	skips  map[httpMethod]map[httpPath]bool
}

type RoutePermission struct {
	Permission  string
	Description string
}

func (r *RouteTable) AddSkip(method httpMethod, handler handler) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()

	return
}

func (r *RouteTable) AddRoute(m httpMethod, h handler, d description) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()

	return
}

func (r *RouteTable) Has(m httpMethod, p httpPath) (ok bool) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	if r.routes == nil {
		return
	}

	v, ok := r.routes[method]
	if !ok {
		return
	}

	vv, ok := v[p]
	if !ok {
		return
	}

	return
}

func (r *RouteTable) Diff(req *http.Request) (desc string, ok bool) {
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()

	return
}
