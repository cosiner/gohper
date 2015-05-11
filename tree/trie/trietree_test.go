package trie

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestTrieTree(t *testing.T) {
	tt := testing2.Wrap(t)
	tree := new(TrieTree)
	tree.AddPath("abcde", 123)
	tree.AddPath("bcdef", 234)
	tree.AddPath("efghi", 456)
	tree.AddPath("fghij", 789)
	tt.True(tree.MatchValue("efghi") != nil)
	tt.Eq(tree.MatchValue("abcde").(int), 123)
	tt.Eq(tree.MatchValue("bcdef").(int), 234)
	tt.Eq(tree.MatchValue("efghi").(int), 456)
	tt.Eq(tree.MatchValue("fghij").(int), 789)
	tt.Eq(tree.MatchValue("fghia"), nil)
	tt.Eq(tree.MatchValue("fasdahia"), nil)
	tt.Eq(tree.MatchValue("csaghia"), nil)
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
		if tree.MatchValue(s).(bool) != true {
			b.Fail()
		}
	}
}
