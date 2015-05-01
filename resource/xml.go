package resource

import (
	"encoding/xml"
	"io"

	eio "github.com/cosiner/gohper/lib/io"
)

type XML struct{}

func (XML) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (XML) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func (XML) Pool([]byte) {}

func (XML) Send(w io.Writer, key string, value interface{}) error {
	var err error
	if key != "" {
		eio.ErrPtrWrite(&err, w, XMLTagStart)
		eio.ErrPtrWriteString(&err, w, key)
		eio.ErrPtrWrite(&err, w, XMLTagEnd)
		switch s := value.(type) {
		case string:
			eio.ErrPtrWriteString(&err, w, s)
		case []byte:
			eio.ErrPtrWrite(&err, w, s)
		default:
			if err == nil {
				err = xml.NewEncoder(w).Encode(value)
			}
		}
		eio.ErrPtrWrite(&err, w, XMLTagCloseStart)
		eio.ErrPtrWriteString(&err, w, key)
		eio.ErrPtrWrite(&err, w, XMLTagEnd)
	} else {
		err = xml.NewEncoder(w).Encode(value)
	}
	return err
}

func (XML) Receive(r io.Reader, value interface{}) error {
	return xml.NewDecoder(r).Decode(value)
}

var (
	XMLTagStart      = []byte("<")
	XMLTagEnd        = []byte(">")
	XMLTagCloseStart = []byte("</")
)
