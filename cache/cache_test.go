package cache

import (
	"testing"

	"github.com/cosiner/golib/test"
)

func TestRandomCache(t *testing.T) {
	cache, _ := NewCache(RANDOM, "maxsize=1")
	cache.Set("aaa", "sss")
	test.AssertEq(t, "sss", cache.Get("aaa").(string))

	cache.Set("bbb", "bbb")
	test.AssertEq(t, "bbb", cache.Get("bbb").(string))
	test.AssertEq(t, 1, cache.Cap())
	test.AssertEq(t, 1, cache.Size())
	test.AssertEq(t, nil, cache.Get("aaa"))

	test.AssertFalse(t, cache.Update("dddd", 133))
}

func TestLRUCache(t *testing.T) {
	tt := test.WrapTest(t)
	cache, _ := NewCache(LRU, "maxsize=3")
	cache.Set("a", "a")
	cache.Set("b", "b")
	cache.Set("c", "c")
	tt.AssertEq(3, cache.Size())
	tt.AssertEq(3, cache.Cap())
	tt.AssertEq("a", cache.Get("a").(string))

	cache.Set("d", "nc")
	tt.AssertEq(nil, cache.Get("b"))
	cache.Set("e", "nc")
	tt.AssertEq(nil, cache.Get("c"))
}

func TestRedisCache(t *testing.T) {
	tt := test.WrapTest(t)
	cache, err := NewCache(REDIS, "addr=127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	cache.Set("aaa", "123")
	tt.Log(cache.Get("aaa"))
}
