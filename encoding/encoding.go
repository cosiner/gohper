package encoding

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/hex"
	"io"
	"io/ioutil"
)

var (
	HEX       = Hex{}
	Base64Std = &Base64{
		Encoding: base64.StdEncoding,
	}
	Base64URL = &Base64{
		Encoding: base64.URLEncoding,
	}
	Gzip = Compress{
		BufSize: 64,
		NewReader: func(r io.Reader) (io.ReadCloser, error) {
			rd, err := gzip.NewReader(r)
			return rd, err
		},
		NewWriter: func(w io.Writer) io.WriteCloser {
			return gzip.NewWriter(w)
		},
	}
	Zlib = Compress{
		BufSize: 64,
		NewReader: func(r io.Reader) (io.ReadCloser, error) {
			rd, err := zlib.NewReader(r)
			return rd, err
		},
		NewWriter: func(w io.Writer) io.WriteCloser {
			return zlib.NewWriter(w)
		},
	}
)

type Encoding interface {
	Encode([]byte) []byte
	Decode([]byte) ([]byte, error)
}

type Hex struct{}

func (Hex) Encode(src []byte) []byte {
	l := len(src)
	dst := make([]byte, hex.EncodedLen(l))
	l = hex.Encode(dst, src)
	return dst[:l]
}

func (Hex) Decode(src []byte) ([]byte, error) {
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
	l := len(src)
	dst := make([]byte, b.Encoding.EncodedLen(l))
	b.Encoding.Encode(dst, src)
	return dst
}

func (b *Base64) Decode(src []byte) ([]byte, error) {
	l := len(src)
	dst := make([]byte, b.Encoding.DecodedLen(l))
	l, err := b.Encoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}

	return dst[:l], nil
}

type Compress struct {
	BufSize   int
	NewWriter func(io.Writer) io.WriteCloser
	NewReader func(io.Reader) (io.ReadCloser, error)
}

func (c Compress) Encode(src []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, c.BufSize))
	w := c.NewWriter(buf)
	w.Write(src)
	w.Close()

	return buf.Bytes()
}

func (c Compress) Decode(src []byte) ([]byte, error) {
	r, err := c.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}
	dst, err := ioutil.ReadAll(r)
	r.Close()

	return dst, err
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
	newP := make(Pipe, len(p)+len(encs))
	copy(newP, encs)
	copy(newP[len(encs):], p)
	return newP
}

func (p Pipe) Append(encs ...Encoding) Pipe {
	return append(p, encs...)
}
