package session

import (
	"bytes"
	"crypto/hmac"
	"encoding/binary"
	"hash"
	"time"

	"github.com/cosiner/gohper/encoding"
	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/time2"
)

const (
	ErrBadKey           = errors.Err("bad key")
	ErrExpiredKey       = errors.Err("expired key")
	ErrInvalidSignature = errors.Err("invalid signature")
)

func NewCipher(signKey []byte, ttl time.Duration, hash func() hash.Hash, encs ...encoding.Encoding) encoding.Encoding {
	return encoding.Pipe(encs).Prepend(&Cipher{
		signKey:      signKey,
		ttl:          ttl,
		hash:         hash,
		sigLen:       hash().Size() / 2,
		timestampLen: 8,
	})
}

type Cipher struct {
	signKey      []byte
	ttl          time.Duration
	hash         func() hash.Hash
	sigLen       int
	timestampLen int
}

func (c *Cipher) encrypt(now uint64, str []byte) []byte {
	hash := hmac.New(c.hash, c.signKey)

	hash.Write(str)
	timestamp := make([]byte, c.timestampLen)
	binary.BigEndian.PutUint64(timestamp, now)
	hash.Write(c.signKey)

	sig := hash.Sum(nil)[:c.sigLen]
	result := make([]byte, c.sigLen+c.timestampLen+len(str))
	copy(result, sig)
	copy(result[c.sigLen:], timestamp)
	copy(result[c.sigLen+c.timestampLen:], str)

	return result
}

func (c *Cipher) Encode(str []byte) []byte {
	n := uint64(time2.Now().Add(c.ttl).UnixNano())
	return c.encrypt(n, str)
}

func (c *Cipher) Decode(str []byte) ([]byte, error) {
	hdrLen := c.sigLen + c.timestampLen
	if len(str) < hdrLen {
		return nil, ErrBadKey
	}

	tm := binary.BigEndian.Uint64(str[c.sigLen:])
	if c.ttl != 0 && uint64(time2.Now().UnixNano()) > tm {
		return nil, ErrExpiredKey
	}

	data := str[hdrLen:]
	sig := c.encrypt(tm, data)
	if !bytes.Equal(sig, str) {
		return nil, ErrInvalidSignature
	}

	return data, nil
}
