// Package encoding supply some utility functions for encoding
package encoding

import (
	"bytes"
	"encoding/gob"
)

// GobEncode encode parameter value to bytes use gob encoder
func GobEncode(v interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 100))
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(v); err == nil {
		return buffer.Bytes(), nil
	} else {
		return nil, err
	}
}

// GobDecode decode bytes to given parameter use gob decoder
func GobDecode(bs []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(bs)).Decode(v)
}
