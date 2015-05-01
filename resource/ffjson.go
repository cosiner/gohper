package resource

import (
	"io"
	"io/ioutil"

	eio "github.com/cosiner/gohper/lib/io"
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
	var (
		data []byte
		err  error
	)
	if key == "" { // send value only
		// alwayse use json marshal, simple string use Response.WriteString
		data, err = ffjson.Marshal(value)
		if err == nil {
			eio.ErrPtrWrite(&err, w, data)
			ffjson.Pool(data)
		}
		return err
	}

	eio.ErrPtrWrite(&err, w, JSONObjStart) // send key
	eio.ErrPtrWrite(&err, w, eio.Bytes(key))
	if err == nil {
		if s, is := value.(string); is { // send string value
			eio.ErrPtrWrite(&err, w, JSONQuoteMid)
			eio.ErrPtrWrite(&err, w, eio.Bytes(s))
			eio.ErrPtrWrite(&err, w, JSONQuoteEnd)
		} else { // send other value
			if data, err = ffjson.Marshal(value); err == nil {
				eio.ErrPtrWrite(&err, w, JSONObjMid)
				eio.ErrPtrWrite(&err, w, data)
				eio.ErrPtrWrite(&err, w, JSONObjEnd)
				ffjson.Pool(data)
			}
		}
	}
	return err
}

func (FFJSON) Unmarshal(data []byte, v interface{}) error {
	return ffjson.Unmarshal(data, v)
}

func (FFJSON) Receive(r io.Reader, v interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err == nil {
		err = ffjson.Unmarshal(data, v)
	}
	return err
}
