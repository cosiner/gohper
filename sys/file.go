package sys

import (
	"bytes"
	"io"
	"os"

	. "github.com/cosiner/golib/errors"
)

const (
	// FILE_BBUFSIZe is the buffer size of slice to store file content
	FILE_BUFSIZE     = 4096
	FILE_PERM        = 0644
	DIR_PERM         = 0755
	CR           int = os.O_CREATE
	AP           int = os.O_APPEND
	TC           int = os.O_TRUNC
	EX           int = os.O_EXCL
	WR           int = os.O_WRONLY
	RD           int = os.O_RDONLY
	RW           int = os.O_RDWR
	TW           int = TC | WR
	CW           int = CR | WR
	CTW          int = CR | WR | TC
	CEW          int = CR | EX | WR
)

// IsExist check whether or not file/dir exist
func IsExist(fname string) bool {
	_, err := os.Stat(ExpandHome(fname))
	return err == nil
}

// IsFile check whether or not file exist
func IsFile(fname string) bool {
	fi, err := os.Stat(ExpandHome(fname))
	return err == nil && !fi.IsDir()
}

// IsDir check whether or not given name is a dir
func IsDir(fname string) bool {
	fi, err := os.Stat(ExpandHome(fname))
	return err == nil && fi.IsDir()
}

// IsFileOrNotExist check whether given name is a file or not exist
func IsFileOrNotExist(fname string) bool {
	return !IsDir(fname)
}

// IsDirOrNotExist check whether given is a directory or not exist
func IsDirOrNotExist(dirName string) bool {
	return !IsFile(dirName)
}

// IsSymlink check whether or not given name is a symlink
func IsSymlink(fname string) bool {
	fi, err := os.Lstat(ExpandHome(fname))
	return err == nil && (fi.Mode()&os.ModeSymlink == os.ModeSymlink)
}

// IsFileModifiedAfter check whether or not file is modified by the function
func IsFileModifiedAfter(fname string, fn func()) bool {
	fname = ExpandHome(fname)
	fi1, err := os.Stat(fname)
	fn()
	fi2, _ := os.Stat(fname)
	return err == nil && !fi1.ModTime().Equal(fi2.ModTime())
}

// FileOpFunc accept a file descriptor, return an error or nil
type FileOpFunc func(*os.File) error

// OpenOrCreate open or create file
func OpenOrCreate(fname string) (*os.File, error) {
	return os.OpenFile(ExpandHome(fname), CW, FILE_PERM)
}

// OpenAndTruncFor open file, clear all content
func OpenAndTruncFor(fname string, fn FileOpFunc) error {
	return OpenFileFor(fname, TW, fn)
}

// OpenOrCreateFor openfile for write
func OpenOrCreateFor(fname string, overwrite bool, fn FileOpFunc) error {
	return OpenFileFor(fname, CW|WriteFlag(overwrite), fn)
}

// OpenOrCreateTruncFor open or create file for write, clear all content
func OpenOrCreateTruncFor(fname string, fn FileOpFunc) error {
	return OpenFileFor(fname, CTW, fn)
}

// OpenForWrite open file for read
func OpenForWrite(fname string, fn FileOpFunc) error {
	return OpenFileFor(fname, WR, fn)
}

// OpenOrCreateForWrite open or create file for write
func OpenOrCreateForWrite(fname string, fn FileOpFunc) error {
	return OpenFileFor(fname, CW, fn)
}

// OpenForRead open file for read
func OpenForRead(fname string, fn FileOpFunc) error {
	return OpenFileFor(fname, RD, fn)
}

// OpenFileForRW open file for write and read
func OpenForRW(fname string, fn FileOpFunc) error {
	return OpenFileFor(fname, RW, fn)
}

// CreateFor create new file for write, if already exist, return error
func CreateFor(fname string, fn FileOpFunc) error {
	return OpenFileFor(fname, CEW, fn)
}

// OpenFileFor openfile use given flag
func OpenFileFor(fname string, flags int, fn FileOpFunc) error {
	fd, err := os.OpenFile(ExpandHome(fname), flags, FILE_PERM)
	if err == nil {
		if fn != nil {
			err = fn(fd)
		}
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

// CopyFile copy src file to dest file
func CopyFile(dst, src string) error {
	if IsDir(dst) || IsDir(src) {
		return Errorf("dest path %s or src path is directory", dst, src)
	}
	return OpenOrCreateTruncFor(dst, func(dstFd *os.File) error {
		return OpenForRead(src, func(srcFd *os.File) error {
			_, err := io.Copy(dstFd, srcFd)
			return err
		})
	})
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
	return OpenAndTruncFor(fname, func(fd *os.File) (err error) {
		_, err = fd.WriteString(content)
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

// ReadOneLine read first line from reader, if no content return errNoContent
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
	filter func([]byte) ([]byte, error)) error {
	return OpenForRead(fname, func(fd *os.File) (err error) {
		return OpenOrCreateFor(dstfile, delContent,
			func(dstFd *os.File) error {
				return FilterContent(fd, dstFd, true, filter)
			})
	})
}

// FilterContent readline from reader, after filter, if parallel set at the end, and writer is non-null
// write content to writer, otherwise,  every read operation will followed by
// a write operation
func FilterContent(rd io.Reader, wr io.Writer,
	parallel bool, filter func([]byte) ([]byte, error)) (err error) {
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
		if parallel { // use buffered writer, don't need extra buffer
			w := BufWriter(wr)
			flushBuffer = func() error {
				return w.Flush()
			}
			bw = w
		} else { // if not parallel, need buffer to save content for later write
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

// CloseFd close io.Closer when it's effective
func CloseFd(fd io.Closer) {
	if fd != nil {
		fd.Close()
	}
}

// WriteFlag return APPEND if not delete content, else TRUNCATE
func WriteFlag(delContent bool) int {
	flag := os.O_APPEND
	if delContent {
		flag = os.O_TRUNC
	}
	return flag
}
