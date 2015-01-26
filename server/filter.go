package server

import (
	. "github.com/cosiner/golib/errors"
)

type (
	// FilterFunc represent common filter function type,
	// returned value means whether execute remains filters
	FilterFunc func(*Response, *Request) bool

	// Filter is an filter that run before or after handler,
	// to modify or check request and response
	// it will be inited on server started, destroyed on server stopped
	Filter interface {
		Init(*Server) error
		Before(*Response, *Request) bool
		After(*Response, *Request) bool
		Destroy()
	}

	// funcFilter is a filter that use user customed filter function
	funcFilter struct {
		before FilterFunc
		after  FilterFunc
	}
)

// emptyFilterFunc is just an empty filter like it's name
func emptyFilterFunc(_ *Response, _ *Request) bool { return true }

// Init init funcFiliter, if one of it's filter function is nil,
// then it will be inited as emptyFilterFunc defined above
func (ff *funcFilter) Init(s *Server) error {
	if ff.before == nil {
		ff.before = emptyFilterFunc
	}
	if ff.after == nil {
		ff.after = emptyFilterFunc
	}
	return nil
}

// Before run before handler
func (ff *funcFilter) Before(resp *Response, req *Request) bool {
	return ff.before(resp, req)
}

// After run after handler
func (ff *funcFilter) After(resp *Response, req *Request) bool {
	return ff.after(resp, req)
}

// Destroy run when destroy filter
func (ff *funcFilter) Destroy() {}

// setFilterFunc setup filter function for funcFilter
func (ff *funcFilter) setFilterFunc(when int, filterFunc FilterFunc) error {
	switch when {
	case _FILTER_BEFORE:
		ff.before = filterFunc
	case _FILTER_AFTER:
		ff.after = filterFunc
	default:
		return Err("Unsupported filter time")
	}
	return nil
}
