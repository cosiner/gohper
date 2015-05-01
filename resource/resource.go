package resource

import (
	"io"
	"strings"
)

const (
	RES_JSON  = "json"
	RES_XML   = "xml"
	RES_HTML  = "html"
	RES_PLAIN = "plain"
)

type (
	Master struct {
		Resources map[string]Resource
		Default   string
		TypeOf    func(typ string) string // Defaultault use 'TypeOf'
	}

	Resource interface {
		Marshal(interface{}) ([]byte, error)
		Pool([]byte)
		Unmarshal([]byte, interface{}) error
		Send(w io.Writer, key string, value interface{}) error
		Receive(r io.Reader, v interface{}) error
	}
)

func TypeOf(typ string) string {
	if typ != "" {
		switch {
		case strings.Contains(typ, RES_JSON):
			return RES_JSON
		case strings.Contains(typ, RES_XML):
			return RES_XML
		case strings.Contains(typ, RES_HTML):
			return RES_HTML
		case strings.Contains(typ, RES_PLAIN):
			return RES_PLAIN
		}
	}
	return ""
}

func NewMaster() Master {
	return Master{
		Resources: make(map[string]Resource),
		Default:   RES_JSON,
		TypeOf:    TypeOf,
	}
}

func (rm *Master) Use(typ string, res Resource) {
	rm.Resources[typ] = res
}

func (rm *Master) DefUse(typ string, res Resource) {
	if res != nil {
		rm.Use(typ, res)
	}
	rm.Default = typ
}

func (rm *Master) Resource(typ string) Resource {
	typ = rm.TypeOf(typ)
	if typ == "" {
		typ = rm.Default
	}
	return rm.Resources[typ]
}
