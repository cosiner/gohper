package encoding

import (
	"encoding/base64"
	"encoding/hex"
)

type Encoding interface {
	Encode([]byte) []byte
	Decode([]byte) ([]byte, error)
}

type Hex struct{}

func (Hex) Encode(src []byte) []byte {
	return HexEncode(src)
}

func (Hex) Decode(src []byte) ([]byte, error) {
	return HexDecode(src)
}

func HexEncode(src []byte) []byte {
	l := len(src)
	dst := make([]byte, hex.EncodedLen(l))
	l = hex.Encode(dst, src)
	return dst[:l]
}

func HexDecode(src []byte) ([]byte, error) {
	l := len(src)
	dst := make([]byte, hex.DecodedLen(l))
	l, err := hex.Decode(dst, src)
	if err != nil {
		return nil, err
	}

	return dst[:l], nil
}

type Base64 struct {
	Encoding *base64.Encoding
}

func (b *Base64) Encode(src []byte) []byte {
	return Base64Encode(b.Encoding, src)
}

func (b *Base64) Decode(src []byte) ([]byte, error) {
	return Base64Decode(b.Encoding, src)
}

func Base64Encode(enc *base64.Encoding, src []byte) []byte {
	l := len(src)
	dst := make([]byte, enc.EncodedLen(l))
	enc.Encode(dst, src)
	return dst
}

func Base64Decode(enc *base64.Encoding, src []byte) ([]byte, error) {
	l := len(src)
	dst := make([]byte, enc.DecodedLen(l))
	l, err := enc.Decode(dst, src)
	if err != nil {
		return nil, err
	}

	return dst[:l], nil
}

type Pipe []Encoding

func (p Pipe) Encode(src []byte) []byte {
	for i := 0; i < len(p); i++ {
		src = p[i].Encode(src)
	}
	return src
}

func (p Pipe) Decode(src []byte) ([]byte, error) {
	var err error
	var dst = src
	for i := len(p) - 1; i >= 0 && err == nil; i-- {
		dst, err = p[i].Decode(dst)
	}
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (p Pipe) Prepend(encs ...Encoding) Pipe {
	newP := make(Pipe, len(p) + len(encs))
	copy(newP, encs)
	copy(newP[len(encs):], p)
	return newP
}

func (p Pipe) Append(encs ...Encoding) Pipe {
	return append(p, encs...)
}
