package regexp

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

// ?m specified match with multi line mode, ^ and & match line begin and end, rather than expression's begin and end
var pattern = MustCompile("(?m)^func (?:\\(.*\\) )?(?P<funcname>\\S+)\\(") // extract function name
var code = `
type T int
func (*T) TT()
}

func ConsoleError(format string, v ...interface{}) {
}

// ignored
// func Parse(fname string) (funcs []string) {
// }

func Parse2() {

}

`

func TestMatch(t *testing.T) {
	tt := testing2.Wrap(t)

	funcs := pattern.AllByName(code, "funcname")
	tt.DeepEq([]string{"TT", "ConsoleError", "Parse2"}, funcs)

	funcs = pattern.AllByIndex(code, 1)
	tt.DeepEq([]string{"TT", "ConsoleError", "Parse2"}, funcs)

	fs := pattern.First(code)
	tt.DeepEq([]string{"func (*T) TT(", "TT"}, fs)

	f := pattern.ByIndex(code, 1)
	tt.Eq("TT", f)
}
