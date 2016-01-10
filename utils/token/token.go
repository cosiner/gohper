package token

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

type Cipher struct {
	signKey []byte
	ttl     time.Duration
	hash    func() hash.Hash
	sigLen  int
	hdrLen  int
}

func NewCipher(signKey []byte, ttl time.Duration, hash func() hash.Hash, encs ...encoding.Encoding) encoding.Encoding {
	sigLen := hash().Size() / 3
	return encoding.Pipe(encs).Prepend(&Cipher{
		signKey: signKey,
		ttl:     ttl,
		hash:    hash,
		sigLen:  sigLen,
		hdrLen:  sigLen + 8,
	})
}

// | signature | deadline | str
func (c *Cipher) encrypt(deadline uint64, b []byte) []byte {
	result := make([]byte, c.hdrLen+len(b))
	binary.BigEndian.PutUint64(result[c.sigLen:c.hdrLen], deadline)
	copy(result[c.hdrLen:], b)

	hash := hmac.New(c.hash, c.signKey)
	hash.Write(b)
	hash.Write(result[c.sigLen : c.hdrLen])
	copy(result, hash.Sum(nil)[:c.sigLen])

	return result
}

func (c *Cipher) Encode(b []byte) []byte {
	deadline := uint64(time2.Now().Add(c.ttl).Unix())
	return c.encrypt(deadline, b)
}

func (c *Cipher) Decode(b []byte) ([]byte, error) {
	if len(b) < c.hdrLen {
		return nil, ErrBadKey
	}

	deadline := binary.BigEndian.Uint64(b[c.sigLen:c.hdrLen])
	if c.ttl != 0 && uint64(time2.Now().Unix()) > deadline {
		return nil, ErrExpiredKey
	}

	data := b[c.hdrLen:]
	encData := c.encrypt(deadline, data)
	if !bytes.Equal(encData, b) {
		return nil, ErrInvalidSignature
	}

	return data, nil
}
