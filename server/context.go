package server

import (
	"net/http"
)

// context is the request specified enviroment of request and response
type context struct {
	srv     *Server
	sess    *Session
	w       http.ResponseWriter
	request *http.Request
	resp    *Response
	req     *Request
	AttrContainer
}

// newContext create a new context
func newContext(s *Server, w http.ResponseWriter, request *http.Request) *context {
	return &context{
		srv:           s,
		w:             w,
		request:       request,
		AttrContainer: make(Values),
	}
}

// setup set up response and request
func (ctx *context) setup(resp *Response, req *Request) {
	ctx.resp = resp
	ctx.req = req
}

// destroy destroy all reference the context keep
func (ctx *context) destroy() {
	ctx.srv = nil
	ctx.sess = nil
	ctx.w = nil
	ctx.request = nil
	ctx.resp = nil
	ctx.req = nil
	ctx.AttrContainer = nil
}

// Server return the only running server
func (ctx *context) Server() *Server {
	return ctx.srv
}

// Session return session that request belong to
// if it's not exist in server's session manager and session store or
// has been expired, then create a new session and set up session cookie
// to client
func (ctx *context) Session() (sess *Session) {
	if sess = ctx.sess; sess == nil { // no session
		if id := ctx.req.cookieSessionId(); id != "" {
			sess = ctx.srv.session(id) // get session from server store
		}
		if sess == nil { // server stored session has been expired, create new session
			sess := ctx.srv.newSession()
			ctx.resp.setSessionCookie(sess.sessionId()) // write session cookie to response
		}
		ctx.sess = sess
	}
	return sess
}
