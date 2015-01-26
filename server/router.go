package server

import (
	"github.com/cosiner/golib/regexp/urlmatcher"
)

type (
	// handlerRoute is route of handler
	handlerRoute struct {
		*urlmatcher.Matcher
		Handler
	}

	// filterRoute is route of filter
	filterRoute struct {
		*urlmatcher.Matcher
		Filter
	}

	// Router is a url router, the later added handler is first matched,
	// and match only one, filters is matched as the order of added,
	// and more than one can be matched
	Router struct {
		handlerRoutes []*handlerRoute
		filterRoutes  []*filterRoute
	}
)

// NewRouter return a new router
func NewRouter() *Router {
	return new(Router)
}

// initHandler init router's handler with given function
func (rt *Router) initHandler(initFn func(h Handler) bool) {
	for _, r := range rt.handlerRoutes {
		if !initFn(r.Handler) {
			break
		}
	}
	return
}

// initHandler init router's filter with given function
func (rt *Router) initFilter(initFn func(f Filter) bool) {
	for _, r := range rt.filterRoutes {
		if !initFn(r.Filter) {
			break
		}
	}
	return
}

// destroy destroy router and it's handler routes, filter routes
func (rt *Router) destroy() {
	for _, h := range rt.handlerRoutes {
		h.Destroy()
	}
	for _, f := range rt.filterRoutes {
		f.Destroy()
	}
}

// addFuncHandler add function handler to router
// the given pattern and new funcFilter will be staged for later same pattern
// function handler with different method
func (rt *Router) addFuncHandler(pattern, method string, handleFunc HandlerFunc) (err error) {
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

// AddHandler add handler to router
// the compiled url matcher will be staged for later added filter
// that with same pattern
func (rt *Router) AddHandler(pattern string, handler Handler) (err error) {
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

// handler return matched handler and url variables of givel url path
// handler is matched in the reverse order of they are added to router
func (rt *Router) handler(path string) (handler Handler, urlVars map[string]string) {
	routes := rt.handlerRoutes
	for i := len(routes) - 1; i >= 0; i-- {
		r := routes[i]
		if vals, match := r.Match(path); match {
			handler, urlVars = r.Handler, vals
			break
		}
	}
	return
}

// addFuncFilter add function filter to router
// the pattern and funcFilter will be staged for later added filter function
// with same pattern and different filter time
func (rt *Router) addFuncFilter(pattern string, when int, filterFunc FilterFunc) (err error) {
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

// AddFilter add filter to router
// the compiled url matcher will be staged for later added router
// that with same pattern
func (rt *Router) AddFilter(pattern string, filter Filter) (err error) {
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

// filters return matched filters of url path
// the order of filters is same as they are added to router
func (rt *Router) filters(path string) (filters []Filter) {
	for _, r := range rt.filterRoutes {
		if r.MatchOnly(path) {
			filters = append(filters, r.Filter)
		}
	}
	return
}
