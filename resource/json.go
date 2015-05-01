package resource

import (
	"encoding/json"
	"io"

	eio "github.com/cosiner/gohper/lib/io"
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
	var err error
	if key != "" {
		eio.ErrPtrWrite(&err, w, JSONObjStart)
		eio.ErrPtrWriteString(&err, w, key)
		switch s := value.(type) {
		case string:
			eio.ErrPtrWrite(&err, w, JSONQuoteMid)
			eio.ErrPtrWriteString(&err, w, s)
			eio.ErrPtrWrite(&err, w, JSONQuoteEnd)
		case []byte:
			eio.ErrPtrWrite(&err, w, JSONQuoteMid)
			eio.ErrPtrWrite(&err, w, s)
			eio.ErrPtrWrite(&err, w, JSONQuoteEnd)
		default:
			eio.ErrPtrWrite(&err, w, JSONObjMid)
			if err == nil {
				err = json.NewEncoder(w).Encode(value)
			}
			eio.ErrPtrWrite(&err, w, JSONObjEnd)
		}
	} else {
		err = json.NewEncoder(w).Encode(value)
	}
	return err
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
