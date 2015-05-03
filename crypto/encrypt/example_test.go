package encrypt_test

import (
	"fmt"
	"log"
	"reflect"

	"github.com/cosiner/gohper/crypto/encrypt"
)

func ExampleSaltEncode() {
	password := []byte("abcdefg")
	fixsalt := []byte("123456")
	enc, salt, _ := encrypt.Encode(password, fixsalt)

	newEnc := encrypt.SaltEncode(password, fixsalt, salt)

	if reflect.DeepEqual(newEnc, enc) {
		fmt.Println("Password Match")
	} else {
		log.Fatalln("Password don't match")
	}
	// Output: Password Match
}
