// Package encoding supply some utility functions for encoding
package encoding

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"

	. "github.com/cosiner/gohper/lib/errors"

	"github.com/cosiner/gohper/lib/types"
)

type (
	// PowerReadWriter is a power reader and writer which can read and write
	// bytes/string/gob/json/xml data all in one
	PowerReadWriter interface {
		PowerReader
		PowerWriter
	}

	// CleanPowerReadWriter is a clean reader and writer without Read/Write method
	CleanPowerReadWriter interface {
		CleanPowerReader
		CleanPowerWriter
	}

	// PowerReader is a power reader which can read bytes/string/gob/jsob/xml data
	PowerReader interface {
		io.Reader
		CleanPowerReader
	}

	// CleanPowerReader is a clean reader without Read method
	CleanPowerReader interface {
		ReadString() (string, error)
		ReadJSON(interface{}) error
		ReadXML(interface{}) error
		ReadGOB(interface{}) error
	}

	// PowerReader is a power writer which can write bytes/string/gob/jsob/xml data
	PowerWriter interface {
		io.Writer
		CleanPowerWriter
	}

	// CleanPowerWriter is a clean  writer without Write method
	CleanPowerWriter interface {
		WriteString(string) (int, error)
		WriteJSON(interface{}) error
		WriteXML(interface{}) error
		WriteGOB(interface{}) error
	}

	// EncodeFunc encode a interface{} to bytes
	EncodeFunc func(interface{}) ([]byte, error)
	// DecodeFunc decode bytes to a interface{}, interface{} must be pointer
	DecodeFunc func([]byte, interface{}) error

	// powerReadWriter is a power reader and writer for read source and write
	// destination is different, also you can still use for same read source
	// and write destination, but not recommended, replace with powerRW
	powerReadWriter struct {
		powerReader
		powerWriter
	}

	// powerRW is a power reader and writer designed for same read source and
	// write destination
	powerRW struct {
		io.ReadWriter
	}

	// powerReader is a power reader
	powerReader struct {
		io.Reader
	}

	// powerWriter is a power writer
	powerWriter struct {
		io.Writer
	}
)

// NewPowerReader create a power reader from exist reader
func NewPowerReader(r io.Reader) PowerReader {
	return powerReader{r}
}

// NewPowerWriter create a new power writer from exist writer
func NewPowerWriter(w io.Writer) PowerWriter {
	return powerWriter{w}
}

// NewPowerReaderWriter create a power reader writer,
// it performed as read from reader, write to writer
func NewPowerReadWriter(r io.Reader, w io.Writer) PowerReadWriter {
	return powerReadWriter{powerReader{r}, powerWriter{w}}
}

// NewPowerReadWriterInOne create a power reader writer from ReadWriter,
func NewPowerReadWriterInOne(rw io.ReadWriter) PowerReadWriter {
	return powerRW{rw}
}

func (pw powerWriter) WriteString(s string) (int, error) {
	return WriteString(pw, s)
}

func (pw powerWriter) WriteJSON(v interface{}) error {
	return WriteJSON(pw, v)
}

func (pw powerWriter) WriteXML(v interface{}) error {
	return WriteXML(pw, v)
}

func (pw powerWriter) WriteGOB(v interface{}) error {
	return WriteGOB(pw, v)
}

func (pr powerReader) ReadString() (string, error) {
	return ReadString(pr)
}

func (pr powerReader) ReadJSON(v interface{}) error {
	return ReadJSON(pr, v)
}

func (pr powerReader) ReadXML(v interface{}) error {
	return ReadXML(pr, v)
}

func (pr powerReader) ReadGOB(v interface{}) error {
	return ReadGOB(pr, v)
}

func (prw powerRW) WriteString(s string) (int, error) {
	return WriteString(prw, s)
}

func (prw powerRW) WriteJSON(v interface{}) error {
	return WriteJSON(prw, v)
}

func (prw powerRW) WriteXML(v interface{}) error {
	return WriteXML(prw, v)
}

