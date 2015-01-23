package regexp

import "testing"

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
