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
		encoding.CleanPowerWriter
		URL() *url.URL
	}

	// webSocketConn is the actual websocket connection
	webSocketConn struct {
		*websocket.Conn
		encoding.CleanPowerWriter
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
		Conn:             conn,
		CleanPowerWriter: encoding.NewPowerReadWriterInOne(conn),
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
