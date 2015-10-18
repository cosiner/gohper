package url2

import (
	"bytes"
	"net/url"
	"strconv"
	"strings"

	"github.com/cosiner/gohper/slices"
)

var Bufsize = 128
var emptyBytes = []byte("")

// Query parameters to url query string without escape,
// if buf is not nil, and there is more than one parameter, the allocated buffer
// will stored to *buf
func Query(params map[string]string, buf *bytes.Buffer) ([]byte, bool) {
	l := len(params)
	if l == 0 {
		return emptyBytes, false
	}

	if l == 1 {
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

// QueryEscape is same as Query, but escape the query string
func QueryEscape(params map[string]string, buf *bytes.Buffer) ([]byte, bool) {
	s, b := Query(params, buf)
	if len(s) == 0 {
		return nil, b
	}

	return []byte(url.QueryEscape(string(s))), false // TODO: remove bytes convert
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
		return slices.Ints(v).Join("", ",")
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case []uint:
		return slices.Uints(v).Join("", ",")
	case []byte:
		return string(v)
	}

	panic("Only support int, []int, uint, []uint, string, []string, []byte")
}
