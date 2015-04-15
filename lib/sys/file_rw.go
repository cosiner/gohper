package sys

import (
	"bytes"
	"io"
	"os"

	rw "github.com/cosiner/gohper/lib/io"
)

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
		br           = rw.BufReader(rd)
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
			w := rw.BufWriter(wr)
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
