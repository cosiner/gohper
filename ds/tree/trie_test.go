package tree

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestTrieTree(t *testing.T) {
	tt := testing2.Wrap(t)
	tree := Trie{}
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
var tree = func() *Trie {
	t := Trie{}
	for i := 0; i < len(s); i++ {
		t.AddPath(s[i:]+s[:i], true)
	}
	return &t
}()

func BenchmarkTrieTree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if tree.MatchValue(s).(bool) != true {
			b.Fail()
		}
	}
}

func TestPrefix(t *testing.T) {
	tt := testing2.Wrap(t)

	tree := Trie{}

	tree.AddPath("1234", 1)
	tree.AddPath("234", 2)
	tree.AddPath("12", 3)
	tree.AddPath("347", 4)
	tree.AddPath("00", 5)

	tt.Nil(tree.PrefixMatchValue(""))
	tt.Nil(tree.PrefixMatchValue("1"))
	tt.Nil(tree.PrefixMatchValue("2"))
	tt.Nil(tree.PrefixMatchValue("3"))
	tt.Nil(tree.PrefixMatchValue("0"))
	tt.Nil(tree.PrefixMatchValue("13"))
	tt.Nil(tree.PrefixMatchValue("01"))

	tt.Eq(1, tree.PrefixMatchValue("1234").(int))
	tt.Eq(1, tree.PrefixMatchValue("12345").(int))

	tt.Eq(3, tree.PrefixMatchValue("12").(int))
	tt.Eq(3, tree.PrefixMatchValue("123").(int))
	tt.Eq(3, tree.PrefixMatchValue("124").(int))
}
