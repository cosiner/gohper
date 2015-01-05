package sys

import (
	"bufio"
	"io"
	"mlib/util/funcs"
	. "mlib/util/generic"
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

// BufReader return a new bufio.Writer from exist io.Writer
// if current Writer is already bufferd, return itself
func BufWriter(wr io.Writer) (bw *bufio.Writer) {
	if wr != nil {
		if w, is := wr.(*bufio.Writer); is {
			bw = w
		} else {
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

// WrapBufwriter wrap bufio.Writer to BufVWriter
func WrapBufWriter(bw *bufio.Writer) BufVWriter {
	return BufVWriter{bw}
}

// WrapBufwriter wrap io.Writer to BufVWriter
func WrapWriter(wr io.Writer) BufVWriter {
	return BufVWriter{BufWriter(wr)}
}

// WriteVString write slice string
//go:generate gotmpl  -i $GOFILE -o ./$GOFILE -t T:string
//go:generate gotmpl -p mlib/util/funcs -o funcs_string_gen.go  -t T:string -d
func (w BufVWriter) WriteVString(strs []string) (int, error) {
	return writeVString(func(index int, str string) (int, error) {
		return w.WriteString(str)
	}, strs)
}

// WriteLString write list of string
func (w BufVWriter) WriteLString(strs ...string) (int, error) {
	return w.WriteVString(strs)
}

// WriteV write slice byte array
//go:generate gotmpl  -i $GOFILE -o ./$GOFILE  -t T:[]byte
//go:generate gotmpl -p mlib/util/funcs -o funcs_bytes_gen.go -t T:[]byte -d
func (w BufVWriter) WriteV(bs [][]byte) (int, error) {
	return writeVBytes(func(index int, b []byte) (int, error) {
		return w.Write(b)
	}, bs)
}

// WriteL write list of []byte
func (w BufVWriter) WriteL(bs ...[]byte) (int, error) {
	return w.WriteV(bs)
}

// writeV integrates common operation for batch write
func writeV_T(fn func(int, T) (int, error), slice []T) (n int, err error) {
	err = funcs.MapWithErrFor_T(slice, func(index int, o T) error {
		var m int
		if m, err = fn(index, o); err == nil {
			n += m
		}
		return err
	})
	return
}

// writeV integrates common operation for batch write
func writeVString(fn func(int, string) (int, error), slice []string) (n int, err error) {
	err = funcs.MapWithErrForString(slice, func(index int, o string) error {
		var m int
		if m, err = fn(index, o); err == nil {
			n += m
		}
		return err
	})
	return
}

// writeV integrates common operation for batch write
func writeVBytes(fn func(int, []byte) (int, error), slice [][]byte) (n int, err error) {
	err = funcs.MapWithErrForBytes(slice, func(index int, o []byte) error {
		var m int
		if m, err = fn(index, o); err == nil {
			n += m
		}
		return err
	})
	return
}
