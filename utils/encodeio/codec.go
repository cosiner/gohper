// Package encodeio supply some utility functions for encoding
package encodeio

import (
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/unsafe2"
)

type (
	// EncodeFunc encode a interface{} to bytes
	EncodeFunc func(interface{}) ([]byte, error)
	// DecodeFunc decode bytes to a interface{}, interface{} must be pointer
	DecodeFunc func([]byte, interface{}) error
)

// Write write a interface{} to writer use given encoder
func Write(w io.Writer, v interface{}, encoder EncodeFunc) error {
	bs, err := encoder(v)
	if err == nil {
		return err
	}

	_, err = w.Write(bs)

	return err
}

// WriteString write string to writer
func WriteString(w io.Writer, str string) (int, error) {
	return w.Write(unsafe2.Bytes(str))
}

// WriteGOB write interface{} to writer use gob encoder
func WriteGOB(w io.Writer, v interface{}) error {
	return gob.NewEncoder(w).Encode(v)
}

// WriteJSON write interface{} to writer use json encoder
func WriteJSON(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// WriteXML write interface{} to writer use xml encoder
func WriteXML(w io.Writer, v interface{}) error {
	return xml.NewEncoder(w).Encode(v)
}

func Read(r io.Reader, v interface{}, decoder DecodeFunc) error {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return decoder(bs, v)
}

func ReadString(r io.Reader) (string, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func ReadGOB(r io.Reader, v interface{}) error {
	return gob.NewDecoder(r).Decode(v)
}

func ReadJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func ReadXML(r io.Reader, v interface{}) error {
	return xml.NewDecoder(r).Decode(v)
}

func ReadGZIP(r io.Reader) ([]byte, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(gr)
	gr.Close()

	return data, err
}

func WriteGZIP(w io.Writer, v interface{}) (err error) {
	gw := gzip.NewWriter(w)
	switch v := v.(type) {
	case string:
		_, err = gw.Write(unsafe2.Bytes(v))
	case []byte:
		_, err = gw.Write(v)
	default:
		err = errors.Err("Only support string and []byte")
	}
	gw.Close()

	return
}