func (prw powerRW) WriteGOB(v interface{}) error {
	return WriteGOB(prw, v)
}

func (prw powerRW) ReadString() (string, error) {
	return ReadString(prw)
}

func (prw powerRW) ReadJSON(v interface{}) error {
	return ReadJSON(prw, v)
}

func (prw powerRW) ReadXML(v interface{}) error {
	return ReadXML(prw, v)
}

func (prw powerRW) ReadGOB(v interface{}) error {
	return ReadGOB(prw, v)
}

// GOBEncode encode parameter value to bytes use gob encoder
func GOBEncode(v interface{}) (res []byte, err error) {
	buffer := bytes.NewBuffer(make([]byte, 100))
	encoder := gob.NewEncoder(buffer)
	if err = encoder.Encode(v); err == nil {
		res = buffer.Bytes()
	}
	return
}

// GOBDecode decode bytes to given parameter use gob decoder
func GOBDecode(bs []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(bs)).Decode(v)
}

// Write write a interface{} to writer use given encoder
func Write(wr io.Writer, v interface{}, encoder EncodeFunc) error {
	bs, err := encoder(v)
	if err == nil {
		_, err = wr.Write(bs)
	}
	return err
}

// WriteString write string to writer
func WriteString(wr io.Writer, str string) (int, error) {
	return wr.Write([]byte(str))
}

// UnsafeWriteString write string to writer use unsafed convert from string to bytes
func UnsafeWriteString(wr io.Writer, str string) (int, error) {
	return wr.Write(types.UnsafeBytes(str))
}

// WriteGOB write interface{} to writer use gob encoder
func WriteGOB(wr io.Writer, v interface{}) error {
	return Write(wr, v, GOBEncode)
}

// WriteJSON write interface{} to writer use json encoder
func WriteJSON(wr io.Writer, v interface{}) error {
	return Write(wr, v, json.Marshal)
}

// WriteXML write interface{} to writer use xml encoder
func WriteXML(wr io.Writer, v interface{}) error {
	return Write(wr, v, xml.Marshal)
}

// Read read all bytes from reader then decode to given interface address
// interface{} must be pointer
func Read(rd io.Reader, v interface{}, decoder DecodeFunc) error {
	bs, err := ioutil.ReadAll(rd)
	if err == nil {
		err = decoder(bs, v)
	}
	return err
}

// ReadString read string from reader
func ReadString(rd io.Reader) (s string, err error) {
	bs, err := ioutil.ReadAll(rd)
	if err == nil {
		s = string(bs)
	}
	return
}

// UnsafeReadString read bytes from reader and unsafed convert to string
func UnsafeReadString(rd io.Reader) (s string, err error) {
	bs, err := ioutil.ReadAll(rd)
	if err == nil {
		s = types.UnsafeString(bs)
	}
	return
}

// ReadGOB read bytes from reader and decode to interface address use gob decoder
func ReadGOB(rd io.Reader, v interface{}) error {
	return Read(rd, v, GOBDecode)
}

// ReadJSON read bytes from reader and decode to interface address use json decoder
func ReadJSON(rd io.Reader, v interface{}) error {
	return Read(rd, v, json.Unmarshal)
}

// ReadXML read bytes from reader and decode to interface address use xml decoder
func ReadXML(rd io.Reader, v interface{}) error {
	return Read(rd, v, xml.Unmarshal)
}

func ReadGZIP(rd io.Reader, v interface{}) error {
	r, err := gzip.NewReader(rd)
	if err == nil {
		if s, is := v.(*string); !is {
			err = Err("value must be pointer to string")
		} else {
			var bs []byte
			if bs, err = ioutil.ReadAll(r); err == nil {
				*s = string(bs)
			}
		}
	}
	return err
}

func WriteGZIP(wr io.Writer, v interface{}) error {
	w := gzip.NewWriter(wr)
	switch v := v.(type) {
	case string:
		w.Write(types.UnsafeBytes(v))
	case []byte:
		w.Write(v)
	default:
		return Err("Only support string and []byte")
	}
	return nil
}
