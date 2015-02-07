package server

import (
	"net"
	"net/url"

	"github.com/cosiner/golib/encoding"

	"github.com/cosiner/gomodule/websocket"
)

type (
	// webSocketConn represent an websocket connection
	// WebSocket connection is not be managed in server,
	// it's handler's responsibility to close connection
	WebSocketConn interface {
		UrlVarIndexer
		net.Conn
		URL() *url.URL
		WriteString(data string) (int, error)
		WriteJSON(v interface{}) error
		WriteXML(v interface{}) error
		ReadString() (string, error)
		ReadJSON(v interface{}) error
		ReadXML(v interface{}) error
	}

	// webSocketConn is the actual websocket connection
	webSocketConn struct {
		*websocket.Conn
		UrlVarIndexer
	}

	// WebSocketHandlerFunc is the websocket connection handler
	WebSocketHandlerFunc func(WebSocketConn)

	// WebSocketHandler is the handler of websocket connection
	WebSocketHandler interface {
		Init(*Server) error
		Destroy()
		Handle(WebSocketConn)
	}
)

// WebSocketHandlerFunc is a function WebSocketHandler
func (WebSocketHandlerFunc) Init(*Server) error           { return nil }
func (fn WebSocketHandlerFunc) Handle(conn WebSocketConn) { fn(conn) }
func (WebSocketHandlerFunc) Destroy()                     {}

// newWebSocketConn wrap a exist websocket connection and url variables to a
// new webSocketConn
func newWebSocketConn(conn *websocket.Conn) *webSocketConn {
	return &webSocketConn{
		Conn: conn,
	}
}

// setUrlVarIndexer set up url variable vlaues
func (wsc *webSocketConn) setUrlVarIndexer(indexer UrlVarIndexer) *webSocketConn {
	wsc.UrlVarIndexer = indexer
	return wsc
}

// URL return client side url
func (wsc *webSocketConn) URL() *url.URL {
	return wsc.Config().Origin
}

// WriteString write string to client side
func (wsc *webSocketConn) WriteString(data string) (int, error) {
	return encoding.WriteString(wsc, data)
}

// WriteString write JSON data to client side
func (wsc *webSocketConn) WriteJSON(v interface{}) error {
	return encoding.WriteJSON(wsc, v)
}

// WriteString write XML data to client side
func (wsc *webSocketConn) WriteXML(v interface{}) error {
	return encoding.WriteXML(wsc, v)
}

// ReadString read recieved data as string
func (wsc *webSocketConn) ReadString() (string, error) {
	return encoding.ReadString(wsc)
}

// ReadString read recieved data as json data
func (wsc *webSocketConn) ReadJSON(v interface{}) error {
	return encoding.ReadJSON(wsc, v)
}

// ReadString read recieved data as xml data
func (wsc *webSocketConn) ReadXML(v interface{}) error {
	return encoding.ReadXML(wsc, v)
}
