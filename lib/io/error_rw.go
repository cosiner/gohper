package io

import (
	"io"

	"github.com/cosiner/gohper/lib/types"
)

var (
	Bytes  = types.UnsafeBytes
	String = types.UnsafeString
)

func ErrWrite(err error, w io.Writer, data []byte) (int, error) {
	if err != nil {
		return 0, err
	}
	return w.Write(data)
}

func ErrWriteString(err error, w io.Writer, data string) (int, error) {
	if err != nil {
		return 0, err
	}
	return w.Write(Bytes(data))
}

func ErrRead(err error, r io.Reader, data []byte) (int, error) {
	if err != nil {
		return 0, err
	}
	return r.Read(data)
}

func ErrPtrWrite(err *error, w io.Writer, data []byte) (count int) {
	if err != nil && *err == nil {
		count, *err = w.Write(data)
	}
	return
}

func ErrPtrWriteString(err *error, w io.Writer, data string) (count int) {
	if err != nil && *err == nil {
		count, *err = w.Write(Bytes(data))
	}
	return
}

func ErrPtrRead(err *error, r io.Reader, data []byte) (count int) {
	if err != nil && *err == nil {
		count, *err = r.Read(data)
	}
	return
}
