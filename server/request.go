package server

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"

	. "github.com/cosiner/golib/errors"
)

//==============================================================================
//                           Request
//==============================================================================
type Request struct {
	server     *Server
	session    *Session
	response   *Response
	urlStories map[string]string
	params     url.Values
	*http.Request
	header http.Header
}

func newRequest(s *Server, request *http.Request, response *Response) *Request {
	req := &Request{
		server:   s,
		Request:  request,
		response: response,
		header:   request.Header,
	}
	return req
}

func (req *Request) Session() *Session {
	if req.session == nil {
		sess := req.server.newSession()
		req.response.setSessionCookie(sess.sessionId())
	}
	return req.session
}

func (req *Request) setSession(s *Session) {
	req.session = s
}

func (req *Request) Cookie(name string) string {
	if c, err := req.Request.Cookie(name); err == nil {
		return c.Value
	}
	return ""
}

func (req *Request) sessionId() string {
	return req.Cookie(_COOKIE_SESSION)
}

func (req *Request) Header(name string) string {
	return req.header.Get(name)
}

func (req *Request) Param(name string) (value string) {
	params := req.Params(name)
	if len(params) > 0 {
		value = params[0]
	}
	return
}

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

func (req *Request) setUrlStories(urlStories map[string]string) {
	req.urlStories = urlStories
}

func (req *Request) UrlStory(name string) string {
	return req.urlStories[name]
}

func (req *Request) ContentType() string {
	return req.Header(HEADER_CONTENTTYPE)
}

func (req *Request) ResolveJSON() (data map[string]string, err error) {
	return req.unmarshalBody(CONTENTTYPE_JSON, json.Unmarshal)
}

func (req *Request) ResolveXML() (data map[string]string, err error) {
	return req.unmarshalBody(CONTENTTYPE_XML, xml.Unmarshal)
}

type unmarshalFunc func([]byte, interface{}) error

func (req *Request) unmarshalBody(format string, unmarshalFunc unmarshalFunc) (
	data map[string]string, err error) {
	if req.ContentType() != format {
		err = Errorf("Request body is not %s format", format)
	} else {
		var body []byte
		if body, err = ioutil.ReadAll(req.Body); err == nil {
			data = make(map[string]string)
			err = unmarshalFunc(body, data)
		}
	}
	return
}
