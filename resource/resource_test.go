package resource

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func resType(typ string) string {
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
	return RES_JSON
}

func TestResourceMaster(t *testing.T) {
	tt := test.Wrap(t)
	rm := NewMaster()
	rm.DefUse(RES_JSON, JSON{})
	rm.Use(RES_XML, XML{})

	res := rm.Resources[resType("application/json;charset=utf-8")]
	_, is := res.(JSON)
	tt.True(is)

	res = rm.Resources[resType("application/xml;charset=utf-8")]
	_, is = res.(XML)
	tt.True(is)

	res = rm.Resources[resType("abcdefghijklmn")]
	_, is = res.(JSON)
	tt.True(is)
}

func TestURLEncode(t *testing.T) {
	tt := test.Wrap(t)
	enc := NewURLEncode(nil)
	// enc.Send(os.Stdout, "", map[string]string{
	// 	"AA":  "DD",
	// 	"111": "123",
	// })
	// os.Stdout.WriteString("\n")
	// enc.Send(os.Stdout, "AAA", "ddd")
	// os.Stdout.WriteString("\n")
	// enc.Send(os.Stdout, "AAA", "ddd&zzz=AAA")
	// os.Stdout.WriteString("\n")
	type User struct {
		Name string `encode:"nm"`
		Age  int    `encode:"Age"`
	}
	buf := bytes.NewBuffer(make([]byte, 0))
	enc.Send(buf, "", User{"123", 123})
	user := &User{}
	enc.Receive(buf, user)
	tt.Eq("123", user.Name)
	tt.Eq(123, user.Age)
}
