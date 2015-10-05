package io2

import (
	"bufio"
	"io"

	"github.com/cosiner/gohper/unsafe2"
)

type ErrorReader struct {
	io.Reader
	Error error
	Count int
}

func NonEOF(err error) error {
	if err == io.EOF {
		return nil
	}

	return err
}

func NewErrorReader(r io.Reader) *ErrorReader {
	if er, is := r.(*ErrorReader); is {
		return er
	}

	return &ErrorReader{
		Reader: r,
	}
}

func (r *ErrorReader) Read(data []byte) (int, error) {
	return r.ReadDo(data, nil)
}

func (r *ErrorReader) ReadDo(data []byte, f func([]byte)) (int, error) {
	var i int
	if r.Error == nil {
		i, r.Error = r.Reader.Read(data)
		if f != nil {
			f(data)
		}
	}
	r.Count += i

	return i, r.Error
}

func (r *ErrorReader) ClearError() {
	r.Error = nil
}

type ErrorWriter struct {
	io.Writer
	Error error
	Count int
}

func NewErrorWriter(w io.Writer) *ErrorWriter {
	if ew, is := w.(*ErrorWriter); is {
		return ew
	}

	return &ErrorWriter{
		Writer: w,
	}
}

func (w *ErrorWriter) Write(data []byte) (int, error) {
	return w.WriteDo(data, nil)
}

func (w *ErrorWriter) WriteString(s string) (int, error) {
	return w.WriteDo(unsafe2.Bytes(s), nil)
}

func (w *ErrorWriter) WriteDo(data []byte, f func([]byte)) (int, error) {
	if w.Error != nil {
		return 0, w.Error
	}

	var i int
	i, w.Error = w.Writer.Write(data)
	if f != nil {
		f(data)
	}
	w.Count += i

	return i, w.Error
}

func (w *ErrorWriter) ClearError() {
	w.Error = nil
}

func ErrWriteByte(err *error, w *bufio.Writer, b byte) {
	if *err == nil {
		*err = w.WriteByte(b)
	}
}

func ErrWrite(err *error, w io.Writer, b []byte) {
	if *err == nil {
		_, *err = w.Write(b)
	}
}

func ErrWriteString(err *error, w io.Writer, s string) {
	if *err == nil {
		_, *err = w.Write(unsafe2.Bytes(s))
	}
}
