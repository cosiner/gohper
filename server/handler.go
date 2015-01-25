package server

import (
	"net/http"
	"strings"

	. "github.com/cosiner/golib/errors"
)

//==============================================================================
//                        Request Method
//==============================================================================
const (
	GET            = "GET"
	POST           = "POST"
	DELETE         = "DELETE"
	PUT            = "PUT"
	UNKNOWN_METHOD = "UNKNOWN"
)

func parseRequestMethod(s string) string {
	if s == "" {
		return GET
	}
	return strings.ToUpper(s)
}

type MethodIndicator interface {
	Method(method string) HandlerFunc
}

type methodIndicator struct {
	Handler
}

func (s methodIndicator) Method(method string) (handleFunc HandlerFunc) {
	switch method {
	case GET:
		handleFunc = s.Get
	case POST:
		handleFunc = s.Post
	case DELETE:
		handleFunc = s.Delete
	case PUT:
		handleFunc = s.Put
	}
	return
}

//==============================================================================
//                         Handler
//==============================================================================
type HandlerFunc func(*Response, *Request)

func ErrorHandlerBuilder(httpCode int) HandlerFunc {
	return func(resp *Response, req *Request) {
		resp.WriteHeader(httpCode)
	}
}

var (
	forbiddenHandler        = ErrorHandlerBuilder(http.StatusForbidden)
	notFoundHandler         = ErrorHandlerBuilder(http.StatusNotFound)
	methodNotAllowedHandler = ErrorHandlerBuilder(http.StatusMethodNotAllowed)
)

type Handler interface {
	Init(*Server) error
	Get(*Response, *Request)
	Post(*Response, *Request)
	Delete(*Response, *Request)
	Put(*Response, *Request)
	Destroy()
}

type EmptyHandler int

func (eh EmptyHandler) Init(s *Server) error                { return nil }
func (eh EmptyHandler) Get(resp *Response, req *Request)    {}
func (eh EmptyHandler) Post(resp *Response, req *Request)   {}
func (eh EmptyHandler) Delete(resp *Response, req *Request) {}
func (eh EmptyHandler) Put(resp *Response, req *Request)    {}
func (eh EmptyHandler) Destroy()                            {}

type funcHandler struct {
	EmptyHandler
	get    HandlerFunc
	post   HandlerFunc
	delete HandlerFunc
	put    HandlerFunc
}

func (fh *funcHandler) Method(method string) (handleFunc HandlerFunc) {
	switch method {
	case GET:
		handleFunc = fh.get
	case POST:
		handleFunc = fh.post
	case DELETE:
		handleFunc = fh.delete
	case PUT:
		handleFunc = fh.put
	}
	return
}

func (fh *funcHandler) setMethod(method string, handleFunc HandlerFunc) error {
	switch method {
	case GET:
		fh.get = handleFunc
	case POST:
		fh.post = handleFunc
	case PUT:
		fh.put = handleFunc
	case DELETE:
		fh.delete = handleFunc
	default:
		return Err("Not supported request method")
	}
	return nil
}
