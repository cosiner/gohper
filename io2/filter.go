package io2

import (
	"bytes"
	"io"
)

// FileBufferSIze is the buffer size of slice to store file content
const FileBufferSIze = 4096

type LineFilterFunc func(num int, line []byte) (newLine []byte, err error)

// NopLineFilter do nothing, just return the line
func NopLineFilter(_ int, line []byte) ([]byte, error) {
	return line, nil
}

// FilterRead filter line from reader
func FilterRead(r io.Reader, filter func(int, []byte) error) error {
	return Filter(r, nil, false, func(linenum int, line []byte) ([]byte, error) {
		return nil, filter(linenum, line)
	})
}

// Filter readline from reader, after filter, if sync set at the end, and writer is non-null,
// write content to writer, otherwise,  every read operation will followed by
// a write operation
//
// return an error to stop filter, io.EOF means normal stop
func Filter(r io.Reader, w io.Writer, sync bool, filter LineFilterFunc) error {
	var (
		save  = func([]byte) {}
		flush = func() error { return nil }
	)

	if w != nil {
		var bufw interface {
			Write([]byte) (int, error)
			WriteString(string) (int, error)
		}

		if sync { // use bufio.Writer
			bw := BufWriter(w)
			flush = func() error {
				return bw.Flush()
			}
			bufw = bw
		} else { // use bytes.Buffer
			bw := bytes.NewBuffer(make([]byte, 0, FileBufferSIze))
			flush = func() error {
				_, err := w.Write(bw.Bytes())
				return err
			}
			bufw = bw
		}

		save = func(line []byte) {
			bufw.Write(line)
			bufw.WriteString("\n")
		}
	}

	if filter == nil {
		filter = NopLineFilter
	}

	var (
		line    []byte
		linenum int
		err     error
		br      = BufReader(r)
	)

	for err == nil {
		if line, _, err = br.ReadLine(); err == nil {
			linenum++
			if line, err = filter(linenum, line); err == nil || err == io.EOF {
				if line != nil {
					save(line)
				}
			}
		}
	}

	if err == io.EOF {
		flush()
		err = nil
	}

	return err
}
