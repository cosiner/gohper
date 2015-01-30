package urlmatcher

import (
	"net/url"

	. "github.com/cosiner/golib/errors"
	"github.com/cosiner/golib/test"

	"testing"
)

func TestMatch(t *testing.T) {
	tt := test.WrapTest(t)
	url := new(url.URL)
	// Match given regexp
	matcher, err := NewMatcher("/{name:.*}aa/{age:.*}")
	OnErrPanic(err)
	url.Path = "/testaa/123"
	matchMap, match := matcher.Match(url)
	tt.AssertEq("Match00", "/(?P<name>[a-zA-Z0-9]*)aa/(?P<age>[a-zA-Z0-9]*)", matcher.Pattern())
	tt.AssertTrue("Match0", match)
	tt.AssertEq("MATCH1", "test", matchMap["name"])
	tt.AssertEq("MATCH2", "123", matchMap["age"])

	// Match other regexp
	matcher, err = NewMatcher("/{name}Abc/{age}")
	OnErrPanic(err)
	tt.AssertEq("Match77", "/(?P<name>[a-zA-Z0-9]*)Abc/(?P<age>[a-zA-Z0-9]*)", matcher.Pattern())
	url.Path = "/LosuAbc/123"
	matchMap, match = matcher.Match(url)
	tt.AssertTrue("Match7", match)
	tt.AssertEq("Match8", "Losu", matchMap["name"])
	tt.AssertEq("Match9", "123", matchMap["age"])

	// Match Literal
	matcher, err = NewMatcher("/user/123")
	OnErrPanic(err)
	tt.AssertEq("Match33", "/user/123", matcher.Pattern())
	url.Path = "/user/123"
	_, match = matcher.Match(url)
	tt.AssertTrue("Match32", matcher.IsLiteral())
	tt.AssertTrue("Match3", match)
	url.Path = "/user"
	_, match = matcher.Match(url)
	tt.AssertFalse("Match4", match)
	url.Path = "/user/12"
	_, match = matcher.Match(url)
	tt.AssertFalse("Match5", match)

	// Compile
	typ, pat := Compile("{proto}://", '{', '}', true)
	tt.AssertEq("Compile1", REGEXP, typ)
	tt.AssertEq("Compile2", "(?P<proto>[a-zA-Z0-9]*)://", pat)
	typ, pat = Compile("{proto}://{id:\\d\\{1,2\\}}", '{', '}', true)
	tt.AssertEq("Compile1", REGEXP, typ)
	tt.AssertEq("Compile2", "(?P<proto>[a-zA-Z0-9]*)://(?P<id>\\d{1,2})", pat)
}
