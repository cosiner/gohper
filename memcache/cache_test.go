package memcache

import (
	"testing"

	"github.com/cosiner/golib/test"
)

func TestRandomCache(t *testing.T) {
	cache := Cacher(RANDOM, 1)
	cache.Set("aaa", "sss")
	test.AssertEq(t, "RandomCache", "sss", cache.Get("aaa").(string))

	cache.Set("bbb", "bbb")
	test.AssertEq(t, "RandomCache", "bbb", cache.Get("bbb").(string))
	test.AssertEq(t, "RandomCache", nil, cache.Get("aaa"))

	test.AssertFalse(t, "RandomCache", cache.Update("dddd", 133))
}

func TestLRUCache(t *testing.T) {
	tt := test.WrapTest(t)
	cache := Cacher(LRU, 3)
	cache.Set("a", "a")
	cache.Set("b", "b")
	cache.Set("c", "c")
	tt.AssertEq("LRUCACHE", 3, cache.Len())
	tt.AssertEq("LRUCACHE", 3, cache.Cap())
	tt.AssertEq("LRUCache", "a", cache.Get("a").(string))

	cache.Set("d", "nc")
	tt.AssertEq("LRUCache", nil, cache.Get("b"))
	cache.Set("e", "nc")
	tt.AssertEq("LRUCache", nil, cache.Get("c"))
}
