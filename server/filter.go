package server

import (
	. "github.com/cosiner/golib/errors"
)

const (
	_BEFORE = iota
	_AFTER
)

type FilterFunc func(*Response, *Request) bool

type Filter interface {
	Init(*Server) error
	Before(*Response, *Request) bool
	After(*Response, *Request) bool
	Destroy()
}

func emptyFilterFunc(_ *Response, _ *Request) bool { return true }

type funcFilter struct {
	before FilterFunc
	after  FilterFunc
}

func (ff *funcFilter) Init(s *Server) error {
	if ff.before == nil {
		ff.before = emptyFilterFunc
	}
	if ff.after == nil {
		ff.after = emptyFilterFunc
	}
	return nil
}

func (ff *funcFilter) Before(resp *Response, req *Request) bool {
	return ff.before(resp, req)
}

func (ff *funcFilter) After(resp *Response, req *Request) bool {
	return ff.after(resp, req)
}

func (ff *funcFilter) Destroy() {}

func (ff *funcFilter) setFilterFunc(when int, filterFunc FilterFunc) error {
	switch when {
	case _BEFORE:
		ff.before = filterFunc
	case _AFTER:
		ff.after = filterFunc
	default:
		return Err("Unsupported filter time")
	}
	return nil
}
