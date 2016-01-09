// Package encodeio help read from/write to config file
package encodeio

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/cosiner/gohper/encoding"
	"github.com/cosiner/gohper/os2/file"
)

var (
	commentPrefix = []byte("//")
)

func Read(fname string, v interface{}, codec encoding.Codec) error {
	return file.Read(fname, func(fd *os.File) error {
		return codec.Decode(fd, v)
	})
}

func ReadJSON(fname string, v interface{}) error {
	return Read(fname, v, encoding.JSON{})
}

func ReadJSONWithComment(fname string, v interface{}) error {
	return file.Read(fname, func(fd *os.File) error {
		buf := bytes.NewBuffer(make([]byte, 0, 1024))
		br := bufio.NewReader(fd)

		for {
			line, _, err := br.ReadLine()
			if err == io.EOF {
				if len(line) == 0 {
					break
				}
			} else if err != nil {
				return err
			}

			if !bytes.HasPrefix(bytes.TrimSpace(line), commentPrefix) {
				buf.Write(line)
			}
		}
		return encoding.JSON{}.Unmarshal(buf.Bytes(), v)
	})
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
