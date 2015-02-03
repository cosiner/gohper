package filters

import (
	"github.com/cosiner/golib/crypto"
	"github.com/cosiner/golib/types"
	. "github.com/cosiner/gomodule/server"
)

type (
	CookieRequestWrapper struct {
		Request
	}

	CookieResponseWrapper struct {
		Response
	}
)

// SecretKey is the key of signature, to use this filter, you must setup it
var SecretKey []byte

func (reqw CookieRequestWrapper) SecureCookie(name string) string {
	userCookie := reqw.Request.SecureCookie(name)
	if userCookie != "" {
		userCookie = types.UnsafeString(
			crypto.VerifySecret(SecretKey,
				types.UnsafeBytes(userCookie)))
	}
	return userCookie
}

func (resw CookieResponseWrapper) SetSecureCookieWithExpire(name string, value string, lifetime int) {
	encVal := crypto.SignSecret(SecretKey, types.UnsafeBytes(value))
	resw.Response.SetSecureCookieWithExpire(name, types.UnsafeString(encVal), lifetime)
}

func SecureCookieFilter(req Request, resp Response, chain FilterChain) {
	chain.Filter(CookieRequestWrapper{req},
		CookieResponseWrapper{resp})
}
