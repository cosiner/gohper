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
	HEADER_CONTENTTYPE   = "Content-Type"
	HEADER_CONTENTLENGTH = "Content-Length"
	HEADER_SETCOOKIE     = "Set-Cookie"

	// Request Method
	GET            = "GET"
	POST           = "POST"
	DELETE         = "DELETE"
	PUT            = "PUT"
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
