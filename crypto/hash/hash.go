package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

func Hash(hashNew func() hash.Hash, key, data []byte, tohex bool) []byte {
	var h hash.Hash
	if len(key) != 0 {
		h = hmac.New(hashNew, key)
	} else {
		h = hashNew()
	}
	h.Write(data)
	sum := h.Sum(nil)
	if tohex {
		hexSum := make([]byte, hex.EncodedLen(len(sum)))
		n := hex.Encode(hexSum, sum)
		sum = hexSum[:n]
	}

	return sum
}

func MD5(key, data []byte, tohex bool) []byte {
	return Hash(md5.New, key, data, tohex)
}

func SHA1(key, data []byte, tohex bool) []byte {
	return Hash(sha1.New, key, data, tohex)
}

func SHA256(key, data []byte, tohex bool) []byte {
	return Hash(sha256.New, key, data, tohex)
}

func SHA512(key, data []byte, tohex bool) []byte {
	return Hash(sha512.New, key, data, tohex)
}
