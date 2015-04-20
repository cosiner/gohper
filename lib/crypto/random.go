package crypto

import (
	"crypto/rand"
	"math/big"

	"github.com/cosiner/gohper/lib/errors"
)

// Charset of characters to use for generating random strings
const (
	ALPHABET     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	NUMERALS     = "1234567890"
	ALPHANUMERIC = ALPHABET + NUMERALS
	ASCII        = ALPHANUMERIC + "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"

	numNegativeError  = errors.Err("Number cannot be negative")
	charsetEmptyError = errors.Err("Charset cannot be empty")
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

// RandBytesASCII generate random ASCII bytes
func RandBytesASCII(n int) ([]byte, error) {
	return RandBytesInCharset(n, ASCII)
}

// RandBytesNumberal generate random ASCII bytes
func RandBytesNumberal(n int) ([]byte, error) {
	return RandBytesInCharset(n, NUMERALS)
}

// RandBytesAlphabet generate random ALPHABET bytes
func RandBytesAlphabet(n int) ([]byte, error) {
	return RandBytesInCharset(n, ALPHABET)
}

// RandBytesAlphanumeric generate random ALPHABET and numberic bytes
func RandBytesAlphanumeric(n int) ([]byte, error) {
	return RandBytesInCharset(n, ALPHANUMERIC)
}

// RandASCII generate random ASCII string
func RandASCII(n int) (string, error) {
	return RandInCharset(n, ASCII)
}

// RandNumberal generate random numberal string
func RandNumberal(n int) (string, error) {
	return RandInCharset(n, NUMERALS)
}

// RandAlphabet generate random ALPHABET string
func RandAlphabet(n int) (string, error) {
	return RandInCharset(n, ALPHABET)
}

// RandAlphanumeric generate random ALPHABET and numberic string
func RandAlphanumeric(n int) (string, error) {
	return RandInCharset(n, ALPHANUMERIC)
}

// RandInt generate random integer
func RandInt(max int) int {
	i, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return int(i.Int64())
}
