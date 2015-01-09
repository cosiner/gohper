package crypto

import (
	"crypto/rand"
	"math/big"

	. "github.com/cosiner/golib/errors"
)

// Charset of characters to use for generating random strings
const (
	Alphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	Numerals     = "1234567890"
	Alphanumeric = Alphabet + Numerals
	Ascii        = Alphanumeric + "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
)

var (
	numNegativeError  = Err("Number cannot be negative")
	charsetEmptyError = Err("Charset cannot be empty")
)

// RandBytesInCharset generate a random bytes in gived charset
func RandBytesInCharset(n int, charset string) ([]byte, error) {
	if n <= 0 {
		return nil, numNegativeError
	}
	var l int
	if l = len(charset); l == 0 {
		return nil, charsetEmptyError
	}

	result := make([]byte, n)
	charnum := big.NewInt(int64(l))
	for i := 0; i < n; i++ {
		r, err := rand.Int(rand.Reader, charnum)
		if err != nil {
			return nil, err
		}
		result[i] = charset[int(r.Int64())]
	}
	return result, nil
}

// RandInCharset is same as RandBytesInCharset, return string
func RandInCharset(n int, charset string) (string, error) {
	b, err := RandBytesInCharset(n, charset)
	if err == nil {
		return string(b), nil
	}
	return "", err
}

// RandBytesAscii generate random ascii bytes
func RandBytesAscii(n int) ([]byte, error) {
	return RandBytesInCharset(n, Ascii)
}

// RandBytesNumberal generate random ascii bytes
func RandBytesNumberal(n int) ([]byte, error) {
	return RandBytesInCharset(n, Numerals)
}

// RandBytesAlphabet generate random alphabet bytes
func RandBytesAlphabet(n int) ([]byte, error) {
	return RandBytesInCharset(n, Alphabet)
}

// RandBytesAlphanumeric generate random alphabet and numberic bytes
func RandBytesAlphanumeric(n int) ([]byte, error) {
	return RandBytesInCharset(n, Alphanumeric)
}

// RandAscii generate random ascii string
func RandAscii(n int) (string, error) {
	return RandInCharset(n, Ascii)
}

// RandNumberal generate random numberal string
func RandNumberal(n int) (string, error) {
	return RandInCharset(n, Numerals)
}

// RandAlphabet generate random alphabet string
func RandAlphabet(n int) (string, error) {
	return RandInCharset(n, Alphabet)
}

// RandAlphanumeric generate random alphabet and numberic string
func RandAlphanumeric(n int) (string, error) {
	return RandInCharset(n, Alphanumeric)
}

// RandInt generate random integer
func RandInt(max int) int {
	i, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return int(i.Int64())
}
