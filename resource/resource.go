package resource

import (
	"io"
	"strings"
)

const (
	RES_JSON      = "json"
	RES_XML       = "xml"
	RES_HTML      = "html"
	RES_PLAIN     = "plain"
	RES_URLENCODE = "urlencoded"

	// Content Type
	CONTENTTYPE_PLAIN     = "text/plain;charset=utf-8"
	CONTENTTYPE_HTML      = "text/html;charset=utf-8"
	CONTENTTYPE_XML       = "text/xml;charset=utf-8"
	CONTENTTYPE_JSON      = "application/json;charset=utf-8"
	CONTENTTYPE_URLENCODE = "application/x-www-form-urlencoded;charset=utf-8"
)

type (
	Master struct {
		Resources map[string]Resource
		Default   string
		TypeOf    func(typ string) string // default use 'ResourceType'
	}

	Resource interface {
		Marshal(interface{}) ([]byte, error)
		Pool([]byte)
		Unmarshal([]byte, interface{}) error
		Send(w io.Writer, key string, value interface{}) error
		Receive(r io.Reader, v interface{}) error
	}
)

func ResourceType(typ string) string {
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
		case strings.Contains(typ, RES_URLENCODE):
			return RES_URLENCODE
		}
	}
	return ""
}

func ContentType(typ string) string {
	switch typ {
	case RES_JSON:
		return CONTENTTYPE_JSON
	case RES_XML:
		return CONTENTTYPE_XML
	case RES_HTML:
		return CONTENTTYPE_HTML
	case RES_PLAIN:
		return CONTENTTYPE_PLAIN
	case RES_URLENCODE:
		return CONTENTTYPE_URLENCODE
	}
	return ""
}

func NewMaster() Master {
	return Master{
		Resources: make(map[string]Resource),
		Default:   RES_JSON,
		TypeOf:    ResourceType,
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
