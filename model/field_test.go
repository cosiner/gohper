package model

import "testing"

func BenchmarkFieldSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewFieldSet()
	}
}
