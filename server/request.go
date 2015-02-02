package server

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cosiner/golib/encoding"
)

type (
	Request interface {
		RemoteAddr() string
		URL() *url.URL
		Method() string
		ContentType() string
		Header(name string) string
		Cookie(name string) string
		Session() *Session
		Server() *Server
		Param(name string) (value string)
		Params(name string) []string
		UrlVar(name string) (value string)
		ScanUrlVars(vars ...*string)
		Forward(addr string) error
		ReadString() (s string)
		Read(data []byte) (int, error)
		ReadAll() (bs []byte)
		ReadJSON(v interface{}) error
		ReadXML(v interface{}) (err error)
		AttrContainer
	}

	// request represent an income request
	request struct {
		*context
		request *http.Request
		method  string
		indexer VarIndexer
		// urlVars is the values extract from url path for REST url
		urlVars []string
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
		context: ctx,
		request: requ,
		header:  requ.Header,
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
	req.urlVars = nil
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

// RemoteAddr return remote address
func (req *request) RemoteAddr() string {
	return req.request.RemoteAddr
}

// URL return request url
func (req *request) URL() *url.URL {
	return req.request.URL
}

// cookieSessionId extract session id from cookie
func (req *request) cookieSessionId() string {
	return req.Cookie(_COOKIE_SESSION)
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

// setUrlVars setup request url variables
func (req *request) setUrlVars(indexer VarIndexer, urlVars []string) *request {
	req.indexer = indexer
	req.urlVars = urlVars
	return req
}

// UrlVar return url variable value with name
func (req *request) UrlVar(name string) (value string) {
	if values := req.urlVars; values != nil {
		value = req.indexer.ValueOf(values, name)
	}
	return
}

// ScanUrlVars scan url variable values into given address
func (req *request) ScanUrlVars(vars ...*string) {
	if values := req.urlVars; values != nil {
		req.indexer.ScanInto(values, vars...)
	}
}

// ContentType extract content type form request header
func (req *request) ContentType() string {
	return parseContentType(req.Header(HEADER_CONTENTTYPE))
}

// Forward forward to given address use exist request and response
func (req *request) Forward(addr string) error {
	u, err := url.Parse(addr)
	if err == nil {
		req.Server().processHttpRequest(u, req, req.resp, true)
	}
	return err
}

// ReadString read request body as string
func (req *request) ReadString() (s string) {
	s, _ = encoding.ReadString(req.request.Body)
	return
}

func (req *request) Read(data []byte) (int, error) {
	return req.request.Body.Read(data)
}

// ReadAll read request body as bytes
func (req *request) ReadAll() (bs []byte) {
	bs, _ = ioutil.ReadAll(req.request.Body)
	return
}

// ReadJSON read json data into given parameter,
// parameter MUST BE POINTER, request's content type MUST BE JSON
func (req *request) ReadJSON(v interface{}) error {
	return encoding.ReadJSON(req, v)
}

// ReadXML read xml data into given parameter
// parameter MUST BE POINTER, request's content type MUST BE XML
func (req *request) ReadXML(v interface{}) (err error) {
	return encoding.ReadXML(req, v)
}

// hasSession return whether a request has own session
func (req *request) hasSession() bool {
	return req.sess != nil
}
