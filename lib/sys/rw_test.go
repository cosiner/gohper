package sys

import (
	"bytes"
	"testing"

	"github.com/cosiner/gohper/lib/test"
)

func TestBufReader(t *testing.T) {

}

func TestBufWriter(t *testing.T) {

}

func TestWrapBufWriter(t *testing.T) {

}

func TestWrapWriter(t *testing.T) {

}

func TestWriteVString(t *testing.T) {

}

func TestWriteLString(t *testing.T) {

}

func TestWriteV(t *testing.T) {

}

func TestWriteL(t *testing.T) {
	tt := test.Wrap(t)
	buffer := bytes.NewBuffer([]byte{})
	bw := NewBufVWriter(buffer)
	bw.WriteLString("abc", "de")
	bw.Flush()
	tt.True(bytes.Equal([]byte("abcde"), buffer.Bytes()))
}
