package encrypt_test

import (
	"fmt"

	"github.com/cosiner/gohper/crypto/encrypt"
)

func ExampleSaltEncode() {
	password := []byte("abcdefg")
	salt := []byte("123456")
	enc, randSalt, _ := encrypt.Encode(password, salt)

	if encrypt.Verify(password, salt, randSalt, enc) {
		fmt.Println("Password Match")
	} else {
		fmt.Println("Password don't match")
	}
	// Output: Password Match
}
