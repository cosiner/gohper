package resource

import (
	"io"

	io2 "github.com/cosiner/gohper/lib/io"
	"github.com/pquerna/ffjson/ffjson"
)

type FFJSON struct{}

func (FFJSON) Marshal(v interface{}) ([]byte, error) {
	return ffjson.Marshal(v)
}

func (FFJSON) Pool(data []byte) {
	ffjson.Pool(data)
}

func (FFJSON) Send(w io.Writer, key string, value interface{}) error {
	ew := io2.NewErrorWriter(w)
	var data []byte

	if key == "" {
		data, ew.Error = ffjson.Marshal(value)
		ew.WriteDo(data, ffjson.Pool)
		return ew.Error
	}

	ew.Write(JSONObjStart)
	ew.WriteString(key)
	if ew.Error == nil {
		if s, is := value.(string); is { // send string value
			ew.Write(JSONQuoteMid)
			ew.WriteString(s)
			ew.Write(JSONQuoteEnd)
		} else { // send other value
			data, ew.Error = ffjson.Marshal(value)
			ew.Write(JSONObjMid)
			ew.WriteDo(data, ffjson.Pool)
			ew.Write(JSONObjEnd)
		}
	}
	return ew.Error
}

func (FFJSON) Unmarshal(data []byte, v interface{}) error {
	return ffjson.Unmarshal(data, v)
}

func (FFJSON) Receive(r io.Reader, v interface{}) error {
	return ffjson.NewDecoder().DecodeReader(r, v)
}
