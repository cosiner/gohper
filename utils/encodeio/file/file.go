// Package file help read from/write to config file
package file

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/cosiner/gohper/bytes2"
	"github.com/cosiner/gohper/os2/file"
	"github.com/cosiner/gohper/strings2"
	"github.com/cosiner/gohper/unsafe2"
	"github.com/cosiner/gohper/utils/encodeio"
	"github.com/cosiner/gohper/utils/pair"
)

type WriteMode bool

const TRUNC WriteMode = true
const APPEND WriteMode = false

func (m WriteMode) Write(fname string, v interface{}, encoder encodeio.EncodeFunc) error {
	return file.OpenOrCreate(fname, bool(m), func(fd *os.File) error {
		return encodeio.Write(fd, v, encoder)
	})
}

// WriteString write string to writer
func (m WriteMode) WriteString(fname, str string) (c int, err error) {
	err = file.OpenOrCreate(fname, bool(m), func(fd *os.File) error {
		c, err = fd.Write(unsafe2.Bytes(str))

		return err
	})

	return
}

// WriteGOB write interface{} to writer use gob encoder
func (m WriteMode) WriteGOB(fname string, v interface{}) error {
	return file.OpenOrCreate(fname, bool(m), func(fd *os.File) error {
		return encodeio.WriteGOB(fd, v)
	})
}

// WriteJSON write interface{} to writer use json encoder
func (m WriteMode) WriteJSON(fname string, v interface{}) error {
	return file.OpenOrCreate(fname, bool(m), func(fd *os.File) error {
		return encodeio.WriteJSON(fd, v)
	})
}

// WriteXML write interface{} to writer use xml encoder
func (m WriteMode) WriteXML(fname string, v interface{}) error {
	return file.OpenOrCreate(fname, bool(m), func(fd *os.File) error {
		return encodeio.WriteXML(fd, v)
	})
}

func (m WriteMode) WriteGZIP(fname string, v interface{}) (err error) {
	return file.OpenOrCreate(fname, bool(m), func(fd *os.File) error {
		return encodeio.WriteGZIP(fd, v)
	})
}

func Read(fname string, v interface{}, decoder encodeio.DecodeFunc) error {
	return file.Read(fname, func(fd *os.File) error {
		return encodeio.Read(fd, v, decoder)
	})
}

func ReadString(fname string) (s string, err error) {
	err = file.Read(fname, func(fd *os.File) error {
		s, err = encodeio.ReadString(fd)

		return err
	})

	return
}

func ReadGOB(fname string, v interface{}) error {
	return file.Read(fname, func(fd *os.File) error {
		return gob.NewDecoder(fd).Decode(v)
	})
}

func ReadJSON(fname string, v interface{}) error {
	return file.Read(fname, func(fd *os.File) error {
		return json.NewDecoder(fd).Decode(v)
	})
}

func ReadCommenttedJSON(fname, comment string, v interface{}) error {
	var trim = func(line, delim []byte) []byte {
		if len(delim) > 0 && bytes.HasPrefix(bytes.TrimSpace(line), delim) {
			return nil
		}

		return line
	}

	return file.Read(fname, func(fd *os.File) error {
		bs, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}
		bs = bytes2.MultipleLineOperate(bs, unsafe2.Bytes(comment), trim)

		return json.Unmarshal(bs, v)
	})
}

func ReadXML(fname string, v interface{}) error {
	return file.Read(fname, func(fd *os.File) error {
		return xml.NewDecoder(fd).Decode(v)
	})
}

func ReadGZIP(fname string) (data []byte, err error) {
	err = file.Read(fname, func(fd *os.File) error {
		data, err = encodeio.ReadGZIP(fd)

		return err
	})

	return
}

func ReadProperties(fname string) (map[string]string, error) {
	props := make(map[string]string)
	err := file.Filter(fname, func(_ int, line []byte) ([]byte, error) {
		p := pair.Parse(unsafe2.String(line), "=").Trim()
		if p.HasKey() {
			props[p.Key] = strings2.TrimAfter(p.Value, "#")
		}

		return line, nil
	})

	return props, err
}
