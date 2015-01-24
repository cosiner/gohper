package context

import (
	"github.com/cosiner/golib/regexp/urlmatcher"
)

type route struct {
	*urlmatcher.Matcher
	Handler
}
type Router []*route

func (rt *Router) InitHandler(initFn func(h Handler) bool) {
	for _, r := range *rt {
		if !initFn(r.Handler) {
			break
		}
	}
	return
}

func (rt *Router) AddFuncRoute(pattern string, method string, handleFunc HandlerFunc) (err error) {
	fHandler := strachFuncHandler(pattern)
	if fHandler == nil {
		fHandler = newFuncHandler()
		if err = fHandler.setMethod(method, handleFunc); err == nil {
			if err = rt.AddRoute(pattern, fHandler); err == nil {
				strachAddFuncHandler(pattern, fHandler)
			}
		}
	} else {
		err = fHandler.setMethod(method, handleFunc)
	}
	return
}

func (rt *Router) AddRoute(pattern string, handler Handler) (err error) {
	var matcher *urlmatcher.Matcher
	if matcher, err = urlmatcher.Compile(pattern); err == nil {
		*rt = append(*rt,
			&route{
				Matcher: matcher,
				Handler: handler,
			})
	}
	return
}

func (rt *Router) Handler(url string) (handler Handler, urlStories map[string]string) {
	routes := *rt
	for i := len(routes) - 1; i >= 0; i-- {
		r := routes[i]
		if vals, match := r.Match(url); match {
			handler, urlStories = r.Handler, vals
			break
		}
	}
	return
}
