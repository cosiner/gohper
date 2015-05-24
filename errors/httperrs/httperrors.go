package httperrs

import (
	"net/http"

	"github.com/cosiner/gohper/errors"
)

type Error interface {
	error
	HTTPCode() int
}

type HTTPError struct {
	error
	code int
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

func (e HTTPError) HTTPCode() int {
	return e.code
}

// 400
func BadRequest(err error) Error {
	return New(err, http.StatusBadRequest)
}

// 401
func UnAuth(err error) Error {
	return New(err, http.StatusUnauthorized)
}

// 402
func PaymentRequired(err error) Error {
	return New(err, http.StatusPaymentRequired)
}

// 403
func Forbidden(err error) Error {
	return New(err, http.StatusForbidden)
}

// 404
func NotFound(err error) Error {
	return New(err, http.StatusNotFound)
}

// 405
func NotAllowedMethod(err error) Error {
	return New(err, http.StatusMethodNotAllowed)
}

// 406
func NotAcceptable(err error) Error {
	return New(err, http.StatusNotAcceptable)
}

// 407
func ProxyUnAuth(err error) Error {
	return New(err, http.StatusProxyAuthRequired)
}

// 408
func Timeout(err error) Error {
	return New(err, http.StatusRequestTimeout)
}

// 409
func Conflict(err error) Error {
	return New(err, http.StatusConflict)
}

// 410
func Gone(err error) Error {
	return New(err, http.StatusGone)
}

// 411
func NoLength(err error) Error {
	return New(err, http.StatusLengthRequired)
}

// 400
func BadRequestS(err error) Error {
	return New(err, http.StatusBadRequest)
}

// 401
func UnAuthS(err string) Error {
	return NewS(err, http.StatusUnauthorized)
}

// 402
func PaymentRequiredS(err string) Error {
	return NewS(err, http.StatusPaymentRequired)
}

// 403
func ForbiddenS(err string) Error {
	return NewS(err, http.StatusForbidden)
}

// 404
func NotFoundS(err string) Error {
	return NewS(err, http.StatusNotFound)
}

// 405
func NotAllowedMethodS(err string) Error {
	return NewS(err, http.StatusMethodNotAllowed)
}

// 406
func NotAcceptableS(err string) Error {
	return NewS(err, http.StatusNotAcceptable)
}

// 407
func ProxyUnAuthS(err string) Error {
	return NewS(err, http.StatusProxyAuthRequired)
}

// 408
func TimeoutS(err string) Error {
	return NewS(err, http.StatusRequestTimeout)
}

// 409
func ConflictS(err string) Error {
	return NewS(err, http.StatusConflict)
}

// 410
func GoneS(err string) Error {
	return NewS(err, http.StatusGone)
}

// 411
func NoLengthS(err string) Error {
	return NewS(err, http.StatusLengthRequired)
}
