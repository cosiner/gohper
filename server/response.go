package server

import (
	"net/http"

	"github.com/cosiner/golib/encoding"

	"github.com/cosiner/golib/types"
)

type (
	// Response represent a response of request to user
	Response struct {
		*context
		w http.ResponseWriter
		*types.WriterChain
		header http.Header
	}

	// marshalFunc is the marshal function type
	marshalFunc func(interface{}) ([]byte, error)
)

// newResponse create a new Response, and set default content type to HTML
func newResponse(ctx *context, w http.ResponseWriter) *Response {
	resp := &Response{
		context:     ctx,
		w:           w,
		WriterChain: types.NewWriterChain(w),
		header:      w.Header(),
	}
	resp.SetContentType(CONTENTTYPE_HTML)
	return resp
}

// destroy destroy all reference that response keep
func (resp *Response) destroy() {
	resp.context.destroy()
	resp.w = nil
	resp.WriterChain = nil
	resp.header = nil
}

// SetHeader setup response header
func (resp *Response) SetHeader(name, value string) {
	resp.header.Set(name, value)
}

// AddHeader add a value to response header
func (resp *Response) AddHeader(name, value string) {
	resp.header.Add(name, value)
}

// SetContentType set content type of response
func (resp *Response) SetContentType(typ string) {
	resp.SetHeader(HEADER_CONTENTTYPE, typ)
}

// contentType return current content type of response
func (resp *Response) contentType() string {
	return resp.header.Get(HEADER_CONTENTTYPE)
}

// newCookie create a new Cookie and return it's displayed string
// parameter lifetime is time by second
func (*Response) newCookie(name, value string, lifetime int) string {
	return (&http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: lifetime,
	}).String()
}

// SetCookie setup response cookie, default age is default browser opened time
func (resp *Response) SetCookie(name, value string) {
	resp.SetCookieWithExpire(name, value, 0)
}

// SetCookieWithExpire setup response cookie with lifetime
func (resp *Response) SetCookieWithExpire(name, value string, lifetime int) {
	resp.SetHeader(HEADER_SETCOOKIE, resp.newCookie(name, value, lifetime))
}

// DeleteClientCookie delete user briwser's cookie by name
func (resp *Response) DeleteClientCookie(name string) {
	resp.SetCookieWithExpire(name, "", -1)
}

// setSessionCookie setup session cookie
func (resp *Response) setSessionCookie(id string) {
	resp.SetCookie(_COOKIE_SESSION, id)
}

// Redirect redirect to new url
func (resp *Response) Redirect(url string) {
	http.Redirect(resp.w, resp.request, url, http.StatusTemporaryRedirect)
}

// PermanentRedirect permanently redirect current request url to new url
func (resp *Response) PermanentRedirect(url string) {
	http.Redirect(resp.w, resp.request, url, http.StatusMovedPermanently)
}

// Report Error report an http error with given status code
func (resp *Response) ReportError(statusCode int) {
	resp.w.WriteHeader(statusCode)
}

// Render render template with context
func (resp *Response) Render(tmpl string) error {
	return resp.Server().renderTemplate(resp, tmpl, resp.context)
}

// BaseWriter return base response writer
func (resp *Response) BaseWriter() http.ResponseWriter {
	return resp.w
}

// WriteString write sting to client
func (resp *Response) WriteString(data string) (int, error) {
	return encoding.WriteString(resp, data)
}

// WriteJSON write json data to client, and setup content type to json
func (resp *Response) WriteJSON(val interface{}) error {
	return encoding.WriteJSON(resp, val)
}

// WriteXML write xml data to client, and setup content type to xml
func (resp *Response) WriteXML(val interface{}) error {
	return encoding.WriteXML(resp, val)
}

// Flush flush response's output
func (resp *Response) Flush() {
	if flusher, is := resp.WriterChain.Writer.(http.Flusher); is {
		flusher.Flush()
	}
}
