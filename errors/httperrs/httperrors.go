package httperrs

import (
	"net/http"

	"github.com/cosiner/gohper/errors"
)

const (
	StatusTooManyRequests = 429
)

type Error interface {
	error
	Code() int
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

func (e HTTPError) Code() int {
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
func NoPayment(err error) Error {
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

// 429
func Frequently(err error) Error {
	return New(err, StatusTooManyRequests)
}

// 500
func Server(err error) Error {
	return New(err, http.StatusInternalServerError)
}

// 501
func NotImplemented(err error) Error {
	return New(err, http.StatusNotImplemented)
}

//503
func Service(err error) Error {
	return New(err, http.StatusServiceUnavailable)
}

// 400
func BadRequestS(err string) Error {
	return NewS(err, http.StatusBadRequest)
}

// 401
func UnAuthS(err string) Error {
	return NewS(err, http.StatusUnauthorized)
}

// 402
func NoPaymentS(err string) Error {
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

// 429
func FrequentlyS(err string) Error {
	return NewS(err, StatusTooManyRequests)
}

// 500
func ServerS(err string) Error {
	return NewS(err, http.StatusInternalServerError)
}

// 501
func NotImplementedS(err string) Error {
	return NewS(err, http.StatusNotImplemented)
}

//503
func ServiceS(err string) Error {
	return NewS(err, http.StatusServiceUnavailable)
}
