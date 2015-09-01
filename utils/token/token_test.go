package session

import (
	"crypto/md5"
	"testing"
	"time"

	"github.com/cosiner/gohper/testing2"
)

func TestCipher(t *testing.T) {
	tt := testing2.Wrap(t)
	c := Cipher{
		SecretKey: "123456dddddddddddddddddddddddddddddddddddddddddddddddddddddd",
		TTL:       time.Second * 3600,
		Hash:      md5.New,
		Seperator: ".",
	}

	for _, s := range []string{"a", "b", "c", "d", ""} {
		tok := c.Encrypt(s)
		ds, err := c.Decrypt(tok)
		tt.Eq(s, ds).Nil(err)
	}
}
