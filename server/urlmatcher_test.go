package server

import (
	. "github.com/cosiner/golib/errors"
	"github.com/cosiner/golib/test"

	"testing"
)

func TestMatch(t *testing.T) {
	tt := test.WrapTest(t)
	// Match given regexp
	matcher, err := NewMatcher("/{name:.*}aa/{age:.*}")
	OnErrPanic(err)
	values, match := matcher.Match("/testaa/123")
	tt.AssertEq("Match00", "^/(?P<name>[a-zA-Z0-9_]*)aa/(?P<age>[a-zA-Z0-9_]*)$", matcher.Pattern())
	tt.AssertTrue("Match0", match)
	tt.AssertEq("MATCH1", "test", matcher.ValueOf(values, "name"))
	tt.AssertEq("MATCH2", "123", matcher.ValueOf(values, "age"))

	// Match other regexp
	matcher, err = NewMatcher("/{name}Abc/{age}")
	OnErrPanic(err)
	tt.AssertEq("Match77", "^/(?P<name>[a-zA-Z0-9_]*)Abc/(?P<age>[a-zA-Z0-9_]*)$", matcher.Pattern())
	values, match = matcher.Match("/LosuAbc/123")
	tt.AssertTrue("Match7", match)
	tt.AssertEq("Match8", "Losu", matcher.ValueOf(values, "name"))
	tt.AssertEq("Match9", "123", matcher.ValueOf(values, "age"))

	// Match Literal
	matcher, err = NewMatcher("/user/123")
	OnErrPanic(err)
	tt.AssertEq("Match33", "/user/123", matcher.Pattern())
	_, match = matcher.Match("/user/123")
	tt.AssertTrue("Match32", matcher.IsLiteral())
	tt.AssertTrue("Match3", match)
	_, match = matcher.Match("/user")
	tt.AssertFalse("Match4", match)
	_, match = matcher.Match("/user/12")
	tt.AssertFalse("Match5", match)

	// Compile
	cp := &UrlCompiler{RegLeft: '|', RegRight: '|', NeedEscape: true}
	typ, pat := cp.Compile("|proto|://www\\.|site:.*|\\.com/")
	tt.AssertEq("Compile1", REGEXP, typ)
	tt.AssertEq("Compile2", "^(?P<proto>[a-zA-Z0-9_]*)://www\\.(?P<site>[a-zA-Z0-9_]*)\\.com/$", pat)

	cp = &UrlCompiler{RegLeft: '+', RegRight: '-', NeedEscape: true, NoReplace: true}
	typ, pat = cp.Compile("+proto-://+id:\\d{1,2}-")
	tt.AssertEq("Compile3", REGEXP, typ)
	tt.AssertEq("Compile4", "^(?P<proto>.*)://(?P<id>\\d{1,2})$", pat)
}

func BenchmarkLiteralMatcher(b *testing.B) {
	matcher, _ := NewMatcher("/user/123")
	url := "/user"
	for i := 0; i < b.N; i++ {
		matcher.MatchOnly(url)
	}
}

func BenchmarkRegexpMatcher(b *testing.B) {
	compiler := &UrlCompiler{
		Matchany: "[0-9]",
	}
	matcher, _ := NewMatcherWith("/{user}/{id}", compiler)
	url := "/user/123"
	for i := 0; i < b.N; i++ {
		_, _ = matcher.Match(url)
	}
}

func BenchmarkRegexpPrefixMatch(b *testing.B) {
	compiler := &UrlCompiler{
		Matchany: "[a-z]",
	}
	matcher, _ := NewMatcherWith("/123", compiler)
	url := "/123/123"
	for i := 0; i < b.N; i++ {
		_ = matcher.PrefixMatchOnly(url)
	}
}
