package server

import (
	"github.com/cosiner/golib/regexp/urlmatcher"
)

type handlerRoute struct {
	*urlmatcher.Matcher
	Handler
}

type filterRoute struct {
	*urlmatcher.Matcher
	Filter
}

type Router struct {
	handlerRoutes []*handlerRoute
	filterRoutes  []*filterRoute
}

func NewRouter() Router {
	return Router{}
}

func (rt Router) initHandler(initFn func(h Handler) bool) {
	for _, r := range rt.handlerRoutes {
		if !initFn(r.Handler) {
			break
		}
	}
	return
}

func (rt Router) initFilter(initFn func(f Filter) bool) {
	for _, r := range rt.filterRoutes {
		if !initFn(r.Filter) {
			break
		}
	}
	return
}

func (rt Router) addFuncHandler(pattern, method string, handleFunc HandlerFunc) (err error) {
	if fHandler := strach.funcHandler(pattern); fHandler == nil {
		fHandler = new(funcHandler)
		if err = fHandler.setMethod(method, handleFunc); err == nil {
			if err = rt.AddHandler(pattern, fHandler); err == nil {
				strach.setFuncHandler(pattern, fHandler)
			}
		}
	} else {
		err = fHandler.setMethod(method, handleFunc)
	}
	return
}

func (rt Router) AddHandler(pattern string, handler Handler) (err error) {
	matcher := strach.routeMatcher(pattern)
	if matcher == nil {
		if matcher, err = urlmatcher.Compile(pattern); err == nil {
			strach.setRouteMatcher(pattern, matcher)
		}
	}
	if err == nil {
		rt.handlerRoutes = append(rt.handlerRoutes,
			&handlerRoute{
				Matcher: matcher,
				Handler: handler,
			})
	}
	return
}
func (rt Router) handler(path string) (handler Handler, urlStories map[string]string) {
	routes := rt.handlerRoutes
	for i := len(routes) - 1; i >= 0; i-- {
		r := routes[i]
		if vals, match := r.Match(path); match {
			handler, urlStories = r.Handler, vals
			break
		}
	}
	return
}
func (rt Router) addFuncFilter(pattern string, when int, filterFunc FilterFunc) (err error) {
	if fFilter := strach.funcFilter(pattern); fFilter == nil {
		fFilter = new(funcFilter)
		if err = fFilter.setFilterFunc(when, filterFunc); err == nil {
			if err = rt.AddFilter(pattern, fFilter); err == nil {
				strach.setFuncFilter(pattern, fFilter)
			}
		}
	} else {
		err = fFilter.setFilterFunc(when, filterFunc)
	}
	return
}

func (rt Router) AddFilter(pattern string, filter Filter) (err error) {
	matcher := strach.routeMatcher(pattern)
	if matcher == nil {
		if matcher, err = urlmatcher.Compile(pattern); err == nil {
			strach.setRouteMatcher(pattern, matcher)
		}
	}
	if err == nil {
		rt.filterRoutes = append(rt.filterRoutes,
			&filterRoute{
				Matcher: matcher,
				Filter:  filter,
			})
	}
	return
}

func (rt Router) filters(path string) (filters []Filter) {
	for _, r := range rt.filterRoutes {
		if r.MatchOnly(path) {
			filters = append(filters, r.Filter)
		}
	}
	return
}
