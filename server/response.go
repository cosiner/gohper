package server

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

type (
	// Response represent a response of request to user
	Response struct {
		server  *Server
		request *http.Request
		http.ResponseWriter
		header http.Header
	}
	// marshalFunc is the marshal function type
	marshalFunc func(interface{}) ([]byte, error)
)

// newResponse create a new Response, and set default content type to HTML
func newResponse(s *Server, request *http.Request, w http.ResponseWriter) *Response {
	resp := &Response{
		server:         s,
		request:        request,
		ResponseWriter: w,
		header:         w.Header(),
	}
	resp.SetContentType(CONTENTTYPE_HTML)
	return resp
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
// parameter expire is time by second
func newCookie(name, value string, expire int) string {
	return (&http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: expire,
	}).String()
}

// SetCookie setup response cookie, default age is default browser opened time
func (resp *Response) SetCookie(name, value string) {
	resp.SetCookieWithExpire(name, value, 0)
}

// SetCookieWithExpire setup response cookie with expire
func (resp *Response) SetCookieWithExpire(name, value string, expire int) {
	resp.SetHeader(HEADER_SETCOOKIE, newCookie(name, value, expire))
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
	http.Redirect(resp, resp.request, url, http.StatusTemporaryRedirect)
}

// PermanentRedirect permanently redirect current request url to new url
func (resp *Response) PermanentRedirect(url string) {
	http.Redirect(resp, resp.request, url, http.StatusMovedPermanently)
}

// Render render data with a template, and setup content type to html
func (resp *Response) Render(tmplName string, val interface{}) error {
	resp.SetContentType(CONTENTTYPE_HTML)
	return resp.server.RenderTemplate(resp, tmplName, val)
}

// WriteString write sting to client
func (resp *Response) WriteString(str string) (err error) {
	_, err = io.WriteString(resp, str)
	return
}

// WriteJSON write json data to client, and setup content type to json
func (resp *Response) WriteJSON(val interface{}) error {
	return resp.marshalValue(CONTENTTYPE_JSON, json.Marshal, val)
}

// WriteXML write xml data to client, and setup content type to xml
func (resp *Response) WriteXML(val interface{}) error {
	return resp.marshalValue(CONTENTTYPE_XML, xml.Marshal, val)
}

// marshalValue marshal value, and write it to client, setup response's content
// type to given format
func (resp *Response) marshalValue(format string, marshalFunc marshalFunc,
	val interface{}) error {

	bs, err := marshalFunc(val)
	if err == nil {
		if _, err = resp.Write(bs); err == nil {
			resp.SetContentType(format)
		}
	}
	return err
}
