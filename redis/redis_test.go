package redis

import (
	"testing"

	. "github.com/cosiner/golib/errors"

	"github.com/cosiner/golib/test"
)

func TestRedis(t *testing.T) {
	tt := test.WrapTest(t)
	store, err := NewRedisStore2("addr='127.0.0.1:6379'")
	if err != nil {
		t.Log(err)
		return
	}
	store.Set("test", "123")
	store.Incr("test")
	s, err := ToInt(store.Get("test"))
	OnErrExit(err)
	tt.AssertEq("Redis1", 124, s)
	store.Set("test", struct{ Name string }{"aaa"})
	v, err := ToString(store.Get("test"))
	OnErrExit(err)
	tt.AssertEq("Redis2", "{aaa}", v)
	// store.HSet("userexist", "someone", false)
	tt.Log(store.IsHExist("userexist", "someone"))
	store.HRemove("userexist", "someone")
	tt.Log(store.IsHExist("userexist", "someone"))
	tt.Log(store.HGet("userexist", "someone"))
	// tt.Log(ToBool(store.HGet("userexist", "someone")))
}
