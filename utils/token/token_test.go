package session

import (
	"crypto/md5"
	"testing"
	"time"
)

func BenchmarkTimingEncrypt(b *testing.B) {
	c := Cipher{
		SecretKey: "123456",
		TTL:       time.Millisecond * 1000,
		Hash:      md5.New,
		Timefmt:   "20060102150405",
		Seperator: ".",
	}

	key := c.Encrypt("abcdefg")
	b.Log(key)
	for i := 0; i < b.N; i++ {
		c.Decrypt(key)
	}
}
