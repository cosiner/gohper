package resource

import (
	"encoding/xml"
	"io"

	io2 "github.com/cosiner/gohper/lib/io"
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
	if key == "" {
		return xml.NewEncoder(w).Encode(value)
	}

	ew := io2.NewErrorWriter(w)
	ew.Write(XMLTagStart)
	ew.WriteString(key)
	ew.Write(XMLTagEnd)
	switch s := value.(type) {
	case string:
		ew.WriteString(s)
	case []byte:
		ew.Write(s)
	default:
		if ew.Error == nil {
			ew.Error = xml.NewEncoder(w).Encode(value)
		}
	}
	ew.Write(XMLTagCloseStart)
	ew.WriteString(key)
	ew.Write(XMLTagEnd)
	return ew.Error
}

func (XML) Receive(r io.Reader, value interface{}) error {
	return xml.NewDecoder(r).Decode(value)
}

var (
	XMLTagStart      = []byte("<")
	XMLTagEnd        = []byte(">")
	XMLTagCloseStart = []byte("</")
)
