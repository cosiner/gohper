package url2

import (
	"net/url"
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestQueryEscape(t *testing.T) {
	tt := testing2.Wrap(t)
	qs := map[string]string{
		"A": "DD",
		"Z": "DD",
	}
	q, _ := QueryEscape(qs, nil)
	s := string(q)
	tt.NE(url.QueryEscape("A=DD&Z=DD") == s, url.QueryEscape("Z=DD&A=DD") == s)
}

func TestParamString(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.Eq(Param(1), "1")
	tt.Eq(Param("ddd"), "ddd")
	tt.Eq(Param([]int{1, 2, 3}), "1,2,3")
	tt.Eq(Param([]string{"1", "2", "3"}), "1,2,3")
	tt.Eq(Param([]byte("ddd")), "ddd")

	defer tt.Recover()

	Param(int64(2))
}
