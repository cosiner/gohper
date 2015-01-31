package server

type (
	// FilterFunc represent common filter function type,
	FilterFunc func(*Request, *Response, *FilterChain)

	// Filter is an filter that run before or after handler,
	// to modify or check request and response
	// it will be inited on server started, destroyed on server stopped
	Filter interface {
		Init(*Server) error
		Destroy()
		Filter(*Request, *Response, *FilterChain)
	}

	// FilterChain represent a chain of filter, the last is final handler
	FilterChain struct {
		index   int
		filters []Filter
		handler HandlerFunc
	}
)

// FilterFunc is a function Filter
func (FilterFunc) Init(*Server) error { return nil }
func (FilterFunc) Destroy()           {}
func (fn FilterFunc) Filter(req *Request, resp *Response, chain *FilterChain) {
	fn(req, resp, chain)
}

// newFilterChain create a chain of filter
func newFilterChain(filters []Filter, handler HandlerFunc) *FilterChain {
	return &FilterChain{
		index:   0,
		filters: filters,
		handler: handler,
	}
}

// Filter call next filter, if there is no next filter,then call final handler
func (chain *FilterChain) Filter(req *Request, resp *Response) {
	index, filters := chain.index, chain.filters
	chain.index++
	if index == len(filters) {
		chain.handler(req, resp)
	} else {
		filters[index].Filter(req, resp, chain)
	}
}
