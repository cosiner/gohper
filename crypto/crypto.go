// Package crypto implements a encrypt algorithm using sha256,
// hex tool and random data generator
// encrypt algorithm: give message and fix salt,
// return sha256(msg, sha256(fixSalt + randSalt)) as result
package crypto

import "crypto/sha256"

// SaltBits is the bit count of random generated Salt
const SaltBits = sha256.Size

// ShaEncrypt use sha256 to encrypt message
func ShaEncrypt(msg, fixSalt []byte) (dst []byte, randSalt []byte, err error) {
	randSalt, err = RandBytesAlphanumeric(SaltBits)
	if err == nil {
		dst = ShaEncryptWithSalt(msg, fixSalt, randSalt)
	}
	return
}

// ShaEncryptWithSalt use gived salt as random salt
// sha256(sha256(msg) + sha256(fixSalt + salt)), randomString is SaltBits length
func ShaEncryptWithSalt(msg, fixSalt, salt []byte) []byte {
	var sha = sha256.New()
	sha.Write(msg)
	encMsg := sha.Sum(nil)
	sha.Reset()
	sha.Write(fixSalt)
	sha.Write(salt)
	newSalt := sha.Sum(nil)
	sha.Reset()
	sha.Write(encMsg)
	sha.Write(newSalt)
	return sha.Sum(nil)
}
