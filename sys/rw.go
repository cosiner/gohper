package sys

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/cosiner/golib/types"
)

var _newLine = []byte("\n")

// WriteBytesln write bytes to writer and append a newline character
func WriteBytesln(w io.Writer, bs []byte) (int, error) {
	c, e := w.Write(bs)
	if e == nil {
		_, e = w.Write(_newLine)
		if e == nil {
			c += 1
		}
	}
	return c, e
}

// WriteBytesln write string to writer and append a newline character
func WriteStrln(w io.Writer, s string) (int, error) {
	return WriteBytesln(w, types.UnsafeBytes(s))
}

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

// WriteVString write slice of string
func (w BufVWriter) WriteVString(strs []string) (int, error) {
	return filterVString(func(index int, str string) (int, error) {
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
	return filterVBytes(func(index int, b []byte) (int, error) {
		return w.Write(b)
	}, bs)
}

// WriteL write list of []byte
func (w BufVWriter) WriteL(bs ...[]byte) (int, error) {
	return w.WriteV(bs)
}
func filterVString(filter func(int, string) (int, error), slice []string) (n int, err error) {
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

func filterVBytes(filter func(int, []byte) (int, error), slice [][]byte) (n int, err error) {
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

// FileOverwrite delete all content in file, and write new content to it
func FileOverwrite(fname string, content string) error {
	return OpenAndTruncFor(fname, func(fd *os.File) (err error) {
		_, err = fd.WriteString(content)
		return
	})
}

// ReadFirstLine read first line from file
func ReadFirstLine(fname string) (line string, err error) {
	err = FilterFileContent(fname, false, func(_ int, l []byte) ([]byte, error) {
		line = string(l)
		return nil, io.EOF
	})
	return
}

// ReadOneLineFrom read first line from reader, if no content return errNoContent
func ReadOneLineFrom(r io.Reader) (line string, err error) {
	err = FilterContent(r, nil, false, func(_ int, l []byte) ([]byte, error) {
		line = string(l)
		return nil, io.EOF
	})
	return
}

// FilterFileContent filter file content with given filter
// if rewrite is true, content after filter will overwrite exist file
// if recieved io.EOF or other error, will stop read next line,
// and not rewrite file, io.EOF means an early end.
func FilterFileContent(fname string, rewrite bool,
	filter func(int, []byte) ([]byte, error)) error {
	var openfunc func(string, FileOpFunc) error
	if rewrite {
		openfunc = OpenForRW
	} else {
		openfunc = OpenForRead
	}
	return openfunc(fname, func(fd *os.File) (err error) {
		var wr io.Writer
		if rewrite {
			wr = fd
		}
		return FilterContent(fd, wr, false, filter)
	})
}

// FilterFileContentTo filter file content with given filter, then write result
// to dest file
func FilterFileContentTo(fname, dstfile string, delContent bool,
	filter func(int, []byte) ([]byte, error)) error {
	return OpenForRead(fname, func(fd *os.File) (err error) {
		return OpenOrCreateFor(dstfile, delContent,
			func(dstFd *os.File) error {
				return FilterContent(fd, dstFd, true, filter)
			})
	})
}

// FilterLine filter line from reader
// if want to stop filter, just return an io.EOF and it will accept as end
// if filter return any other error, it will just be returned to caller
func FilterLine(rd io.Reader, filter func(int, []byte) error) error {
	return FilterContent(rd, nil, false, func(linenum int, line []byte) ([]byte, error) {
		return nil, filter(linenum, line)
	})
}

// FilterContent readline from reader, after filter, if parallel set at the end, and writer is non-null
// write content to writer, otherwise,  every read operation will followed by
// a write operation
func FilterContent(rd io.Reader, wr io.Writer,
	parallel bool, filter func(int, []byte) ([]byte, error)) (err error) {
	var (
		saveToBuffer = func(line []byte) {}
		flushBuffer  = func() error { return nil }
		line         []byte
		linenum      int
		br           = BufReader(rd)
		earlyStop    bool
	)
	if wr != nil { // writer non-null means need out put filter line to writer
		var bw interface {
			Write([]byte) (int, error)
			WriteString(string) (int, error)
		}
		saveToBuffer = func(line []byte) {
			if line != nil {
				bw.Write(line)
				bw.WriteString("\n")
			}
		}
		if parallel { // use buffered writer, don't need extra buffer
			w := BufWriter(wr)
			flushBuffer = func() error {
				return w.Flush()
			}
			bw = w
		} else { // if not parallel, need buffer to save content for later write
			w := bytes.NewBuffer(make([]byte, 0, FileBufferSIze))
			flushBuffer = func() error {
				TruncateAndSeekWriter(wr)
				_, err := wr.Write(w.Bytes())
				return err
			}
			bw = w
		}
	}
	if filter == nil {
		filter = func(_ int, line []byte) ([]byte, error) {
			return line, nil
		}
	}
	for {
		if line, _, err = br.ReadLine(); err == nil {
			linenum++
			if line, err = filter(linenum, line); err == nil {
				saveToBuffer(line)
				continue
			}
			earlyStop = true
		}
		break
	}
	if !earlyStop {
		flushBuffer()
	}
	if err == io.EOF {
		err = nil
	}
	return
}
