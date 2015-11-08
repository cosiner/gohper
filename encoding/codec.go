package encoding

import (
	"encoding/json"
	"io"
)

type Codec interface {
	Encode(io.Writer, interface{}) error
	Marshal(interface{}) ([]byte, error)
	Decode(io.Reader, interface{}) error
	Unmarshal([]byte, interface{}) error
	Pool([]byte)
}

type JSON struct{}

func (JSON) Encode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (JSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (JSON) Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func (JSON) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (JSON) Pool([]byte) {}

var DefaultCodec Codec = JSON{}

func Encode(w io.Writer, v interface{}) error {
	return DefaultCodec.Encode(w, v)
}

func Marshal(v interface{}) ([]byte, error) {
	return DefaultCodec.Marshal(v)
}

func Decode(r io.Reader, v interface{}) error {
	return DefaultCodec.Decode(r, v)
}

func Unmarshal(data []byte, v interface{}) error {
	return DefaultCodec.Unmarshal(data, v)
}
func Pool(data []byte) {
	DefaultCodec.Pool(data)
}
