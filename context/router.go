package context

import (
	"github.com/cosiner/golib/regexp/urlmatcher"
)

type Router struct {
	routes []*struct {
		*urlmatcher.Matcher
		Handler
	}
}

func (rt *Router) InitHandler(initFn func(h Handler) bool) {
	for _, r := range rt.routes {
		if !initFn(r.Handler) {
			break
		}
	}
	return
}

func (rt *Router) AddFuncRoute(pattern string, method string, handleFunc HandlerFunc) {
	fHandler := strachFuncHandler(pattern)
	if fHandler == nil {
		fHandler = newFuncHandler()
		strachAddFuncHandler(pattern, fHandler)
		rt.AddRoute(pattern, fHandler)
	}
	switch method {
	case GET:
		fHandler.Get = handleFunc
	case POST:
		fHandler.Post = handleFunc
	case PUT:
		fHandler.Put = handleFunc
	case DELETE:
		fHandler.Delete = handleFunc
	}
}

func (rt *Router) addRoute(pattern string, handler Handler) (err error) {
	var matcher *urlmatcher.Matcher
	if matcher, err = urlmatcher.Compile(pattern); err == nil {
		rt.routes = append(rt.routes,
			&Route{
				Matcher: matcher,
				Handler: handler,
			})
	}
	return
}

func (rt *Router) handler(url string) (handler Handler, urlStories map[string]string) {
	routes := rt.routes
	for i := len(routes) - 1; i >= 0; i-- {
		r := routes[i]
		if vals, match := r.Match(url); match {
			handler, urlStories = r.Handler, vals
			break
		}
	}
	return
}
