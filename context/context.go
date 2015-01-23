package context

import (
	"net/http"

	"github.com/cosiner/gomodule/memcache"
	"github.com/cosiner/gomodule/rediscache"
)

type Server struct {
	context      memcache.MemCache
	sessionStore *rediscache.RedisCache
}

type Handler interface {
}

func (s Server) ServeHttp(w http.ResponseWriter, req *http.Request) {
	http.Response
}

func (s *Server) Start(addr string) error {
	http.ListenAndServe(addr, s)
}

type Session struct {
	id      uint64
	server  *Server
	context *rediscache.RedisCache
}

type Request struct {
	session *Session
	context memcache.MemCache
	*http.Request
}

func (req *Request) Session() *Session {
	return req.session
}

func (req *Request) setSession(s *Session) {
	req.session = s
}

func (req *Request) Cookie(name string) string {
}

type Response struct {
	wr http.ResponseWriter
}
