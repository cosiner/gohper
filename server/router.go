package server

import (
	"net/url"

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
	// Router is responsible for route manage and match
	Router interface {
		// Init init handlers and filters, parameter function's return value
		// indicate whether continue init next handler
		Init(func(Handler) (still bool), func(Filter) (still bool))
		// Destroy destroy router, also responsible for destroy all handlers and filters
		Destroy()
		// AddFuncHandler add a function handler, method are defined as constant string
		AddFuncHandler(pattern string, method string, handler HandlerFunc) error
		// AddHandler add a handler
		AddHandler(pattern string, handler Handler) error
		// MatchHandler match given url to find a handler, also return match url variables
		MatchHandler(url *url.URL) (handler Handler, urlVars map[string]string)
		// AddFilter add a filter
		AddFilter(pattern string, filter Filter) error
		// MatchFilters match given url to find all matched filters
		MatchFilters(url *url.URL) []Filter
	}

	// router is a url router, the later added handler is first matched,
	// and match only one, filters is matched as the order of added,
	// and more than one can be matched
	router struct {
		handlerRoutes []*handlerRoute
		filterRoutes  []*filterRoute
	}
)

// NewRouter return a new router
func NewRouter() Router {
	return new(router)
}

// Init init router's handlers and filters with given function
func (rt *router) Init(initHandler func(Handler) bool, initFilter func(Filter) bool) {
	for _, r := range rt.handlerRoutes {
		if !initHandler(r.Handler) {
			break
		}
	}
	for _, r := range rt.filterRoutes {
		if !initFilter(r.Filter) {
			break
		}
	}
	return
}

// Destroy destroy router and it's handler routes, filter routes
func (rt *router) Destroy() {
	for _, h := range rt.handlerRoutes {
		h.Destroy()
	}
	for _, f := range rt.filterRoutes {
		f.Destroy()
	}
}

// AddFuncHandler add function handler to router
// the given pattern and new funcFilter will be staged for later same pattern
// function handler with different method
func (rt *router) AddFuncHandler(pattern, method string, handleFunc HandlerFunc) (err error) {
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
func (rt *router) AddHandler(pattern string, handler Handler) (err error) {
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
func (rt *router) MatchHandler(url *url.URL) (handler Handler, urlVars map[string]string) {
	path, routes := url.Path, rt.handlerRoutes
	for i := len(routes) - 1; i >= 0; i-- {
		r := routes[i]
		if vals, match := r.Match(path); match {
			handler, urlVars = r.Handler, vals
			break
		}
	}
	return
}

// AddFilter add filter to router,
// filter can be FilterFunc for FilterFunc is also a filter
// the compiled url matcher will be staged for later added router
// that with same pattern
func (rt *router) AddFilter(pattern string, filter Filter) (err error) {
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

// MatchFilters return matched filters of url path
// the order of filters is same as they are added to router
func (rt *router) MatchFilters(url *url.URL) (filters []Filter) {
	path := url.Path
	for _, r := range rt.filterRoutes {
		if r.MatchOnly(path) {
			filters = append(filters, r.Filter)
		}
	}
	return
}
