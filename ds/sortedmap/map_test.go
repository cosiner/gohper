package sortedmap

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestMap(t *testing.T) {
	tt := testing2.Wrap(t)

	mp := New()
	kvs := map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
		"D": 4,
		"E": 5,
		"F": 6,
	}

	for k, v := range kvs {
		mp.Set(k, v)
	}
	mp.Set("A", 7)
	for k, v := range kvs {
		val := mp.Get(k).(int)
		if k == "A" {
			v = 7
		}

		tt.Eq(v, val)
	}

	tt.False(mp.HasKey("G"))
	tt.Nil(mp.Get("G"))
	mp.Delete("G")
	tt.Eq(8, mp.DefGet("G", 8).(int))
	tt.Eq(7, mp.DefGet("A", 7).(int))

	mp.Delete("A")
	delete(kvs, "A")
	for k, v := range kvs {
		val := mp.Get(k).(int)
		if k == "A" {
			v = 7
		}

		tt.Eq(v, val)
	}

	mp.Clear()
	tt.DeepEq(make(map[string]int), mp.Indexes)
	tt.DeepEq([]Element{}, mp.Values)
}
