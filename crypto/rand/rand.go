package rand

import (
	"crypto/rand"
	"math/big"

	"github.com/cosiner/gohper/errors"
)

// Charset of characters to use for generating random strings
const (
	ALPHABET     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	NUMERALS     = "1234567890"
	ALPHANUMERIC = ALPHABET + NUMERALS
	ASCII        = ALPHANUMERIC + "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"

	ErrNegativeNum  = errors.Err("Number cannot be negative")
	ErrEmptyCharset = errors.Err("Charset cannot be empty")
)

type BytesFunc func(n int, charset string) ([]byte, error)
type StringFunc func(n int, charset string) (string, error)

// B generate a random bytes in gived charset
var B BytesFunc = func(n int, charset string) ([]byte, error) {
	if n <= 0 {
		return nil, ErrNegativeNum
	}
	var l int
	if l = len(charset); l == 0 {
		return nil, ErrEmptyCharset
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

// String is same as Bytes, return string
var S StringFunc = func(n int, charset string) (string, error) {
	b, err := B(n, charset)
	if err == nil {
		return string(b), nil
	}
	return "", err
}

// ASCIIBytes generate random ASCII bytes
func (f BytesFunc) ASCII(n int) ([]byte, error) {
	return f(n, ASCII)
}

// NumberalBytes generate random ASCII bytes
func (f BytesFunc) Numberal(n int) ([]byte, error) {
	return f(n, NUMERALS)
}

// AlphabetBytes generate random ALPHABET bytes
func (f BytesFunc) Alphabet(n int) ([]byte, error) {
	return f(n, ALPHABET)
}

// AlphanumericBytes generate random ALPHABET and numberic bytes
func (f BytesFunc) Alphanumeric(n int) ([]byte, error) {
	return f(n, ALPHANUMERIC)
}

// Ascii generate random ASCII string
func (f StringFunc) Ascii(n int) (string, error) {
	return f(n, ASCII)
}

// Numberal generate random numberal string
func (f StringFunc) Numberal(n int) (string, error) {
	return f(n, NUMERALS)
}

// Alphabet generate random ALPHABET string
func (f StringFunc) Alphabet(n int) (string, error) {
	return f(n, ALPHABET)
}

// Alphanumeric generate random ALPHABET and numberic string
func (f StringFunc) Alphanumeric(n int) (string, error) {
	return f(n, ALPHANUMERIC)
}
