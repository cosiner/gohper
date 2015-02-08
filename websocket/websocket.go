// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package websocket implements a client and server for the WebSocket protocol
// as specified in RFC 6455.
package websocket

import (
	"bufio"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	ProtocolVersionHybi13    = 13
	ProtocolVersionHybi      = ProtocolVersionHybi13
	SupportedProtocolVersion = "13"

	ContinuationFrame = 0
	TextFrame         = 1
	BinaryFrame       = 2
	CloseFrame        = 8
	PingFrame         = 9
	PongFrame         = 10
	UnknownFrame      = 255
)

// ProtocolError represents WebSocket protocol errors.
type ProtocolError struct {
	ErrorString string
}

func (err *ProtocolError) Error() string { return err.ErrorString }

var (
	ErrBadProtocolVersion   = &ProtocolError{"bad protocol version"}
	ErrBadScheme            = &ProtocolError{"bad scheme"}
	ErrBadStatus            = &ProtocolError{"bad status"}
	ErrBadUpgrade           = &ProtocolError{"missing or bad upgrade"}
	ErrBadWebSocketOrigin   = &ProtocolError{"missing or bad WebSocket-Origin"}
	ErrBadWebSocketLocation = &ProtocolError{"missing or bad WebSocket-Location"}
	ErrBadWebSocketProtocol = &ProtocolError{"missing or bad WebSocket-Protocol"}
	ErrBadWebSocketVersion  = &ProtocolError{"missing or bad WebSocket Version"}
	ErrChallengeResponse    = &ProtocolError{"mismatch challenge/response"}
	ErrBadFrame             = &ProtocolError{"bad frame"}
	ErrBadFrameBoundary     = &ProtocolError{"not on frame boundary"}
	ErrNotWebSocket         = &ProtocolError{"not websocket protocol"}
	ErrBadRequestMethod     = &ProtocolError{"bad method"}
	ErrNotSupported         = &ProtocolError{"not supported"}
)

// Addr is an implementation of net.Addr for WebSocket.
type Addr struct {
	*url.URL
}

// Network returns the network type for a WebSocket, "websocket".
func (addr *Addr) Network() string { return "websocket" }

// Config is a WebSocket configuration
type Config struct {
	// A WebSocket server address.
	Location *url.URL

	// A Websocket client origin.
	Origin *url.URL

	// WebSocket subprotocols.
	Protocol []string

	// WebSocket protocol version.
	Version int

	// TLS config for secure WebSocket (wss).
	TlsConfig *tls.Config

	// Additional header fields to be sent in WebSocket opening handshake.
	Header http.Header

	handshakeData map[string]string
}

// serverHandshaker is an interface to handle WebSocket server side handshake.
type serverHandshaker interface {
	// ReadHandshake reads handshake request message from client.
	// Returns http response code and error if any.
	ReadHandshake(buf *bufio.Reader, req *http.Request, checkedWebsocket bool) (code int, err error)

	// AcceptHandshake accepts the client handshake request and sends
	// handshake response back to client.
	AcceptHandshake(buf *bufio.Writer) (err error)

	// NewServerConn creates a new WebSocket connection.
	NewServerConn(buf *bufio.ReadWriter, rwc io.ReadWriteCloser, request *http.Request) (conn *Conn)
}

// frameReader is an interface to read a WebSocket frame.
type frameReader interface {
	// Reader is to read payload of the frame.
	io.Reader

	// PayloadType returns payload type.
	PayloadType() byte

	// HeaderReader returns a reader to read header of the frame.
	HeaderReader() io.Reader

	// TrailerReader returns a reader to read trailer of the frame.
	// If it returns nil, there is no trailer in the frame.
	TrailerReader() io.Reader

	// Len returns total length of the frame, including header and trailer.
	Len() int
}

// frameReaderFactory is an interface to creates new frame reader.
type frameReaderFactory interface {
	NewFrameReader() (r frameReader, err error)
}

// frameWriter is an interface to write a WebSocket frame.
type frameWriter interface {
	// Writer is to write payload of the frame.
	io.WriteCloser
}

