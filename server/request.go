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
		server   *Server
		session  *Session
		response *Response
		// urlVars is the values extract from url path for REST url
		urlVars map[string]string
		// params represent user request's parameters,
		// for GET it's exist in url, for other method, parse from form
		params url.Values
		*http.Request
		header http.Header
	}

	// unmarshalFunc is the type of unmarshal function
	unmarshalFunc func([]byte, interface{}) error
)

// newRequest create a new request
// the parameter response is only used for request to setup session cookie
// when user call request.Session()
func newRequest(s *Server, request *http.Request, response *Response) *Request {
	return &Request{
		server:   s,
		Request:  request,
		response: response,
		header:   request.Header,
	}
}

// Server return the running server
func (req *Request) Server() *Server {
	return req.server
}

// Session return the session that request blong to
func (req *Request) Session() (sess *Session) {
	if sess = req.session; sess == nil { // no session
		if id := req.cookieSessionId(); id != "" {
			sess = req.server.session(id) // get session from server store
		}
		if sess == nil { // server stored session has been expired, create new session
			sess := req.server.newSession()
			req.response.setSessionCookie(sess.sessionId()) // write session cookie to response
		}
		req.session = sess
	}
	return
}

// Cookie return cookie value with given name
func (req *Request) Cookie(name string) string {
	if c, err := req.Request.Cookie(name); err == nil {
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
	return parseContentType(req.header.Get(name))
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
	params := req.params
	if params == nil {
		switch req.Method {
		case GET:
			params = req.URL.Query()
		default:
			req.ParseForm()
			params = req.PostForm
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
	return req.Header(HEADER_CONTENTTYPE)
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
		if body, err = ioutil.ReadAll(req.Body); err == nil {
			err = unmarshalFunc(body, v)
		}
	}
	return
}
