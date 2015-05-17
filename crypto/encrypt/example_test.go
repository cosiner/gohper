package encrypt_test

import (
	"fmt"
	"log"

	"github.com/cosiner/gohper/crypto/encrypt"
)

func ExampleSaltEncode() {
	password := []byte("abcdefg")
	salt := []byte("123456")
	enc, randSalt, _ := encrypt.Encode(password, salt)

	if encrypt.Verify(password, salt, randSalt, enc) {
		fmt.Println("Password Match")
	} else {
		log.Fatalln("Password don't match")
	}
	// Output: Password Match
}
