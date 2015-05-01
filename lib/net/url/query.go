package url

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/cosiner/gohper/lib/types"

	"bytes"
)

var Bufsize = 128

// Query encode parameters to url query string without escape,
// if buf is not nil, and there is more than one parameter, the allocated buffer
// will stored to *buf
func Query(params map[string]string, buf **bytes.Buffer) string {
	if l := len(params); l == 0 {
		return ""
	} else if l == 1 {
		for k, v := range params {
			return k + "=" + v
		}
	}
	var nbuf *bytes.Buffer
	if buf == nil {
		nbuf = bytes.NewBuffer(make([]byte, 0, Bufsize))
	} else if *buf == nil {
		*buf = bytes.NewBuffer(make([]byte, 0, Bufsize))
		nbuf = *buf
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
	return nbuf.String()
}

// QueryEscape is same as Query, but escape the query string
func QueryEscape(params map[string]string, buf **bytes.Buffer) string {
	s := Query(params, buf)
	if s != "" {
		s = url.QueryEscape(s)
	}
	return s
}

// ParamString convert i to params string, join with ",",
// Only support int, []int, uint, []uint, string, []string, []byte
func ParamString(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v, ",")
	case int:
		return strconv.Itoa(v)
	case []int:
		return types.JoinInt(v, ",")
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case []uint:
		return types.JoinUint(v, ",")
	case []byte:
		return string(v)
	}
	panic("Only support int, []int, uint, []uint, string, []string, []byte")
}
