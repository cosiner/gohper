package server

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/cosiner/golib/encoding"
)

type (
	// Request represent an income request
	Request struct {
		*context
		request *http.Request
		method  string
		// urlVars is the values extract from url path for REST url
		urlVars map[string]string
		// params represent user request's parameters,
		// for GET it's exist in url, for other method, parse from form
		params url.Values
		header http.Header
	}

	// unmarshalFunc is the type of unmarshal function
	unmarshalFunc func([]byte, interface{}) error
)

// newRequest create a new request
func newRequest(ctx *context, request *http.Request) *Request {
	return &Request{
		context: ctx,
		request: request,
		method:  parseRequestMethod(request.Method),
		header:  request.Header,
	}
}

// destroy destroy all reference that request keep
func (req *Request) destroy() {
	req.context.destroy()
	req.request = nil
	req.urlVars = nil
	req.params = nil
	req.header = nil
}

// setMethod set up request's method
func (req *Request) setMethod(method string) {
	req.method = method
}

// Method return method of request
func (req *Request) Method() string {
	return req.method
}

// Cookie return cookie value with given name
func (req *Request) Cookie(name string) string {
	if c, err := req.request.Cookie(name); err == nil {
		return c.Value
	}
	return ""
}

// RemoteAddr return remote address
func (req *Request) RemoteAddr() string {
	return req.request.RemoteAddr
}

// URL return request url
func (req *Request) URL() *url.URL {
	return req.request.URL
}

// cookieSessionId extract session id from cookie
func (req *Request) cookieSessionId() string {
	return req.Cookie(_COOKIE_SESSION)
}

// Header return header value with name
func (req *Request) Header(name string) string {
	return req.header.Get(name)
}

// Param return request parameter with name
func (req *Request) Param(name string) (value string) {
	params := req.Params(name)
	if len(params) > 0 {
		value = params[0]
	}
	return
}

// Params return request parameters with name
func (req *Request) Params(name string) []string {
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
func (req *Request) setUrlVars(urlVars map[string]string) {
	req.urlVars = urlVars
}

// UrlVar return url variable value with name
func (req *Request) UrlVar(name string) (value string) {
	if req.urlVars != nil {
		value = req.urlVars[name]
	}
	return
}

// ContentType extract content type form request header
func (req *Request) ContentType() string {
	return parseContentType(req.Header(HEADER_CONTENTTYPE))
}

// Forward forward to given address use exist request and response
func (req *Request) Forward(addr string) error {
	u, err := url.Parse(addr)
	if err == nil {
		req.Server().processHttpRequest(u, req, req.resp, true)
	}
	return err
}

// ReadString read request body as string
func (req *Request) ReadString() (s string) {
	s, _ = encoding.ReadString(req.request.Body)
	return
}

func (req *Request) Read(data []byte) (int, error) {
	return req.request.Body.Read(data)
}

// ReadAll read request body as bytes
func (req *Request) ReadAll() (bs []byte) {
	bs, _ = ioutil.ReadAll(req.request.Body)
	return
}

// ReadJSON read json data into given parameter,
// parameter MUST BE POINTER, request's content type MUST BE JSON
func (req *Request) ReadJSON(v interface{}) error {
	return encoding.ReadJSON(req, v)
}

// ReadXML read xml data into given parameter
// parameter MUST BE POINTER, request's content type MUST BE XML
func (req *Request) ReadXML(v interface{}) (err error) {
	return encoding.ReadXML(req, v)
}
