package encoding

import (
	"errors"

	"github.com/cosiner/gohper/unsafe2"
)

type RawBytes []byte

func (r RawBytes) Marshal() ([]byte, error) {
	return []byte(r), nil
}

func (r *RawBytes) Unmarshal(data []byte) error {
	if r == nil {
		return errors.New("encoding.RawBytes: Unmarshal on nil pointer")
	}
	*r = append((*r)[0:0], data...)
	return nil
}

func (r RawBytes) MarshalJSON() ([]byte, error) {
	return r.Marshal()
}

func (r *RawBytes) UnmarshalJSON(data []byte) error {
	return r.Unmarshal(data)
}

type RawString string

func (r RawString) Marshal() ([]byte, error) {
	return unsafe2.Bytes(string(r)), nil
}

func (r RawString) Unmarshal(data []byte) error {
	panic("encoding.RawString: Unmarshal on immutable string")
}

func (r RawString) MarshalJSON() ([]byte, error) {
	return r.Marshal()
}

func (r RawString) UnmarshalJSON(data []byte) error {
	return r.Unmarshal(data)
}
