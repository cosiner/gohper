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
		urlVars []string
		indexer VarIndexer
	}

	// WebSocketHandlerFunc is the websocket connection handler
	WebSocketHandlerFunc func(*WebSocketConn)

	// WebSocketHandler is the handler of websocket connection
	WebSocketHandler interface {
		Init(*Server) error
		Destroy()
		Handle(*WebSocketConn)
	}
)

// WebSocketHandlerFunc is a function WebSocketHandler
func (WebSocketHandlerFunc) Init(*Server) error            { return nil }
func (fn WebSocketHandlerFunc) Handle(conn *WebSocketConn) { fn(conn) }
func (WebSocketHandlerFunc) Destroy()                      {}

// newWebSocketConn wrap a exist websocket connection and url variables to a
// new WebSocketConn
func newWebSocketConn(conn *websocket.Conn) *WebSocketConn {
	return &WebSocketConn{
		Conn: conn,
	}
}

// UrlVar return variable value in url
func (wsc *WebSocketConn) UrlVar(name string) (value string) {
	if values := wsc.urlVars; values != nil {
		value = wsc.indexer.ValueOf(values, name)
	}
	return
}

// ScanUrlVars scan url variable values into given address
func (wsc *WebSocketConn) ScanUrlVars(vars ...*string) {
	if values := wsc.urlVars; values != nil {
		wsc.indexer.ScanInto(values, vars...)
	}
}

// setUrlVars set up url variable vlaues
func (wsc *WebSocketConn) setUrlVars(indexer VarIndexer, vars []string) *WebSocketConn {
	wsc.indexer = indexer
	wsc.urlVars = vars
	return wsc
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
