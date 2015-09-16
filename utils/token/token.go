package session

import (
	"crypto/hmac"
	"encoding/hex"
	"hash"
	"strconv"
	"strings"
	"time"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/time2"
	"github.com/cosiner/gohper/unsafe2"
)

const (
	ErrBadKey           = errors.Err("bad key")
	ErrExpiredKey       = errors.Err("expired key")
	ErrInvalidSignature = errors.Err("invalid signature")
)

type Cipher struct {
	SecretKey string
	TTL       time.Duration
	Hash      func() hash.Hash
	Seperator string
}

func (c *Cipher) segSep() string {
	if c.Seperator == "" {
		return "."
	}

	return c.Seperator
}

func (c *Cipher) encrypt(now int64, str string) string {
	hash := hmac.New(c.Hash, unsafe2.Bytes(c.SecretKey))

	if now == 0 {
		now = time2.Now().Add(c.TTL).UnixNano()
	}
	hash.Write(unsafe2.Bytes(str))
	hash.Write(unsafe2.Bytes(c.segSep()))
	nows := strconv.FormatInt(now, 10)
	hash.Write(unsafe2.Bytes(nows))
	hash.Write(unsafe2.Bytes(c.segSep()))
	hash.Write(unsafe2.Bytes(c.SecretKey))
	sig := hash.Sum(nil)

	sigStr := hex.EncodeToString(sig[:hash.Size()/2])
	return str + c.segSep() + nows + c.segSep() + sigStr[:len(sigStr)/2]
}

func (c *Cipher) Encrypt(str string) string {
	return c.encrypt(0, str)
}

func (c *Cipher) Decrypt(str string) (string, error) {
	segs := strings.Split(str, c.segSep())
	if len(segs) != 3 {
		return "", ErrBadKey
	}

	tm, err := strconv.ParseInt(segs[1], 10, 64)
	if err != nil {
		return "", ErrBadKey
	}
	if time2.Now().UnixNano() > tm {
		return "", ErrExpiredKey
	}

	sig := c.encrypt(tm, segs[0])
	if sig != str {
		return "", ErrInvalidSignature
	}

	return segs[0], nil
}
