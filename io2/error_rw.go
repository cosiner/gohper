package io2

import (
	"io"

	"github.com/cosiner/gohper/unsafe2"
)

type ErrorReader struct {
	io.Reader
	Error error
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
	return i, r.Error
}

func (r *ErrorReader) ClearError() {
	r.Error = nil
}

type ErrorWriter struct {
	io.Writer
	Error error
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
	var i int
	if w.Error == nil {
		i, w.Error = w.Writer.Write(data)
		if f != nil {
			f(data)
		}
	}
	return i, w.Error
}

func (w *ErrorWriter) ClearError() {
	w.Error = nil
}
