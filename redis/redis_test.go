package redis

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestRedis(t *testing.T) {
	cache, err := NewRedisStore("addr='127.0.0.1:6379'")
	if err != nil {
		t.Log(err)
		return
	}
	cache.Set("test", "test")
	cache.Set("test", struct{ Name string }{"aaa"})
	cache.Incr("test")
	t.Log(redis.Bytes(cache.Get("test")))
}
