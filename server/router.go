package server

import (
	"fmt"
	"net/url"
)

type (
	// Router is responsible for route manage and match
	Router interface {
		// Init init handlers and filters, parameter function's return value
		// indicate whether continue init next handler
		Init(func(Handler) (still bool),
			func(Filter) (still bool),
			func(WebSocketHandler) (still bool))
		// Destroy destroy router, also responsible for destroy all handlers and filters
		Destroy()
		// AddFuncHandler add a function handler, method are defined as constant string
		AddFuncHandler(pattern string, method string, handler HandlerFunc) error
		// AddHandler add a handler
		AddHandler(pattern string, handler Handler) error
		// MatchHandler match given url to find a handler, also return match url variables
		MatchHandler(url *url.URL) (handler Handler, indexer VarIndexer, values []string)
		// AddFuncFilter add function filter
		AddFuncFilter(pattern string, filter FilterFunc) error
		// AddFilter add a filter
		AddFilter(pattern string, filter Filter) error
		// MatchFilters match given url to find all matched filters
		MatchFilters(url *url.URL) []Filter
		// AddWebsocketFuncHandler add a websocket functionhandler
		AddWebsocketFuncHandler(pattern string, handler WebSocketHandlerFunc) error
		// AddWebsocketHandler add a websocket handler
		AddWebsocketHandler(pattern string, handler WebSocketHandler) error
		// MatchWebSocketHandler match given url to find a matched websocket handler
		MatchWebSocketHandler(url *url.URL) (handler WebSocketHandler, indexer VarIndexer, values []string)
	}

	// handlerRoute is route of handler
	handlerRoute struct {
		Matcher
		Handler
	}

	// filterRoute is route of filter
	filterRoute struct {
		Matcher
		Filter
	}

	// websocketRoute is route of websocket handler
	websocketRoute struct {
		Matcher
		WebSocketHandler
	}

	// router is a url router, the later added handler is first matched,
	// and match only one, filters is matched as the order of added,
	// and more than one can be matched
	// filter's pattern can only be literal
	// currently router only support url path match
	router struct {
		handlerRoutes   []*handlerRoute
		filterRoutes    []*filterRoute
		websocketRoutes []*websocketRoute
	}
)

// NewRouter return a new router
func NewRouter() Router {
	return new(router)
}

// Init init router's handlers and filters with given function
func (rt *router) Init(initHandler func(Handler) bool,
	initFilter func(Filter) bool,
	initWebSocketHandler func(WebSocketHandler) bool) {
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
	for _, r := range rt.websocketRoutes {
		if !initWebSocketHandler(r.WebSocketHandler) {
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
	for _, w := range rt.websocketRoutes {
		w.Destroy()
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
	fmt.Println(len(rt.handlerRoutes), "--")
	return
}

// AddHandler add handler to router
// the compiled url matcher will be staged for later added filter
// that with same pattern
func (rt *router) AddHandler(pattern string, handler Handler) (err error) {
	matcher, err := rt.buildMatcher(pattern, false)
	if err == nil {
		rt.handlerRoutes = append(rt.handlerRoutes,
			&handlerRoute{
				Matcher: matcher,
				Handler: handler,
			})
	}
	return
}

// Matchhandler return matched handler and url variables of givel url path
// handler is matched in the reverse order of they are added to router
// it's perform full matched
func (rt *router) MatchHandler(url *url.URL) (handler Handler, indexer VarIndexer, values []string) {
	path, routes := url.Path, rt.handlerRoutes
	for i := len(routes) - 1; i >= 0; i-- {
		r := routes[i]
		if vals, match := r.Match(path); match {
			handler, indexer, values = r.Handler, r.Matcher, vals
			break
		}
	}
	return
}

// AddFuncFilter add function filter
func (rt *router) AddFuncFilter(pattern string, filter FilterFunc) error {
	return rt.AddFilter(pattern, filter)
}

// AddFilter add filter to router,
// filter can be FilterFunc for FilterFunc is also a filter
// the compiled url matcher will be staged for later added router
// that with same pattern
func (rt *router) AddFilter(pattern string, filter Filter) (err error) {
	matcher, err := rt.buildMatcher(pattern, true)
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
// MatchFilters is perform prefix match
func (rt *router) MatchFilters(url *url.URL) []Filter {
	path, routes := url.Path, rt.filterRoutes
	index, l := 0, len(routes)
	filters := make([]Filter, l)
	for _, r := range routes {
		if r.PrefixMatchOnly(path) {
			filters[index] = r.Filter
			index++
		}
	}
	return filters[:index]
}

// AddWebsocketFuncHandler add a websocket function handler
func (rt *router) AddWebsocketFuncHandler(pattern string, handler WebSocketHandlerFunc) error {
	return rt.AddWebsocketHandler(pattern, handler)
}

// AddWebsocketHandler add a websocket handler
func (rt *router) AddWebsocketHandler(pattern string, handler WebSocketHandler) error {
	matcher, err := rt.buildMatcher(pattern, false)
	if err == nil {
		rt.websocketRoutes = append(rt.websocketRoutes,
			&websocketRoute{
				Matcher:          matcher,
				WebSocketHandler: handler,
			})
	}
	return err
}

// MatchWebSocketHandler match given url to find a matched websocket handler
// the match order is reverse of handlers are added to router
// it's perform full matched
func (rt *router) MatchWebSocketHandler(url *url.URL) (handler WebSocketHandler,
	indexer VarIndexer, values []string) {
	path, routes := url.Path, rt.websocketRoutes
	for i := len(routes) - 1; i >= 0; i-- {
		r := routes[i]
		if vals, match := r.Match(path); match {
			handler, indexer, values = r.WebSocketHandler, r.Matcher, vals
			break
		}
	}
	return
}

// buildMatcher build a matcher from pattern, if matcher already exist for same
// pattern in strach, just use the exist one, else create a new one
func (*router) buildMatcher(pattern string, forceLiteral bool) (matcher Matcher, err error) {
	if matcher = strach.routeMatcher(pattern); matcher == nil {
		newMatcherFunc := NewMatcher
		if forceLiteral {
			newMatcherFunc = NewLiteralMatcher
		}
		if matcher, err = newMatcherFunc(pattern); err == nil {
			strach.setRouteMatcher(pattern, matcher)
		}
	} else if forceLiteral && !matcher.IsLiteral() {
		matcher, err = nil, nonLiteralError
	}
	return
}
