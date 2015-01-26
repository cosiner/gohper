package urlmatcher

import (
	. "github.com/cosiner/golib/errors"
	"github.com/cosiner/golib/test"

	"testing"
)

func TestMatch(t *testing.T) {
	tt := test.WrapTest(t)
	// Match given regexp
	matcher, err := Compile("/{name:.*}aa/{age:.*}")
	OnErrPanic(err)
	matchMap, match := matcher.Match("/testaa/123")
	tt.AssertEq("Match00", "/(?P<name>.*)aa/(?P<age>.*)", matcher.Pattern())
	tt.AssertTrue("Match0", match)
	tt.AssertEq("MATCH1", "test", matchMap["name"])
	tt.AssertEq("MATCH2", "123", matchMap["age"])

	// Match other regexp
	matcher, err = Compile("/{name}Abc/{age}")
	OnErrPanic(err)

	tt.AssertEq("Match77", "/(?P<name>[^/]*)Abc/(?P<age>[^/]*)", matcher.Pattern())
	matchMap, match = matcher.Match("/LosuAbc/123")
	tt.AssertTrue("Match7", match)
	tt.AssertEq("Match8", "Losu", matchMap["name"])
	tt.AssertEq("Match9", "123", matchMap["age"])

	// Match Literal
	matcher, err = Compile("/user/123")
	OnErrPanic(err)
	tt.AssertEq("Match33", "/user/123", matcher.Pattern())
	_, match = matcher.Match("/user/123")
	tt.AssertTrue("Match3", match)
	_, match = matcher.Match("/user")
	tt.AssertFalse("Match4", match)
	_, match = matcher.Match("/user/12")
	tt.AssertFalse("Match5", match)
}
