package server

import (
	"net/url"

	"github.com/cosiner/golib/encoding"

	"github.com/cosiner/gomodule/websocket"
)

type (
	// WebSocketConn represent an websocket connection
	// WebSocket connection is not be managed in server,
	// it's handler's responsibility to close connection
	WebSocketConn struct {
		*websocket.Conn
		urlVars map[string]string
	}

	// WebSocketHandler is the handler of websocket connection
	WebSocketHandler interface {
		Init(*Server) error
		Destroy()
		Handle(*WebSocketConn)
	}
)

// newWebSocketConn wrap a exist websocket connection and url variables to a
// new WebSocketConn
func newWebSocketConn(conn *websocket.Conn, urlVars map[string]string) *WebSocketConn {
	return &WebSocketConn{
		Conn:    conn,
		urlVars: urlVars,
	}
}

// UrlVar return variable value in url
func (wsc *WebSocketConn) UrlVar(name string) string {
	return wsc.urlVars[name]
}

// setUrlVars set up url variable vlaues
func (wsc *WebSocketConn) setUrlVars(vars map[string]string) {
	wsc.urlVars = vars
}

// URL return client side url
func (wsc *WebSocketConn) URL() *url.URL {
	return wsc.Config().Origin
}

// WriteString write string to client side
func (wsc *WebSocketConn) WriteString(data string) (int, error) {
	return encoding.WriteString(wsc, data)
}

// WriteString write JSON data to client side
func (wsc *WebSocketConn) WriteJSON(v interface{}) error {
	return encoding.WriteJSON(wsc, v)
}

// WriteString write XML data to client side
func (wsc *WebSocketConn) WriteXML(v interface{}) error {
	return encoding.WriteXML(wsc, v)
}

// ReadString read recieved data as string
func (wsc *WebSocketConn) ReadString() (string, error) {
	return encoding.ReadString(wsc)
}

// ReadString read recieved data as json data
func (wsc *WebSocketConn) ReadJSON(v interface{}) error {
	return encoding.ReadJSON(wsc, v)
}

// ReadString read recieved data as xml data
func (wsc *WebSocketConn) ReadXML(v interface{}) error {
	return encoding.ReadXML(wsc, v)
}
