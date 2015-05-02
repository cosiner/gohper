package resource

import (
	"encoding/json"
	"io"

	io2 "github.com/cosiner/gohper/lib/io"
)

type JSON struct{}

func (JSON) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (JSON) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (JSON) Pool([]byte) {}

func (JSON) Send(w io.Writer, key string, value interface{}) error {
	if key == "" {
		return json.NewEncoder(w).Encode(value)
	}

	ew := io2.NewErrorWriter(w)
	ew.Write(JSONObjStart)
	ew.WriteString(key)
	switch s := value.(type) {
	case string:
		ew.Write(JSONQuoteMid)
		ew.WriteString(s)
		ew.Write(JSONQuoteEnd)
	case []byte:
		ew.Write(JSONQuoteMid)
		ew.Write(s)
		ew.Write(JSONQuoteEnd)
	default:
		ew.Write(JSONObjMid)
		if ew.Error == nil {
			ew.Error = json.NewEncoder(w).Encode(value)
		}
		ew.Write(JSONObjEnd)
	}
	return ew.Error
}

func (JSON) Receive(r io.Reader, value interface{}) error {
	d := json.NewDecoder(r)
	d.UseNumber()
	return d.Decode(value)
}

var (
	JSONObjStart = []byte(`{"`)
	JSONObjMid   = []byte(`":`)
	JSONQuoteMid = []byte(`":"`)
	JSONObjEnd   = []byte("}")
	JSONQuoteEnd = []byte(`"}`)
)
