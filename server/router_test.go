package server

import (
	"net/url"
	"testing"

	. "github.com/cosiner/golib/errors"
)

func rout() *router {
	rt := new(router)
	OnErrPanic(rt.AddFuncFilter("/", func(req Request, resp Response, chain FilterChain) {
		chain.Filter(req, resp)
	}))
	OnErrPanic(rt.AddFuncHandler("/user/{id:\\d+}", "GET", func(_ Request, _ Response) {}))
	return rt
}

var rt = rout()

func BenchmarkRouter(b *testing.B) {
	u := &url.URL{Path: "/user/123"}
	for i := 0; i < b.N; i++ {
		_, _, _ = rt.MatchHandler(u)
	}
}
