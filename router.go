package server

import (
	"net/url"
	"strings"

	. "github.com/cosiner/golib/errors"
)

type (
	pathNode struct {
		Var   string
		Index int
	}

	route struct {
		Vars             []pathNode
		Handler          Handler
		Filters          []Filter
		WebSocketHandler WebSocketHandler
	}

	router struct {
		sub   map[string]*router
		route *route
	}
)

func (r *route) Values(path []string) (values []string) {
	vars := r.Vars
	if l := len(vars); l != 0 {
		values = make([]string, len(vars))
		for i := 0; i < l; i++ {
			values[i] = path[vars[i].Index]
		}
	}
	return
}

func (r *route) ValuesMap(path []string) (values map[string]string) {
	vars := r.Vars
	if l := len(vars); l != 0 {
		values = make(map[string]string, len(vars))
		for i := 0; i < l; i++ {
			node := vars[i]
			values[node.Var] = path[node.Index]
		}
	}
	return
}

func (r *route) ValuesScan(path []string, scanVars ...*string) {
	vars := r.Vars
	if l, l2 := len(vars), len(scanVars); l != 0 && l2 != 0 {
		for i := 0; i < l && i < l2; i++ {
			(*scanVars[i]) = path[vars[i].Index]
		}
	}
	return
}

func (rt *router) Init(initHandler func(Handler) bool,
	initFilter func(Filter) bool,
	initWebSocketHandler func(WebSocketHandler) bool) {
	rt.init(initHandler, initFilter, initWebSocketHandler)
}

func (rt *router) init(initHandler func(Handler) bool,
	initFilter func(Filter) bool,
	initWebSocketHandler func(WebSocketHandler) bool) bool {
	route, continu := rt.route, true
	if route == nil {
		return true
	}
	if route.Handler != nil {
		continu = initHandler(route.Handler)
	}
	if route.WebSocketHandler != nil {
		continu = initWebSocketHandler(route.WebSocketHandler)
	}
	if continu {
		for _, f := range route.Filters {
			if continu = initFilter(f); !continu {
				break
			}
		}
	}
	if continu {
		for _, rt := range rt.sub {
			if continu = rt.init(initHandler,
				initFilter,
				initWebSocketHandler); !continu {
				break
			}
		}
	}
	return continu
}

// matchMultiRoute match multiple routes
func (rt *router) matchMultiRoutes(url *url.URL) (path []string, routes []*route, match bool) {
	path = strings.Split(url.Path, "/")
	routes = make([]*route, 0, len(path)+1)
	depth, match := rt.matchMulti(path, routes, 0)
	routes = routes[:depth+1]
	return
}

// matchMulti match multi route, result newDepth means the max matched depth, match means
// whether whole matched
func (rt *router) matchMulti(path []string, routes []*route, depth int) (newDepth int, match bool) {
	routes[0] = rt.route
	newDepth, match = depth, true
	if len(path) != 0 {
		newPath, newRoutes := path[1:], routes[1:]
		match = false
		if newRt := rt.sub(path[0]); newRt != nil {
			newDepth, match = newRt.matchMulti(newPath, newRoutes, depth+1)
		}
		if !match {
			if newRt := rt.sub("*"); newRt != nil {
				newDepth, match = newRt.matchMulti(newPath, newRoutes, depth+1)
			}
		}
	}
	return
}

// matchSingleRoute will only match the longest path, and first match
// non-variabled route
// example: pattern is 1:/user/:op/123 2:/user/delete/234,
// for url path /user/delete/123, pattern 1 will be matched
func (rt *router) matchSingleRoute(url *url.URL) (path []string, route *route) {
	return rt.matchSingle(strings.Split(url.Path, "/")[1:])
}

func (rt *router) matchSingle(path []string) (r *route) {
	if len(path) == 0 {
		r = rt.route
	} else {
		newPath := path[1:]
		if newRt := rt.sub(path[0]); newRt != nil {
			r = newRt.matchSingle(newPath)
		}
		if r == nil {
			if newRt := rt.sub("*"); newRt != nil {
				r = newRt.matchSingle(newPath)
			}
		}
	}
	return
}

