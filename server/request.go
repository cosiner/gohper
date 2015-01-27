package server

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"

	. "github.com/cosiner/golib/errors"
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

// ResolveJSON resolve json data, request's content type MUST BE JSON
func (req *Request) ResolveJSON() (data map[string]string, err error) {
	return req.unmarshalBody(CONTENTTYPE_JSON, json.Unmarshal)
}

// ResolveJSONInto resolve json data into given parameter,
// parameter MUST BE POINTER, request's content type MUST BE JSON
func (req *Request) ResolveJSONInto(v interface{}) error {
	return req.unmarshalBodyInto(CONTENTTYPE_JSON, json.Unmarshal, v)
}

// ResolveXML resolve xml data, request's content type MUST BE XML
func (req *Request) ResolveXML() (data map[string]string, err error) {
	return req.unmarshalBody(CONTENTTYPE_XML, xml.Unmarshal)
}

// ResolveXMLInto resolve xml data into given parameter
// parameter MUST BE POINTER, request's content type MUST BE XML
func (req *Request) ResolveXMLInto(v interface{}) (err error) {
	return req.unmarshalBodyInto(CONTENTTYPE_XML, xml.Unmarshal, v)
}

// unmarshalBody unmarshal request body with given unmarshal function
func (req *Request) unmarshalBody(format string, unmarshalFunc unmarshalFunc) (
	map[string]string, error) {
	data := make(map[string]string)
	err := req.unmarshalBodyInto(format, unmarshalFunc, &data)
	if err != nil {
		data = nil
	}
	return data, err
}

// unmarshanBodyInto unmarshan request body into given parameter's memory space
// parameter MUST BE POINTER, and request's content must equal to given format
func (req *Request) unmarshalBodyInto(format string, unmarshalFunc unmarshalFunc,
	v interface{}) (err error) {
	if req.ContentType() != format {
		err = Errorf("Request body is not %s format", format)
	} else {
		var body []byte
		if body, err = ioutil.ReadAll(req.request.Body); err == nil {
			err = unmarshalFunc(body, v)
		}
	}
	return
}

// Forward forward to given address use exist request and response
func (req *Request) Forward(addr string) error {
	u, err := url.Parse(addr)
	if err == nil {
		req.forwarded = true
		req.Server().serve(u, req.resp, req, true)
	}
	return err
}
