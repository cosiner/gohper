package types

import (
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestTrieTree(t *testing.T) {
	tt := test.Wrap(t)
	tree := new(TrieTree)
	tree.AddPath("abcde", 123)
	tree.AddPath("bcdef", 234)
	tree.AddPath("efghi", 456)
	tree.AddPath("fghij", 789)
	tt.True(tree.Match("efghi") != nil)
	tt.Eq(tree.Match("abcde").(int), 123)
	tt.Eq(tree.Match("bcdef").(int), 234)
	tt.Eq(tree.Match("efghi").(int), 456)
	tt.Eq(tree.Match("fghij").(int), 789)
	tt.Eq(tree.Match("fghia"), nil)
	tt.Eq(tree.Match("fasdahia"), nil)
	tt.Eq(tree.Match("csaghia"), nil)
}

var s = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890"
var tree = func() *TrieTree {
	t := new(TrieTree)
	for i := 0; i < len(s); i++ {
		t.AddPath(s[i:]+s[:i], true)
	}
	return t
}()

func BenchmarkTrieTree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if tree.Match(s).(bool) != true {
			b.Fail()
		}
	}
}
