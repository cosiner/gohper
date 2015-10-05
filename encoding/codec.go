package encoding

import "encoding/json"

type UnmarshalFunc func([]byte, interface{}) error

func (u UnmarshalFunc) Unmarshal(b []byte, v interface{}) error {
	return u(b, v)
}

type MarshalFunc func(interface{}) ([]byte, error)

func (m MarshalFunc) Marshal(v interface{}) ([]byte, error) {
	return m(v)
}

type PoolFunc func([]byte)

func (p PoolFunc) Pool(b []byte) {
	p(b)
}

type Codec interface {
	Unmarshal([]byte, interface{}) error
	Marshal(interface{}) ([]byte, error)
	Pool([]byte)
}

type funcCodec struct {
	UnmarshalFunc
	MarshalFunc
	PoolFunc
}

func NewFuncCodec(unmarshal UnmarshalFunc, marshal MarshalFunc, pool PoolFunc) Codec {
	if pool == nil {
		pool = func([]byte) {}
	}
	return funcCodec{
		UnmarshalFunc: unmarshal,
		MarshalFunc:   marshal,
		PoolFunc:      pool,
	}
}

var JSON = NewFuncCodec(json.Unmarshal, json.Marshal, nil)
