// Package crypto implements a encrypt algorithm using sha256,
// hex tool and random data generator
// encrypt algorithm: give message and fix salt,
// return sha256(msg, sha256(fixSalt + randSalt)) as result
package crypto

import (
	"crypto/sha256"
)

// SaltBits is the bit count of random generated Salt
const SaltBits = sha256.Size

// sha is a sha256 encryptor
var sha = sha256.New()

// ShaEncrypt use sha256 to encrypt message
func ShaEncrypt(msg, fixSalt string) (dst []byte, randSalt []byte, err error) {
	randSalt, err = RandBytesAlphanumeric(SaltBits)
	if err == nil {
		dst, err = ShaEncryptWithSalt(msg, fixSalt, randSalt)
	}
	return
}

// ShaEncryptWithSalt use gived salt as random salt
// sha256(sha256(msg) + sha256(fixSalt + salt)), randomString is SaltBits length
func ShaEncryptWithSalt(msg, fixSalt string, salt []byte) (dst []byte, err error) {
	sha.Reset()
	sha.Write([]byte(msg))
	encMsg := sha.Sum(nil)
	sha.Reset()
	sha.Write([]byte(fixSalt))
	sha.Write(salt)
	newSalt := sha.Sum(nil)
	sha.Reset()
	sha.Write(encMsg)
	sha.Write(newSalt)
	dst = sha.Sum(nil)

	return
}
