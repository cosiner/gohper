package server

import (
	"strings"

	"github.com/cosiner/golib/types"
)

const (
	// Filter time
	_FILTER_BEFORE = iota // _BEFORE means execute filter before handler
	_FILTER_AFTER         // _AFTER means execute filter after handler

	// Http Header
	HEADER_CONTENTTYPE     = "Content-Type"
	HEADER_CONTENTLENGTH   = "Content-Length"
	HEADER_SETCOOKIE       = "Set-Cookie"
	HEADER_REFER           = "Referer"
	HEADER_CONTENTENCODING = "Content-Encoding"
	HEADER_USERAGENT       = "User-Agent"

	// Xsrf
	HEADER_XSRFTOKEN  = "X-XSRFToken"
	HEADER_CSRFTOKEN  = "X-CSRFToken"
	XSRF_NAME         = "_xsrf"
	XSRF_ONERRORTOKEN = "dsajhdoqwARUH20174P2UAsdJDASDKJ"
	XSRF_DEFFLUSH     = 60
	XSRF_DEFLIFETIME  = 300
	XSRF_ERRORCODE    = 101
	XSRF_FORMHEAD     = `<input type="hidden" name="` + XSRF_NAME + `" value="`
	XSRF_FORMEND      = `"/>`

	// ContentEncoding
	ENCODING_GZIP = "gzip"

	// Request Method
	GET            = "GET"
	POST           = "POST"
	DELETE         = "DELETE"
	PUT            = "PUT"
	PATCH          = "PATCH"
	UNKNOWN_METHOD = "UNKNOWN"

	// Content Type
	CONTNTTYPE_PLAIN = "text/plain"
	CONTENTTYPE_HTML = "text/html"
	CONTENTTYPE_XML  = "application/xml"
	CONTENTTYPE_JSON = "application/json"

	// Session
	_COOKIE_SESSION            = "session"
	SESSION_MEM_GCINTERVAL     = "gcinterval"
	SESSION_MEM_RMBACKLOG      = "rmbacklog"
	DEF_SESSION_MEMSTORE_CONF  = "gcinterval=600&rmbacklog=10"
	DEF_SESSION_MEM_GCINTERVAL = 600
	DEF_SESSION_MEM_RMBACKLOG  = 10
	DEF_SESSION_LIFETIME       = 600
)

// parseRequestMethod convert a string to request method, default use GET
// if string is empty
func parseRequestMethod(s string) string {
	if s == "" {
		return GET
	}
	return strings.ToUpper(s)
}

// parseContentType parse content type
func parseContentType(str string) string {
	if str == "" {
		return CONTENTTYPE_HTML
	}
	return types.TrimLower(str)
}

// xsrfFormHTML return xsrf form html string with given token
func xsrfFormHTML(tok string) string {
	return XSRF_FORMHEAD + tok + XSRF_FORMEND
}
