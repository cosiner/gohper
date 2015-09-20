package session

import (
	"crypto/md5"
	"testing"
	"time"

	"github.com/cosiner/gohper/encoding"
	"github.com/cosiner/gohper/testing2"
	"github.com/cosiner/gohper/unsafe2"
)

func TestCipher(t *testing.T) {
	tt := testing2.Wrap(t)
	c := NewCipher([]byte("12345"), time.Second*100, md5.New, encoding.Hex{})

	for _, s := range []string{"a", "b", "c", "d", ""} {
		tok := c.Encode(unsafe2.Bytes(s))
		tt.Log(string(tok))
		ds, err := c.Decode(tok)
		tt.DeepEq(unsafe2.Bytes(s), ds).Nil(err)
	}

	tt.Log()

	for _, s := range []string{"a", "b", "c", "d", ""} {
		tok := c.Encode(unsafe2.Bytes(s))
		ds, err := c.Decode(tok)
		tt.DeepEq(unsafe2.Bytes(s), ds).Nil(err)
	}
}
