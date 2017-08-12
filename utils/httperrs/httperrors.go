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
