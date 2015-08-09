package strings2

import (
	"bytes"
	"sort"
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestSnakeCase(t *testing.T) {
	testing2.
		Expect("_xy_xy").Arg("_xy_xy").
		Expect("_xy_xy").Arg("_xy_xy").
		Expect("_xy_xy").Arg("_xyXy").
		Expect("_xy xy").Arg("_Xy Xy").
		Expect("_xy_xy").Arg("_Xy_Xy").
		Run(t, ToSnake)
}

func TestCamelString(t *testing.T) {
	testing2.
		Expect("XyXy").Arg("xy_xy").
		Expect("Xy__Xy").Arg("xy__Xy").
		Expect("Xy Xy").Arg("xy Xy").
		Expect("XY Xy").Arg("x_y Xy").
		Expect("X_Y XY").Arg("x__y XY").
		Expect("XY XY").Arg("x_y xY").
		Expect("XY XY").Arg("x_y _x_y").
		Expect("  XY").Arg("  x_y").
		Run(t, ToCamel)
}

func TestAbridgeString(t *testing.T) {
	testing2.
		Expect("ABC").Arg("AaaBbbCcc").
		Expect("ABC").Arg("AaaBbbCcc").
		Expect("").Arg("").
		Run(t, ToAbridge)

	testing2.
		Expect("abc").Arg("AaaBbbCcc").
		Expect("abc").Arg("AaaBbbCcc").
		Expect("").Arg("").
		Run(t, ToLowerAbridge)
}

func TestTrimQuote(t *testing.T) {
	testing2.
		Expect("aaa", true).Arg("\"aaa\"").
		Expect("aaa", true).Arg("'aaa'").
		Expect("aaa", true).Arg("`aaa`").
		Expect("", true).Arg("").
		Expect("", testing2.NonNil).Arg("`aaa").
		Run(t, TrimQuote)
}

func TestSplitAtN(t *testing.T) {
	testing2.
		Expect(3).Arg("123123123", "12", 2).
		Expect(6).Arg("123123123", "12", 3).
		Expect(-1).Arg("123123123", "12", 4).
		Run(t, IndexN)
}

func TestSplitAtLastN(t *testing.T) {
	testing2.
		Expect(6).Arg("123123123", "12", 1).
		Expect(3).Arg("123123123", "12", 2).
		Expect(0).Arg("123123123", "12", 3).
		Expect(-1).Arg("123123123", "12", 4).
		Run(t, LastIndexN)
}

func TestJoin(t *testing.T) {
	testing2.
		Expect("abc,abc,abc").Arg("abc", ",", 3).
		Expect("abc,abc").Arg("abc", ",", 2).
		Expect("abc").Arg("abc", ",", 1).
		Expect("").Arg("abc", ",", 0).
		Run(t, RepeatJoin)

	testing2.
		Expect("").Arg([]int{}, ",").
		Expect("1,2,3,4").Arg([]int{1, 2, 3, 4}, ",").
		Run(t, JoinInt)

	testing2.
		Expect("").Arg([]uint{}, ",").
		Expect("1,2,3,4").Arg([]uint{1, 2, 3, 4}, ",").
		Run(t, JoinUint)

	testing2.
		Expect("").Arg([]string{}, "=?", ", ").
		Expect("a=?").Arg([]string{"a"}, "=?", ", ").
		Expect("a=?, b=?, c=?").Arg([]string{"a", "b", "c"}, "=?", ", ").
		Run(t, SuffixJoin)

}

func TestValid(t *testing.T) {
	testing2.Tests().
		True().Arg("", "abcdefghijklmn").
		True().Arg("abc", "abcdefghijklmn").
		False().Arg("ao", "abcdefghijklmn").
		Run(t, IsAllCharsIn)
}

func TestRemoveSpace(t *testing.T) {
	testing2.
		Expect("abcdefg").Arg(`a b
    	c d 	e
    	 	f g`).
		Run(t, RemoveSpace)
}

func TestMergeSpace(t *testing.T) {
	testing2.
		Expect("a b c dd").Arg("   a    b   c  dd   ", true).
		Expect(" a b c dd ").Arg("   a    b   c  dd   ", false).
		Run(t, MergeSpace)
}

func TestIndexNonSpace(t *testing.T) {
	testing2.
		Expect(1).Arg(" 1             ").
		Expect(-1).Arg("             ").
		Expect(0).Arg("a   ").
		Run(t, IndexNonSpace).
		Run(t, LastIndexNonSpace)
}

func TestTrim(t *testing.T) {
	tt := testing2.Wrap(t)

	tt.
		Eq("0123", TrimAfter("01234567", "45")).
		Eq("4567", TrimBefore("01234567", "23"))

	left, right := ("["), ("]")
	testing2.
		Expect("123", true).Arg("[123]", left, right, true).
		Expect(testing2.NoCheck, false).Arg("[123", left, right, true).
		Expect("123", true).Arg("[123", left, right, false).
		Expect("123", true).Arg("123", left, right, true).
		Expect("", true).Arg("", left, right, true).
		Run(t, TrimWrap)

	testing2.
		Expect([]string{"1", "2", "3"}).Arg(" 1 , 2  , 3   ", ",").
		Run(t, SplitAndTrim)

	tt.Eq("abcde", TrimAndToLower("  ABCDE "))
	tt.Eq("ABCDE", TrimAndToUpper("  abcde "))
}

func TestCompare(t *testing.T) {
	testing2.
		Expect(1).Arg("123", "012").
		Expect(-1).Arg("123", "212").
		Expect(0).Arg("123", "123").
		Expect(-1).Arg("123", "1234").
		Expect(1).Arg("123", "12").
		Run(t, Compare)
}

func TestLastIndexByte(t *testing.T) {
	testing2.
		Expect(0).Arg("123", byte('1')).
		Expect(-1).Arg("123", byte('4')).
		Run(t, LastIndexByte)
}

func TestMidSep(t *testing.T) {
	testing2.
		Expect(-1).Arg("123", byte('1')).
		Expect(-1).Arg("123", byte('3')).
		Expect(-1).Arg("123", byte('4')).
		Expect(1).Arg("123", byte('2')).
		Run(t, MidIndex)

	testing2.
		Expect("", "").Arg("123", byte('1')).
		Expect("", "").Arg("123", byte('3')).
		Expect("", "").Arg("123", byte('4')).
		Expect("1", "3").Arg("123", byte('2')).
		Run(t, Seperate)
}

func TestWriteBuffer(t *testing.T) {
	tt := testing2.Wrap(t)
	buf := bytes.NewBuffer(make([]byte, 0, 10))
	WriteStringsToBuffer(buf, []string{"a", "b", "c"}, ",")
	tt.Eq("a,b,c", buf.String())
}

func TestNumMatched(t *testing.T) {
	testing2.
		Expect(0).Arg(IsNotEmpty, "", "", "", "").
		Expect(1).Arg(IsNotEmpty, "1", "", "", "").
		Expect(2).Arg(IsNotEmpty, "1", "2", "", "").
		Run(t, NumMatched)

	testing2.
		Expect(4).Arg(IsEmpty, "", "", "", "").
		Expect(3).Arg(IsEmpty, "1", "", "", "").
		Expect(2).Arg(IsEmpty, "1", "2", "", "").
		Run(t, NumMatched)
}

func TestFilterMap(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.DeepEq([]string{"1"}, Filter(IsNotEmpty, "", "", "1", ""))
	tt.DeepEq([]string{"", "", ""}, Filter(IsEmpty, "", "", "1", ""))

	tt.DeepEq([]string{"empty", "empty", "1", "empty"},
		Map(func(s string) string {
			if s == "" {
				return "empty"
			}
			return s
		}, "", "", "1", ""))

}

func TestFilterInStrings(t *testing.T) {
	tt := testing2.Wrap(t)

	strings := []string{"", "A", "", "B", "", "C"}
	strings2 := FilterInPlace(IsNotEmpty, strings...)
	tt.DeepEq([]string{"A", "B", "C"}, strings2)
	tt.DeepEq([]string{"A", "B", "C", "B", "", "C"}, strings)

	strings2 = ClearEmpty(strings)
	tt.DeepEq([]string{"A", "B", "C", "B", "C"}, strings2)
}

func TestMultipleLineOperate(t *testing.T) {
	tt := testing2.Wrap(t)
	json := `{
	"aa":"bb", //ddd
	"ba":"bb", //ccc
	"ca":"bb" //zzz
}`
	json = MultipleLineOperate(json, "//", TrimAfter)
	tt.Eq(`{
"aa":"bb",
"ba":"bb",
"ca":"bb"
}`, json)
}

func TestMakeSlice(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.DeepEq([]string{}, MakeSlice("0", 0))
	tt.DeepEq([]string{"0"}, MakeSlice("0", 1))
	tt.DeepEq([]string{"0", "0"}, MakeSlice("0", 2))
}

func TestRandString(t *testing.T) {
	tt := testing2.Wrap(t)
	slice := []string{"1", "2", "3", "4"}
	tt.True(sort.SearchStrings(slice, RandIn(slice)) >= 0)
	tt.True(sort.SearchStrings(slice, RandIn(slice)) >= 0)
	tt.True(sort.SearchStrings(slice, RandIn(slice)) >= 0)
	tt.True(sort.SearchStrings(slice, RandIn(slice)) >= 0)
	tt.True(RandIn(nil) == "")
	tt.True(RandIn([]string{}) == "")
}

func TestTrimN(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.Eq("aaa", TrimFirstN("///aaa", "/", 0))
	tt.Eq("aaa", TrimFirstN("///aaa", "/", -1))
	tt.Eq("//aaa", TrimFirstN("///aaa", "/", 1))

	tt.Eq("aaa", TrimLastN("aaa///", "/", 0))
	tt.Eq("aaa", TrimLastN("aaa///", "/", -1))
	tt.Eq("aaa//", TrimLastN("aaa///", "/", 1))
	tt.Eq("aaa/", TrimLastN("aaa///", "/", 2))
	tt.Eq("aaa", TrimLastN("aaa///", "/", 3))
	tt.Eq("aaa", TrimLastN("aaa///", "/", 4))
}

func TestRemoveDuplicate(t *testing.T) {
	tt := testing2.Wrap(t)

	tt.DeepEq([]string{"aa", "bb", "cc"}, RemoveDuplicate([]string{"cc", "aa", "cc", "bb", "aa"}))
	tt.DeepEq([]string{}, RemoveDuplicate([]string{}))
}

func TestSearch(t *testing.T) {
	tt := testing2.Wrap(t)

	strings := []string{"a", "b", "c", "d", "e"}
	tt.Eq(1, Search(strings, "b", 0))
	tt.Eq(-1, Search(strings, "b", 1))
}
