package cache

import (
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestRandomCache(t *testing.T) {
	cache, _ := NewCache(RANDOM, "maxsize=1")
	cache.Set("aaa", "sss")
	test.Eq(t, "sss", cache.Get("aaa").(string))

	cache.Set("bbb", "bbb")
	test.Eq(t, "bbb", cache.Get("bbb").(string))
	test.Eq(t, 1, cache.Cap())
	test.Eq(t, 1, cache.Size())
	test.Eq(t, nil, cache.Get("aaa"))

	test.False(t, cache.Update("dddd", 133))
}

func TestLRUCache(t *testing.T) {
	tt := test.Wrap(t)
	cache, _ := NewCache(LRU, "maxsize=3")
	cache.Set("a", "a")
	cache.Set("b", "b")
	cache.Set("c", "c")
	tt.Eq(3, cache.Size())
	tt.Eq(3, cache.Cap())
	tt.Eq("a", cache.Get("a").(string))

	cache.Set("d", "nc")
	tt.Eq(nil, cache.Get("b"))
	cache.Set("e", "nc")
	tt.Eq(nil, cache.Get("c"))
}

func TestRedisCache(t *testing.T) {
	tt := test.Wrap(t)
	cache, err := NewCache(REDIS, "addr=127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	cache.Set("aaa", "123")
	tt.Log(cache.Get("aaa"))
}
