package url2

import (
	"bytes"
	"net/url"
	"strconv"
	"strings"

	"github.com/cosiner/gohper/strings2"
)

var Bufsize = 128
var emptyBytes = []byte("")

// Encode parameters to url query string without escape,
// if buf is not nil, and there is more than one parameter, the allocated buffer
// will stored to *buf
func Encode(params map[string]string, buf *bytes.Buffer) ([]byte, bool) {
	if l := len(params); l == 0 {
		return emptyBytes, false
	} else if l == 1 {
		for k, v := range params {
			return []byte(k + "=" + v), false
		}
	}
	var nbuf = buf
	if buf == nil {
		nbuf = bytes.NewBuffer(make([]byte, 0, Bufsize))
	}
	var i int
	for k, v := range params {
		if i != 0 {
			nbuf.WriteByte('&')
		}
		i++
		nbuf.WriteString(k)
		nbuf.WriteByte('=')
		nbuf.WriteString(v)
	}
	return nbuf.Bytes(), buf != nil
}

// EscapeEncode is same as Encode, but escape the query string
func EscapeEncode(params map[string]string, buf *bytes.Buffer) ([]byte, bool) {
	s, b := Encode(params, buf)
	if len(s) != 0 {
		return []byte(url.QueryEscape(string(s))), false // TODO: remove bytes convert
	}
	return s, b
}

// Param convert i to params string, join with ",", without escape
//
// Only support int, []int, uint, []uint, string, []string, []byte
func Param(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v, ",")
	case int:
		return strconv.Itoa(v)
	case []int:
		return strings2.JoinInt(v, ",")
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case []uint:
		return strings2.JoinUint(v, ",")
	case []byte:
		return string(v)
	}
	panic("Only support int, []int, uint, []uint, string, []string, []byte")
}
