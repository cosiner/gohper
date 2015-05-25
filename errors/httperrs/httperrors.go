package httperrs

import "github.com/cosiner/gohper/errors"

type Code int

func (c Code) New(err error) Error {
	return New(err, int(c))
}

func (c Code) NewS(err string) Error {
	return NewS(err, int(c))
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

type Error interface {
	error
	Code() int
}

type HTTPError struct {
	error
	code int
}

func Must(err error) Error {
	if err == nil {
		return nil
	}

	return err.(Error)
}

func New(err error, code int) Error {
	if err == nil {
		return nil
	}

	return HTTPError{
		error: err,
		code:  code,
	}
}

func NewS(err string, code int) Error {
	if err == "" {
		return nil
	}

	return HTTPError{
		error: errors.Err(err),
		code:  code,
	}
}

func (e HTTPError) Code() int {
	return e.code
}
