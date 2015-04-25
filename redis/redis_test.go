package redis

import (
	"testing"
)

func TestRedis(t *testing.T) {
	// tt := test.Wrap(t)
	// store, err := New2("addr='127.0.0.1:6379'")
	// if err != nil {
	// 	t.Log(err)
	// 	return
	// }
	// store.Set("test", "123")
	// store.Incr("test")
	// s, err := ToInt(store.Get("test"))
	// OnErrExit(err)
	// tt.Eq(124, s)
	// store.Set("test", struct{ Name string }{"aaa"})
	// v, err := ToString(store.Get("test"))
	// OnErrExit(err)
	// tt.Eq("{aaa}", v)
	// // store.HSet("userexist", "someone", false)
	// tt.Log(store.IsHExist("userexist", "someone"))
	// store.HRemove("userexist", "someone")
	// tt.Log(store.IsHExist("userexist", "someone"))
	// tt.Log(store.HGet("userexist", "someone"))
	// // tt.Log(ToBool(store.HGet("userexist", "someone")))
}
