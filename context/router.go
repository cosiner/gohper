package context

import (
	"github.com/cosiner/golib/regexp/urlmatcher"
)

type Route struct {
	*urlmatcher.Matcher
	Handler
}

type Router struct {
	literalRoutes []*Route
	regexpRoutes  []*Route
}

func (rt *Router) Handle(url string) error {
	m := urlmatcher.LiteralMatch(rt.literalRoutes, url)
	if m != nil {
		matcher := m.(*Route)
	}
}
