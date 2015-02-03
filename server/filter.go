package server

type (
	// FilterFunc represent common filter function type,
	FilterFunc func(Request, Response, FilterChain)

	// Filter is an filter that run before or after handler,
	// to modify or check request and response
	// it will be inited on server started, destroyed on server stopped
	Filter interface {
		Init(*Server) error
		Destroy()
		Filter(Request, Response, FilterChain)
	}

	// FilterChain represent a chain of filter, the last is final handler
	// to continue the chain, must call chain.Filter
	FilterChain interface {
		Filter(Request, Response)
	}

	filterChain struct {
		index   int
		filters []Filter
		handler HandlerFunc
	}
)

// FilterFunc is a function Filter
func (FilterFunc) Init(*Server) error { return nil }
func (FilterFunc) Destroy()           {}
func (fn FilterFunc) Filter(req Request, resp Response, chain FilterChain) {
	fn(req, resp, chain)
}

// NewFilterChain create a chain of filter
// this method is setup to public for which condition there only one route
// need filter, if add to global router, it will make router match slower
// this method can help for these condition
func NewFilterChain(filters []Filter, handler HandlerFunc) FilterChain {
	return &filterChain{
		index:   -1,
		filters: filters,
		handler: handler,
	}
}

// Filter call next filter, if there is no next filter,then call final handler
func (chain *filterChain) Filter(req Request, resp Response) {
	chain.index++
	index, filters := chain.index, chain.filters
	if index == len(filters) {
		if handler := chain.handler; handler != nil {
			handler(req, resp)
		}
	} else {
		filters[index].Filter(req, resp, chain)
	}
}
