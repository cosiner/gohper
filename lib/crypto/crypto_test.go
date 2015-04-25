package crypto

import (
	"testing"
)

func TestEnCrypt(t *testing.T) {
	password := []byte("abcdefg")
	fixsalt := []byte("123456")
	enc, salt, _ := ShaEncrypt(password, fixsalt)
	newEnc := ShaEncryptWithSalt(password, fixsalt, salt)
	if string(newEnc) != string(enc) {
		t.Fail()
	}
}
