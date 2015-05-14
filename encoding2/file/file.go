// Package file help read from/write to config file
package file

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"os"

	"github.com/cosiner/gohper/encoding2"
	"github.com/cosiner/gohper/os2/file"
	"github.com/cosiner/gohper/unsafe2"
)

type WriteMode bool

const TRUNC WriteMode = true
const APPEND WriteMode = false

func (m WriteMode) Write(fname string, v interface{}, encoder encoding2.EncodeFunc) error {
	return file.OpenOrCreate(fname, bool(m), func(fd *os.File) error {
		return encoding2.Write(fd, v, encoder)
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
		return encoding2.WriteGOB(fd, v)
	})
}

// WriteJSON write interface{} to writer use json encoder
func (m WriteMode) WriteJSON(fname string, v interface{}) error {
	return file.OpenOrCreate(fname, bool(m), func(fd *os.File) error {
		return encoding2.WriteJSON(fd, v)
	})
}

// WriteXML write interface{} to writer use xml encoder
func (m WriteMode) WriteXML(fname string, v interface{}) error {
	return file.OpenOrCreate(fname, bool(m), func(fd *os.File) error {
		return encoding2.WriteXML(fd, v)
	})
}

func (m WriteMode) WriteGZIP(fname string, v interface{}) (err error) {
	return file.OpenOrCreate(fname, bool(m), func(fd *os.File) error {
		return encoding2.WriteGZIP(fd, v)
	})
}

func Read(fname string, v interface{}, decoder encoding2.DecodeFunc) error {
	return file.Read(fname, func(fd *os.File) error {
		return encoding2.Read(fd, v, decoder)
	})
}

func ReadString(fname string) (s string, err error) {
	err = file.Read(fname, func(fd *os.File) error {
		s, err = encoding2.ReadString(fd)

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

func ReadXML(fname string, v interface{}) error {
	return file.Read(fname, func(fd *os.File) error {
		return xml.NewDecoder(fd).Decode(v)
	})
}

func ReadGZIP(fname string) (data []byte, err error) {
	err = file.Read(fname, func(fd *os.File) error {
		data, err = encoding2.ReadGZIP(fd)

		return err
	})

	return
}
