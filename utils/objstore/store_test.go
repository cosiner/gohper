package objstore

import (
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkPut(b *testing.B) {
	// tt := testing2.Wrap(b)
	s := New(0, 0)
	for i := 0; i < b.N; i++ {
		si := strconv.Itoa(i)
		s.Put(si, Object{
			Value: si,
		})
	}
	b.Log(s.Size())
}

func initStore(size int) *Store {
	s := New(0, uint(size))

	for i := 0; i < size; i++ {
		si := strconv.Itoa(i)
		s.Put(si, Object{
			Value: si,
		})
		if i&1 == 0 {
			s.Remove(si)
		}
	}
	fmt.Println(len(s.indexes))
	return s
}

var s = initStore(20000)

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s.Get(strconv.Itoa(i))
	}
}
