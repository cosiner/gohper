// Package encrypt implements a encrypt tool with fix salt and random salt.
//
// Algorithm:
//   Give message and fix salt, return Hash(msg, Hash(salt + randSalt)) as result.
package encrypt

import (
	"bytes"
	"crypto/sha256"
	"hash"

	"github.com/cosiner/gohper/crypto/rand"
)

// Encode a message with fixed salt, return the encoded message and random salt
func Encode(hash hash.Hash, msg, salt []byte) ([]byte, []byte, error) {
	if hash == nil {
		hash = sha256.New()
	}
	rand, err := rand.B.Alphanumeric(hash.BlockSize())
	if err != nil {
		return nil, nil, err
	}

	return SaltEncode(hash, msg, salt, rand), rand, err
}

// SaltEncode encode the message with a fixed salt and a random salt, typically
// used to verify
func SaltEncode(hash hash.Hash, msg, fixed, rand []byte) []byte {
	if hash == nil {
		hash = sha256.New()
	}

	hash.Write(msg)
	enc := hash.Sum(nil)

	hash.Reset()
	hash.Write(fixed)
	hash.Write(rand)
	new := hash.Sum(nil)

	hash.Reset()
	hash.Write(enc)
	hash.Write(new)

	return hash.Sum(nil)
}

// Verify will encode msg, salt, randSalt, then compare it with encoded password,
// return true if equal, else false
func Verify(hash hash.Hash, msg, salt, randSalt, pass []byte) bool {
	return bytes.Equal(SaltEncode(hash, msg, salt, randSalt), pass)
}
