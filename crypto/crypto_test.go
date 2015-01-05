package test

import (
	"github.com/cosiner/golib/crypto"
	"testing"
)

func TestEnCrypt(t *testing.T) {
	password := "abcdefg"
	fixsalt := "123456"
	enc, salt, err := crypto.ShaEncrypt(password, fixsalt)

	if err == nil {
		t.Log(crypto.BytesToHexStr(enc), len(enc)*2, crypto.BytesToHexStr(salt), len(salt)*2)
	}
}
