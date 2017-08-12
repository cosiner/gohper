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
}

var JSON Codec = Json{}

type Json struct{}

func (Json) Encode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (Json) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (Json) Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func (Json) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
