// Package encoding supply some utility functions for encoding
package encoding

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"

	"github.com/cosiner/golib/types"
)

type (
	// EncodeFunc encode a interface{} to bytes
	EncodeFunc func(interface{}) ([]byte, error)
	// DecodeFunc decode bytes to a interface{}, interface{} must be pointer
	DecodeFunc func([]byte, interface{}) error
)

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
