package hashcode

import "testing"

var data = [][]byte{
	[]byte("abcdefghijklmnopqrstuvwxyz1234567890"),
}

func BenchmarkSDBM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SDBM(data[0])
	}
}
func BenchmarkRS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RS(data[0])
	}
}
func BenchmarkJS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JS(data[0])
	}
}
func BenchmarkELF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ELF(data[0])
	}
}
func BenchmarkBKDR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BKDR(data[0])
	}
}
func BenchmarkDJB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DJB(data[0])
	}
}
func BenchmarkAP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AP(data[0])
	}
}
