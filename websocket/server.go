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

// Upgrade hajack a normal http connection, convert it to websocket connection
func Upgrade(w http.ResponseWriter, req *http.Request,
	handshakeChecker HandshakeChecker) (conn *Conn, err error) {
	var hs serverHandshaker = new(hybiServerHandshaker)
	rwc, buf, err := w.(http.Hijacker).Hijack()

	code, err := hs.ReadHandshake(buf.Reader, req)
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
