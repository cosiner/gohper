package io2

import (
	"bufio"
	"io"

	"github.com/cosiner/gohper/unsafe2"
)

// BufReader return a new bufio.Reader from exist io.Reader
// if current reader is already bufferd, return itself
func BufReader(r io.Reader) *bufio.Reader {
	if r, is := r.(*bufio.Reader); is {
		return r
	}
	return bufio.NewReader(r)
}

// BufWriter return a new bufio.Writer from exist io.Writer
// if current Writer is already bufferd, return itself
func BufWriter(w io.Writer) *bufio.Writer {
	if w, is := w.(*bufio.Writer); is {
		return w
	}
	return bufio.NewWriter(w)
}

var _newLine = []byte("\n")

func WriteIfString(w io.Writer, v interface{}) (bool, error) {
	switch s := v.(type) {
	case string:
		_, err := w.Write(unsafe2.Bytes(s))
		return true, err
	case []byte:
		_, err := w.Write(s)
		return true, err
	}
	return false, nil
}

// Writeln write bytes to writer and append a newline character
func Writeln(w io.Writer, bs []byte) (int, error) {
	c, e := w.Write(bs)
	if e == nil {
		_, e = w.Write(_newLine)
		if e == nil {
			c++
		}
	}
	return c, e
}

// WriteStringln write string to writer and append a newline character
func WriteStringln(w io.Writer, s string) (int, error) {
	return Writeln(w, unsafe2.Bytes(s))
}

// WriteL write a string list to writer, return total bytes writed
func WriteLString(w io.Writer, strs ...string) (n int, err error) {
	var c int
	for i := range strs {
		if c, err = w.Write(unsafe2.Bytes(strs[i])); err == nil {
			n += c
		} else {
			break
		}
	}
	return
}

// WriteL write a bytes list to writer, return total bytes writed
func WriteL(w io.Writer, bs ...[]byte) (n int, err error) {
	var c int
	for i := range bs {
		if c, err = w.Write(bs[i]); err == nil {
			n += c
		} else {
			break
		}
	}
	return
}
