package regexp

import (
	"regexp"

	"testing"
)

var pattern = MustCompile("(?m)^func\\ (?:\\(.*\\)\\ )?(?P<funcname>\\S+)\\(")
var code = `
type T int
func (*T) TT() {
}

func ConsoleError(format string, v ...interface{}) {
}

// func Parse(fname string) (funcs []string) {
// }
`

func TestMatch(t *testing.T) {
	t.Log(pattern.AllSubmatch(code))
	t.Log(pattern.AllSubmatchAtIndex(code, 1))
	t.Log(pattern.AllSubmatchMap(code))
	t.Log(pattern.AllSubmatchWithName(code, "funcname"))
	t.Log(pattern.SingleSubmatch(code))
	t.Log(pattern.SingleSubmatchAtIndex(code, 1))
	t.Log(pattern.SingleSubmatchMap(code))
	t.Log(pattern.SingleSubmatchWithName(code, "funcname"))
}

func TestNoMatch(t *testing.T) {
	r := MustCompile("1234.*")
	t.Log(r.SingleSubmatchMap("001234a"))
}

func BenchmarkMatchOnly(b *testing.B) {
	pattern := MustCompile("^/[a-zA-Z]*/[0-9]*$")
	for i := 0; i < b.N; i++ {
		_ = pattern.MatchString("/user/123")
	}
}

func BenchmarkMatchNoMap(b *testing.B) {
	pattern := MustCompile("^/[a-zA-Z]*/[0-9]*$")
	for i := 0; i < b.N; i++ {
		_, _ = pattern.SingleSubmatch("/user/123")
	}
}

func BenchmarkRgxMap(b *testing.B) {
	pattern := MustCompile("^/[a-zA-Z]*/[0-9]*$")
	for i := 0; i < b.N; i++ {
		_, _ = pattern.SingleSubmatchMap("/user/123")
	}
}

func BenchmarkRgxSlice(b *testing.B) {
	pattern := regexp.MustCompile("^/[a-zA-Z]*/[0-9]*/[0-9]*/[0-9]*$")
	for i := 0; i < b.N; i++ {
		_ = pattern.FindStringSubmatch("/user/123/123/123")
	}
}
