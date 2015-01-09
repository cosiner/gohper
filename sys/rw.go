package sys

import (
	"bufio"
	"io"
)

// BufReader return a new bufio.Reader from exist io.Reader
// if current reader is already bufferd, return itself
func BufReader(rd io.Reader) (br *bufio.Reader) {
	if rd != nil {
		if r, is := rd.(*bufio.Reader); is {
			br = r
		} else {
			br = bufio.NewReader(rd)
		}
	}
	return
}

// BufWriter return a new bufio.Writer from exist io.Writer
// if current Writer is already bufferd, return itself
func BufWriter(wr io.Writer) (bw *bufio.Writer) {
	if wr != nil {
		switch wr := wr.(type) {
		case *bufio.Writer:
			bw = wr
		default:
			bw = bufio.NewWriter(wr)
		}
	}
	return
}

// BufVWriter is a wrapper of bufio.Writer that supply functions to write
// a slice of string or byte array for batch write
type BufVWriter struct {
	*bufio.Writer
}

// NewBufVWriter wrap io.Writer to BufVWriter
func NewBufVWriter(wr io.Writer) BufVWriter {
	return BufVWriter{BufWriter(wr)}
}

// Filter write slice string
//go:generate gotmpl -p "github.com/cosiner/golib/types" -o ./$GOFILE -f FilterV -t "T:string"
func (w BufVWriter) WriteVString(strs []string) (int, error) {
	return FilterVString(func(index int, str string) (int, error) {
		return w.WriteString(str)
	}, strs)
}

// WriteLString write list of string
func (w BufVWriter) WriteLString(strs ...string) (int, error) {
	return w.WriteVString(strs)
}

// WriteV write slice byte array
//go:generate gotmpl -p "github.com/cosiner/golib/types" -o ./$GOFILE -f FilterV -t "T:[]byte]"
func (w BufVWriter) WriteV(bs [][]byte) (int, error) {
	return FilterVBytes(func(index int, b []byte) (int, error) {
		return w.Write(b)
	}, bs)
}

// WriteL write list of []byte
func (w BufVWriter) WriteL(bs ...[]byte) (int, error) {
	return w.WriteV(bs)
}

func FilterVString(filter func(int, string) (int, error), slice []string) (n int, err error) {
	var m int
	for index, s := range slice {
		if m, err = filter(index, s); err == nil {
			n += m
		} else {
			break
		}
	}
	return
}

func FilterVBytes(filter func(int, []byte) (int, error), slice [][]byte) (n int, err error) {
	var m int
	for index, s := range slice {
		if m, err = filter(index, s); err == nil {
			n += m
		} else {
			break
		}
	}
	return
}
