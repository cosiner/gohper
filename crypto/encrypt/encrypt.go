// Package encrypt implements a encrypt tool with fix salt and random salt.
//
// Algorithm:
//   Give message and fix salt, return Hash(msg, Hash(salt + randSalt)) as result.
package encrypt

import (
	"bytes"
	"crypto/sha256"

	"github.com/cosiner/gohper/crypto/rand"
)

// Hash is the hash function used for Encode and SaltEncode
var Hash = sha256.New

// SaltSize is the size of generated salt in bytes used in Encode
var SaltSize = sha256.Size

// Encode a message with fixed salt, return the encoded message and random salt
func Encode(msg, salt []byte) ([]byte, []byte, error) {
	rand, err := rand.B.Alphanumeric(SaltSize)
	if err != nil {
		return nil, nil, err
	}

	return SaltEncode(msg, salt, rand), rand, err
}

// SaltEncode encode the message with a fixed salt and a random salt, typically
// used to verify
func SaltEncode(msg, fixed, rand []byte) []byte {
	h := Hash()

	h.Write(msg)
	enc := h.Sum(nil)

	h.Reset()
	h.Write(fixed)
	h.Write(rand)
	new := h.Sum(nil)

	h.Reset()
	h.Write(enc)
	h.Write(new)

	return h.Sum(nil)
}

// Verify will encode msg, salt, randSalt, then compare it with encoded password,
// return true if equal, else false
func Verify(msg, salt, randSalt, pass []byte) bool {
	return bytes.Equal(SaltEncode(msg, salt, randSalt), pass)
}
