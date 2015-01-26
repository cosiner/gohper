package server

import (
	"net/http"

	. "github.com/cosiner/golib/errors"
)

type (
	// HandlerFunc is the common request handler function type
	HandlerFunc func(*Response, *Request)

	// Handler is an common interface of request handler
	// it will be inited on server started, destroyed on server stopped
	// and use it's method to process income request
	// if want to custom the relation of method and handler, just embed
	// EmptyHandler to your Handler, then implement interface MethodIndicator
	Handler interface {
		Init(*Server) error
		Get(*Response, *Request)
		Post(*Response, *Request)
		Delete(*Response, *Request)
		Put(*Response, *Request)
		Destroy()
	}

	// EmptyHandler is an empty handler for user to embed
	EmptyHandler struct{}

	// funcHandler is a handler that use user customed handler function
	// it a funcHandler is defined after a normal handler with same pattern,
	// regardless of only one method is used in funcHandler like Get, it's impossible
	// to use the Post of normal handler, whole normal handler is hidded by this
	// funcHandler, and if other method handler like Post, Put is not set,
	// user access of these method is forbiddened
	funcHandler struct {
		EmptyHandler
		get    HandlerFunc
		post   HandlerFunc
		delete HandlerFunc
		put    HandlerFunc
	}

	// MethodIndicator is an interface for user handler to
	// custom method handle functions
	MethodIndicator interface {
		// Method return an method handle function by method name
		// if nill returned, means access forbidden
		Method(method string) HandlerFunc
	}

	// standardIndicator is the standard method indicator
	standardIndicator struct {
		Handler
	}

	// ErrorHandlers is a collection of http error handler
	ErrorHandlers map[int]HandlerFunc
)

// EmptyHandler methods
func (eh EmptyHandler) Init(s *Server) error                { return nil }
func (eh EmptyHandler) Get(resp *Response, req *Request)    {}
func (eh EmptyHandler) Post(resp *Response, req *Request)   {}
func (eh EmptyHandler) Delete(resp *Response, req *Request) {}
func (eh EmptyHandler) Put(resp *Response, req *Request)    {}
func (eh EmptyHandler) Destroy()                            {}

// funcHandler implements MethodIndicator interface for custom method handler
func (fh *funcHandler) Method(method string) (handlerFunc HandlerFunc) {
	switch method {
	case GET:
		handlerFunc = fh.get
	case POST:
		handlerFunc = fh.post
	case DELETE:
		handlerFunc = fh.delete
	case PUT:
		handlerFunc = fh.put
	}
	return
}

// setMethod setup method handler for funcHandler
func (fh *funcHandler) setMethod(method string, handlerFunc HandlerFunc) error {
	switch method {
	case GET:
		fh.get = handlerFunc
	case POST:
		fh.post = handlerFunc
	case PUT:
		fh.put = handlerFunc
	case DELETE:
		fh.delete = handlerFunc
	default:
		return Err("Not supported request method")
	}
	return nil
}

// ErrorHandlerBuilder build an error handler with given http code
func ErrorHandlerBuilder(httpCode int) HandlerFunc {
	return func(resp *Response, req *Request) {
		resp.WriteHeader(httpCode)
	}
}

// indicateMethod indicate handler function from a handler and method
func indicateMethod(handler Handler, method string) HandlerFunc {
	var indicator MethodIndicator
	switch handler := handler.(type) {
	case MethodIndicator:
		indicator = handler
	default:
		indicator = standardIndicator{handler}
	}
	return indicator.Method(method)
}

// Method indicate method handle function, for standardIndicator,
// each method indicate the function with same name, such as GET->Get...
func (s standardIndicator) Method(method string) (handlerFunc HandlerFunc) {
	switch method {
	case GET:
		handlerFunc = s.Get
	case POST:
		handlerFunc = s.Post
	case DELETE:
		handlerFunc = s.Delete
	case PUT:
		handlerFunc = s.Put
	}
	return
}

// NewErrorHandlers create new ErrorHandlers
func NewErrorHandlers() ErrorHandlers {
	return ErrorHandlers{
		http.StatusForbidden:        ErrorHandlerBuilder(http.StatusForbidden),
		http.StatusNotFound:         ErrorHandlerBuilder(http.StatusNotFound),
		http.StatusMethodNotAllowed: ErrorHandlerBuilder(http.StatusMethodNotAllowed),
	}
}

// ForbiddenHandler return forbidden error handler
func (eh ErrorHandlers) ForbiddenHandler() HandlerFunc {
	return eh[http.StatusForbidden]
}

// SetForbiddenHandler set forbidden error handler
func (eh ErrorHandlers) SetForbiddenhandler(handlerFunc HandlerFunc) {
	eh[http.StatusForbidden] = handlerFunc
}

// NotFoundHandler return notfound error handler
func (eh ErrorHandlers) NotFoundHandler() HandlerFunc {
	return eh[http.StatusNotFound]
}

// SetNotFoundHandler set notfound error handler
func (eh ErrorHandlers) SetNotFoundHandler(handlerFunc HandlerFunc) {
	eh[http.StatusNotFound] = handlerFunc
}

// MethodNotAllowedHandler return methodnotallowed error handler
func (eh ErrorHandlers) MethodNotAllowedHandler() HandlerFunc {
	return eh[http.StatusMethodNotAllowed]
}

// SetMethodNotAllowedHandler set methodnotallowed error handler
func (eh ErrorHandlers) SetMethodNotAllowedHandler(handlerFunc HandlerFunc) {
	eh[http.StatusMethodNotAllowed] = handlerFunc
}
