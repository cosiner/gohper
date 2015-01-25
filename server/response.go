package server

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/cosiner/golib/types"
)

//==============================================================================
//                           Respone
//==============================================================================
const (
	CONTNTTYPE_PLAIN = "text/plain"
	CONTENTTYPE_HTML = "text/html"
	CONTENTTYPE_XML  = "text/xml"
	CONTENTTYPE_JSON = "text/json"
)

func parseContentType(str string) string {
	return types.TrimLower(str)
}

type Response struct {
	server *Server
	http.ResponseWriter
	header http.Header
}

func newResponse(s *Server, w http.ResponseWriter) *Response {
	resp := &Response{
		server:         s,
		ResponseWriter: w,
		header:         w.Header(),
	}
	resp.SetContentType(CONTENTTYPE_HTML)
	return resp
}

func (resp *Response) SetHeader(name, value string) {
	resp.header.Set(name, value)
}

func (resp *Response) AddHeader(name, value string) {
	resp.header.Add(name, value)
}

func (resp *Response) SetContentType(typ string) {
	resp.SetHeader(HEADER_CONTENTTYPE, typ)
}

func (resp *Response) contentType() string {
	return resp.header.Get(HEADER_CONTENTTYPE)
}

func newCookie(name, value string) string {
	return (&http.Cookie{Name: name, Value: value}).String()
}

func (resp *Response) SetCookie(name, value string) {
	resp.SetHeader(HEADER_SETCOOKIE, newCookie(name, value))
}

func (resp *Response) setSessionCookie(id string) {
	resp.SetCookie(_COOKIE_SESSION, id)
}

func (resp *Response) Redirect(req *Request, url string) {
	http.Redirect(resp, req.Request, url, http.StatusTemporaryRedirect)
}

func (resp *Response) PermanentRedirect(req *Request, url string) {
	http.Redirect(resp, req.Request, url, http.StatusMovedPermanently)
}

func (resp *Response) Render(tmplName string, val interface{}) error {
	resp.SetContentType(CONTENTTYPE_HTML)
	return resp.server.RenderTemplate(resp, tmplName, val)
}

func (resp *Response) WriteString(str string) (err error) {
	_, err = io.WriteString(resp, str)
	return
}

func (resp *Response) WriteJSON(val interface{}) error {
	return resp.marshalValue(CONTENTTYPE_JSON, json.Marshal, val)
}

func (resp *Response) WriteXML(val interface{}) error {
	return resp.marshalValue(CONTENTTYPE_XML, xml.Marshal, val)
}

type marshalFunc func(interface{}) ([]byte, error)

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
