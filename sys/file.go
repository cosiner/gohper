package sys

import (
	"bytes"
	"io"
	"os"
)

const (
	// FILE_BBUFSIZe is the buffer size of slice to store file content
	FILE_BUFSIZE = 4096
	FILE_WRPERM  = 0644
	DIR_PERM     = 0755
)

// IsExist check whether or not file/dir exist
func IsExist(fname string) bool {
	_, err := os.Stat(ExpandAbs(fname))
	return err == nil
}

// IsDir check whether or not given name is a dir
func IsDir(fname string) bool {
	fi, err := os.Stat(fname)
	return err == nil && fi.IsDir()
}

// IsSymlink check whether or not given name is a symlink
func IsSymlink(fname string) bool {
	fi, err := os.Lstat(fname)
	return err == nil && (fi.Mode()&os.ModeSymlink == os.ModeSymlink)
}

// IsFileModifiedAfter check whether or not file is modified by the function
func IsFileModifiedAfter(fname string, fn func()) bool {
	fi1, err := os.Stat(fname)
	fn()
	fi2, _ := os.Stat(fname)
	return err == nil && !fi1.ModTime().Equal(fi2.ModTime())
}

// OpenOrCreate open or  create file
func OpenOrCreate(fname string) (*os.File, error) {
	return os.OpenFile(fname, os.O_RDWR|os.O_CREATE, FILE_WRPERM)
}

// OpenFileForRead open file, process file with function, then close file
func OpenFileForRead(fname string, fn func(*os.File) error) error {
	return openFileFor(func() (*os.File, error) {
		return os.Open(fname)
	}, fn)
}

// OpenFileFor open file, process file with function, then close file
func OpenFileForModify(fname string, fn func(*os.File) error) error {
	return openFileFor(func() (*os.File, error) {
		return os.OpenFile(fname, os.O_RDWR, FILE_WRPERM)
	}, fn)
}

// openFileFor open file use given opener, if no error, call fn
// opener determin mode and permission of open file
func openFileFor(opener func() (*os.File, error), fn func(*os.File) error) error {
	fd, err := opener()
	if err == nil {
		err = fn(fd)
		fd.Close()
	}
	return err
}

// TruncateAndSeek truncate file size to 0 and seek current positon to 0
func TruncateAndSeek(fd *os.File) {
	if fd != nil {
		fd.Truncate(0)
		fd.Seek(0, os.SEEK_SET)
	}
}

// TruncateAndSeekWriter assume writer is actually os.File,
// do same as TruncateAndSeek, otherwise, do nothing
func TruncateAndSeekWriter(wr io.Writer) {
	if w, is := wr.(*os.File); is {
		TruncateAndSeek(w)
	}
}

// FileOverwrite delete all content in file, and write new content to it
func FileOverwrite(fname string, content string) error {
	return OpenFileForModify(fname, func(fd *os.File) (err error) {
		if err = fd.Truncate(0); err == nil {
			_, err = fd.WriteString(content)
		}
		return
	})
}

// ReadOneLineFrom read first line from file
func ReadOneLine(fname string) (line string, err error) {
	err = FilterFileContent(fname, false, func(l []byte) ([]byte, error) {
		line = string(l)
		return nil, io.EOF
	})
	return
}

// ReadOneLine  read first line from reader, if no content return errNoContent
func ReadOneLineFrom(r io.Reader) (line string, err error) {
	err = FilterContent(r, nil, false, func(l []byte) ([]byte, error) {
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
	filter func([]byte) ([]byte, error)) error {

	return OpenFileForModify(fname, func(fd *os.File) (err error) {
		var wr io.Writer
		if rewrite {
			wr = fd
		}
		return FilterContent(fd, wr, false, filter)
	})
}

// FilterContent readline from reader, after filter, if sync set at the end, and writer is non-null
// write content to writer, otherwise,  every read operation will followed by
// a write operation
func FilterContent(rd io.Reader, wr io.Writer, sync bool, filter func([]byte) ([]byte, error)) (err error) {
	var (
		saveToBuffer = func(line []byte) {}
		flushBuffer  = func() error { return nil }
		line         []byte
		earlyStop    bool
		br           = BufReader(rd)
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
		if sync { // use buffered writer, don't need extra buffer
			w := BufWriter(wr)
			flushBuffer = func() error {
				return w.Flush()
			}
			bw = w
		} else { // if not sync, need buffer to save content for later write
			w := bytes.NewBuffer(make([]byte, 0, FILE_BUFSIZE))
			flushBuffer = func() error {
				TruncateAndSeekWriter(wr)
				_, err := wr.Write(w.Bytes())
				return err
			}
			bw = w
		}
	}
	if filter == nil {
		filter = func(line []byte) ([]byte, error) {
			return line, nil
		}
	}
	for {
		if line, _, err = br.ReadLine(); err == nil {
			if line, err = filter(line); err == nil {
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