// AddFuncHandler add function handler
func (rt *router) AddFuncHandler(pattern, method string, handler HandlerFunc) error {
	fHandler := strach.funcHandler(pattern)
	if fHandler == nil {
		fHandler = new(funcHandler)
		strach.setFuncHandler(pattern, fHandler)
	}
	err := fHandler.setMethod(method, handler)
	if err == nil {
		err = rt.AddHandler(pattern, fHandler)
	}
	return err
}

// AddHandler add a pattern handler
// NOTE: pattern with same path (path of pattern is that all variables are
// replaced with "*") regardless of it's variable names will be treat as a same
// pattern, and it's not allowed
// example /user/:op/:id and /user/:operate/:ident is two same pattern with path
// /user/*/*, websocket handler should also comply with this rule
func (rt *router) AddHandler(pattern string, handler Handler) (err error) {
	return rt.addPattern(pattern, func(r *route) (err error) {
		if r.Handler != nil {
			err = Err("Handler for given pattern already exist")
		} else {
			r.Handler = handler
		}
		return
	})
}

// AddFuncWebSocketHandler add WebSocket function handler
func (rt *router) AddFuncWebSocketHandler(pattern string, handler WebSocketHandlerFunc) error {
	return rt.AddWebSocketHandler(pattern, handler)
}

// AddWebSockethandler add Websocket Handler
func (rt *router) AddWebSocketHandler(pattern string, handler WebSocketHandler) (err error) {
	return rt.addPattern(pattern, func(r *route) (err error) {
		if r.WebSocketHandler != nil {
			err = Err("WebsocketHandler for given pattern already exist")
		} else {
			r.WebSocketHandler = handler
		}
		return
	})
}

// addFuncFilter add a function filter
func (rt *router) AddFuncFilter(pattern string, filter FilterFunc) error {
	return rt.AddFilter(pattern, filter)
}

// AddFilter add filter to router with given pattern
func (rt *router) AddFilter(pattern string, filter Filter) error {
	return rt.addPattern(pattern, func(r *route) error {
		r.Filters = append(r.Filters, filter)
		return nil
	})
}

// addPattern add a pattern to router, if pattern is wrong format
// or there is pattern with same url path exist, error returned
// else return result of parameter op
func (rt *router) addPattern(pattern string, op func(*route) error) error {
	r := strach.route(pattern)
	if r == nil {
		path, vars, e := rt.compile(pattern)
		if e == nil {
			r = &route{
				Vars: vars,
			}
			if e = rt.addRoute(path, r); e == nil {
				strach.setRoute(pattern, r)
			}
		}
		if e != nil {
			return e
		}
	}
	return op(r)
}

// addRoute add a route with given path to router
// path is url path, all variable are replaced with "*"
// if route with same path already exist, error returned
func (rt *router) addRoute(path []string, route *route) (err error) {
	index := 0
	for index != len(path) {
		p := path[index]
		if newRt := rt.sub[p]; newRt == nil {
			newRt = &router{
				route: nil,
				sub:   nil,
			}
			if rt.sub == nil {
				rt.sub = make(map[string]*router)
			}
			rt.sub[p] = newRt
		}
		index++
	}
	if rt.route != nil {
		err = Err("route already exist")
	} else {
		rt.route = route
	}
	return
}

// errWrongFormat means wrong pattern format
var errWrongFormat = Err("Wrong format")

// compile compile given pattern, named variable is represent as :name
// example /:user/:id, and all variable are replaced with "*"
// only the vars record where is the variable and what is the variable name
func (*router) compile(pattern string) (path []string, vars []pathNode, err error) {
	var patternPath []string
	if len(pattern) == 0 || pattern[0] != '/' {
		goto ERROR
	}
	patternPath = strings.Split(pattern, "/")[1:]
	path = make([]string, len(patternPath))
	for i, p := range patternPath {
		if len(p) == 0 {
			goto ERROR
		}
		if p[0] == ':' {
			if p = p[1:]; len(p) == 0 {
				goto ERROR
			}
			vars = append(vars, pathNode{Var: p, Index: i})
			p = "*"
		}
		path[i] = p
	}
	goto END
ERROR:
	path, vars, err = nil, nil, errWrongFormat
END:
	return
}
