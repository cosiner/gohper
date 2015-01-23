package cache

import (
	"testing"

	"github.com/cosiner/golib/test"
)

func TestRandomCache(t *testing.T) {
	cache, _ := Cacher(RANDOM, "maxsize=1")
	cache.Set("aaa", "sss")
	test.AssertEq(t, "RandomCache1", "sss", cache.Get("aaa").(string))

	cache.Set("bbb", "bbb")
	test.AssertEq(t, "RandomCache2", "bbb", cache.Get("bbb").(string))
	test.AssertEq(t, "RandomCache5", 1, cache.Cap())
	test.AssertEq(t, "RandomCache5", 1, cache.Len())
	test.AssertEq(t, "RandomCache3", nil, cache.Get("aaa"))

	test.AssertFalse(t, "RandomCache4", cache.Update("dddd", 133))
}

func TestLRUCache(t *testing.T) {
	tt := test.WrapTest(t)
	cache, _ := Cacher(LRU, "maxsize=3")
	cache.Set("a", "a")
	cache.Set("b", "b")
	cache.Set("c", "c")
	tt.AssertEq("LRUCACHE1", 3, cache.Len())
	tt.AssertEq("LRUCACHE2", 3, cache.Cap())
	tt.AssertEq("LRUCache3", "a", cache.Get("a").(string))

	cache.Set("d", "nc")
	tt.AssertEq("LRUCache4", nil, cache.Get("b"))
	cache.Set("e", "nc")
	tt.AssertEq("LRUCache5", nil, cache.Get("c"))
}

func TestRedisCache(t *testing.T) {
	tt := test.WrapTest(t)
	cache, err := Cacher(REDIS, "addr=127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	cache.Set("aaa", "123")
	tt.Log(cache.Get("aaa"))
}
