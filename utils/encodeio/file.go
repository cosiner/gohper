// Package encodeio help read from/write to config file
package encodeio

import (
	"os"

	"github.com/cosiner/gohper/encoding"
	"github.com/cosiner/gohper/os2/file"
)

func Read(fname string, v interface{}, codec encoding.Codec) error {
	return file.Read(fname, func(fd *os.File) error {
		return codec.Decode(fd, v)
	})
}

func ReadJSON(fname string, v interface{}) error {
	return Read(fname, v, encoding.JSON{})
}

func Write(fname string, v interface{}, codec encoding.Codec) error {
	return file.Write(fname, func(fd *os.File) error {
		return codec.Encode(fd, v)
	})
}

func Trunc(fname string, v interface{}, codec encoding.Codec) error {
	return file.Trunc(fname, func(fd *os.File) error {
		return codec.Encode(fd, v)
	})
}