// frameWriterFactory is an interface to create new frame writer.
type frameWriterFactory interface {
	NewFrameWriter(payloadType byte) (w frameWriter, err error)
}

type frameHandler interface {
	HandleFrame(frame frameReader) (r frameReader, err error)
	WriteClose(status int) (err error)
}

// Conn represents a WebSocket connection.
type Conn struct {
	config  *Config
	request *http.Request

	buf *bufio.ReadWriter
	rwc io.ReadWriteCloser

	rio sync.Mutex
	frameReaderFactory
	frameReader

	wio sync.Mutex
	frameWriterFactory

	frameHandler
	PayloadType        byte
	defaultCloseStatus int
}

// Read implements the io.Reader interface:
// it reads data of a frame from the WebSocket connection.
// if msg is not large enough for the frame data, it fills the msg and next Read
// will read the rest of the frame data.
// it reads Text frame or Binary frame.
func (ws *Conn) Read(msg []byte) (n int, err error) {
	ws.rio.Lock()
	defer ws.rio.Unlock()
again:
	if ws.frameReader == nil {
		frame, err := ws.frameReaderFactory.NewFrameReader()
		if err != nil {
			return 0, err
		}
		ws.frameReader, err = ws.frameHandler.HandleFrame(frame)
		if err != nil {
			return 0, err
		}
		if ws.frameReader == nil {
			goto again
		}
	}
	n, err = ws.frameReader.Read(msg)
	if err == io.EOF {
		if trailer := ws.frameReader.TrailerReader(); trailer != nil {
			io.Copy(ioutil.Discard, trailer)
		}
		ws.frameReader = nil
		goto again
	}
	return n, err
}

// Write implements the io.Writer interface:
// it writes data as a frame to the WebSocket connection.
func (ws *Conn) Write(msg []byte) (n int, err error) {
	ws.wio.Lock()
	defer ws.wio.Unlock()
	w, err := ws.frameWriterFactory.NewFrameWriter(ws.PayloadType)
	if err != nil {
		return 0, err
	}
	n, err = w.Write(msg)
	w.Close()
	if err != nil {
		return n, err
	}
	return n, err
}

// Close implements the io.Closer interface.
func (ws *Conn) Close() error {
	err := ws.frameHandler.WriteClose(ws.defaultCloseStatus)
	if err != nil {
		return err
	}
	return ws.rwc.Close()
}
func (ws *Conn) isClientConn() bool { return ws.request == nil }
func (ws *Conn) isServerConn() bool { return ws.request != nil }

// LocalAddr returns the WebSocket Origin for the connection for client, or
// the WebSocket location for server.
func (ws *Conn) LocalAddr() net.Addr {
	if ws.isClientConn() {
		return &Addr{ws.config.Origin}
	}
	return &Addr{ws.config.Location}
}

// RemoteAddr returns the WebSocket location for the connection for client, or
// the Websocket Origin for server.
func (ws *Conn) RemoteAddr() net.Addr {
	if ws.isClientConn() {
		return &Addr{ws.config.Location}
	}
	return &Addr{ws.config.Origin}
}

func (ws *Conn) Config() *Config {
	return ws.config
}

var errSetDeadline = errors.New("websocket: cannot set deadline: not using a net.Conn")

// SetDeadline sets the connection's network read & write deadlines.
func (ws *Conn) SetDeadline(t time.Time) error {
	if conn, ok := ws.rwc.(net.Conn); ok {
		return conn.SetDeadline(t)
	}
	return errSetDeadline
}

// SetReadDeadline sets the connection's network read deadline.
func (ws *Conn) SetReadDeadline(t time.Time) error {
	if conn, ok := ws.rwc.(net.Conn); ok {
		return conn.SetReadDeadline(t)
	}
	return errSetDeadline
}

// SetWriteDeadline sets the connection's network write deadline.
func (ws *Conn) SetWriteDeadline(t time.Time) error {
	if conn, ok := ws.rwc.(net.Conn); ok {
		return conn.SetWriteDeadline(t)
	}
	return errSetDeadline
}
