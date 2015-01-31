package websocket

import (
	"fmt"
	"net/http"
)

// HandShakeCheck is checker for websocket handshake message
type HandshakeChecker func(*Config, *http.Request) error

// checkOrigin check request origin header
// a connection request without origin should be abort
func checkOrigin(config *Config, req *http.Request) (err error) {
	config.Origin, err = Origin(config, req)
	if err == nil && config.Origin == nil {
		return fmt.Errorf("null origin")
	}
	return err
}

// IsWebSocketRequest check whether a request is websocket request
func IsWebSocketRequest(req *http.Request) bool {
	_, err := isWebSocketRequest(req)
	return err == nil
}

// UpgradeWebsocket upgrade a http request  to websocket connection, the request must
// already have been checked and is actuall a websocket request
func UpgradeWebsocket(w http.ResponseWriter, req *http.Request, handshakeChecker HandshakeChecker) (conn *Conn, err error) {
	return upgrade(w, req, true, handshakeChecker)
}

// Upgrade is same as UpgradeWebsocket, but not required check websocket request
func Upgrade(w http.ResponseWriter, req *http.Request, handshakeChecker HandshakeChecker) (conn *Conn, err error) {
	return upgrade(w, req, false, handshakeChecker)
}

// upgrade hajack a normal http connection, convert it to websocket connection
func upgrade(w http.ResponseWriter, req *http.Request, checkedWebsocket bool,
	handshakeChecker HandshakeChecker) (conn *Conn, err error) {
	config := new(Config)
	var hs serverHandshaker = &hybiServerHandshaker{Config: config}
	rwc, buf, err := w.(http.Hijacker).Hijack()

	code, err := hs.ReadHandshake(buf.Reader, req, checkedWebsocket)
	if err == ErrBadWebSocketVersion {
		fmt.Fprintf(buf, "HTTP/1.1 %03d %s\r\n", code, http.StatusText(code))
		fmt.Fprintf(buf, "Sec-WebSocket-Version: %s\r\n", SupportedProtocolVersion)
		buf.WriteString("\r\n")
		buf.WriteString(err.Error())
		buf.Flush()
		return
	}
	if err != nil {
		fmt.Fprintf(buf, "HTTP/1.1 %03d %s\r\n", code, http.StatusText(code))
		buf.WriteString("\r\n")
		buf.WriteString(err.Error())
		buf.Flush()
		return
	}
	if handshakeChecker == nil {
		handshakeChecker = checkOrigin
	}
	if err = handshakeChecker(config, req); err != nil {
		code = http.StatusForbidden
		fmt.Fprintf(buf, "HTTP/1.1 %03d %s\r\n", code, http.StatusText(code))
		buf.WriteString("\r\n")
		buf.Flush()
		return
	}
	if err = hs.AcceptHandshake(buf.Writer); err != nil {
		code = http.StatusBadRequest
		fmt.Fprintf(buf, "HTTP/1.1 %03d %s\r\n", code, http.StatusText(code))
		buf.WriteString("\r\n")
		buf.Flush()
		return
	}
	conn = hs.NewServerConn(buf, rwc, req)
	return
}
