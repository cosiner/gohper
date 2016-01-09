package http2

import "strings"

func IpOfAddr(addr string) string {
	i := strings.IndexByte(addr, ':')
	if i >= 0 {
		addr = addr[:i]
	}
	return addr
}
