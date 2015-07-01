// Package httperrs provide error types help interactive over http
package httperrs

import "github.com/cosiner/gohper/errors"

type Error interface {
	error
	Code() int
}

type HTTPError struct {
	error
	code int
}

func (e HTTPError) Code() int {
	return e.code
}

type Code int

func (c Code) New(err error) Error {
	if err == nil {
		return nil
	}

	return HTTPError{
		error: err,
		code:  int(c),
	}
}

func (c Code) NewS(err string) Error {
	if err == "" {
		return nil
	}

	return HTTPError{
		error: errors.Err(err),
		code:  int(c),
	}
}

const (
	BadRequest  Code = 400
	UnAuth      Code = 401
	NoPayment   Code = 402
	Forbidden   Code = 403
	NotFound    Code = 404
	NotAllowed  Code = 405
	NotAccept   Code = 406
	ProxyUnAuth Code = 407
	Timeout     Code = 408
	Conflict    Code = 409
	Gone        Code = 410
	NoLength    Code = 411
	Frequently  Code = 429

	Server         Code = 500
	NotImplemented Code = 501
	Service        Code = 503
)

func Must(err error) Error {
	if err == nil {
		return nil
	}

	return err.(Error)
}

func New(err error, code int) Error {
	return Code(code).New(err)
}

func NewS(err string, code int) Error {
	return Code(code).NewS(err)
}
