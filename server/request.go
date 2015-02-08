package server

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/cosiner/golib/encoding"
)

type (
	Request interface {
		RemoteAddr() string
		Refer() string
		UserAgent() string
		URL() *url.URL
		Method() string
		ContentType() string
		ContentEncoding() string
		Header(name string) string
		Cookie(name string) string
		SecureCookie(name string) string
		Session() *Session
		Server() *Server
		Param(name string) (value string)
		Params(name string) []string
		Forward(addr string) error
		encoding.PowerReader
		UrlVarIndexer
		AttrContainer
	}

	// request represent an income request
	request struct {
		*context
		UrlVarIndexer
		encoding.PowerReader
		request *http.Request
		method  string
		// params represent user request's parameters,
		// for GET it's exist in url, for other method, parse from form
		params url.Values
		header http.Header
	}

	// unmarshalFunc is the type of unmarshal function
	unmarshalFunc func([]byte, interface{}) error
)

// newRequest create a new request
func newRequest(ctx *context, requ *http.Request) *request {
	req := &request{
		context:     ctx,
		PowerReader: encoding.NewPowerReader(requ.Body),
		request:     requ,
		header:      requ.Header,
	}
	method := requ.Method
	if m := requ.Header.Get("X-HTTP-Method-Override"); method == POST && m != "" {
		method = m
	}
	req.method = strings.ToUpper(method)
	return req
}

// destroy destroy all reference that request keep
func (req *request) destroy() {
	req.context.destroy()
	req.request = nil
	req.params = nil
	req.header = nil
}

// Method return method of request
func (req *request) Method() string {
	return req.method
}

// Cookie return cookie value with given name
func (req *request) Cookie(name string) string {
	if c, err := req.request.Cookie(name); err == nil {
		return c.Value
	}
	return ""
}

// SecureCookie return secure cookie, currently it's just call Cookie without
// 'Secure', if need this feture, just put an filter before handler
// and override this method
func (req *request) SecureCookie(name string) string {
	return req.Cookie(name)
}

// RemoteAddr return remote address
func (req *request) RemoteAddr() string {
	return req.request.RemoteAddr
}

// Refer return where user from
func (req *request) Refer() string {
	return req.Header(HEADER_REFER)
}

// UserAgent return user's agent identify
func (req *request) UserAgent() string {
	return req.Header(HEADER_USERAGENT)
}

// ContentType extract content type form request header
func (req *request) ContentType() string {
	return req.Header(HEADER_CONTENTTYPE)
}

func (req *request) ContentEncoding() string {
	return req.Header(HEADER_CONTENTENCODING)
}

// URL return request url
func (req *request) URL() *url.URL {
	return req.request.URL
}

// cookieSessionId extract session id from cookie, if enable secure cookie, it will
// use it
func (req *request) cookieSessionId() string {
	return req.SecureCookie(_COOKIE_SESSION)
}

// Header return header value with name
func (req *request) Header(name string) string {
	return req.header.Get(name)
}

// Param return request parameter with name
func (req *request) Param(name string) (value string) {
	params := req.Params(name)
	if len(params) > 0 {
		value = params[0]
	}
	return
}

// Params return request parameters with name
func (req *request) Params(name string) []string {
	params, request := req.params, req.request
	if params == nil {
		switch req.method {
		case GET:
			params = request.URL.Query()
		default:
			request.ParseForm()
			params = request.PostForm
		}
		req.params = params
	}
	return params[name]
}

// setVarIndexer setup request url variables
func (req *request) setVarIndexer(indexer UrlVarIndexer) *request {
	req.UrlVarIndexer = indexer
	return req
}

// Forward forward to given address use exist request and response
func (req *request) Forward(addr string) error {
	u, err := url.Parse(addr)
	if err == nil {
		req.Server().processHttpRequest(u, req, req.resp, true)
	}
	return err
}

// hasSession return whether a request has own session
func (req *request) hasSession() bool {
	return req.sess != nil
}
