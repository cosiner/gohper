package file

import (
	"io"
	"os"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/io2"
)

const (
	FilePerm = 0644
	DirPerm  = 0755
)

// IsExist check whether or not file/dir exist
func IsExist(fname string) bool {
	_, err := os.Stat(fname)
	return err == nil || !os.IsNotExist(err)
}

// IsFile check whether or not file exist
func IsFile(fname string) bool {
	fi, err := os.Stat(fname)
	return err == nil && !fi.IsDir()
}

// IsDir check whether or not given name is a dir
func IsDir(fname string) bool {
	fi, err := os.Stat(fname)
	return err == nil && fi.IsDir()
}

// IsFileOrNotExist check whether given name is a file or not exist
func IsFileOrNotExist(fname string) bool {
	return !IsDir(fname)
}

// IsDirOrNotExist check whether given is a directory or not exist
func IsDirOrNotExist(dir string) bool {
	return !IsFile(dir)
}

// IsSymlink check whether or not given name is a symlink
func IsSymlink(fname string) bool {
	fi, err := os.Lstat(fname)
	return err == nil && (fi.Mode()&os.ModeSymlink == os.ModeSymlink)
}

// IsModifiedAfter check whether or not file is modified by the function
func IsModifiedAfter(fname string, fn func()) bool {
	fi1, err := os.Stat(fname)
	fn()
	fi2, _ := os.Stat(fname)
	return err == nil && !fi1.ModTime().Equal(fi2.ModTime())
}

// TruncSeek truncate file size to 0 and seek current positon to 0
func TruncSeek(fd *os.File) {
	if fd != nil {
		fd.Truncate(0)
		fd.Seek(0, os.SEEK_SET)
	}
}

// Copy copy src file to dest file
func Copy(dst, src string) error {
	if IsDir(dst) || IsDir(src) {
		return errors.Newf("dest path %s or src path %s is directory", dst, src)
	}
	sf, err := os.Open(src)
	if err == nil {
		var df *os.File
		df, err = os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, FilePerm)
		if err == nil {
			_, err = io.Copy(df, sf)
			df.Close()
		}
		sf.Close()
	}
	return err
}

// Overwrite delete all content in file, and write new content to it
func Overwrite(src string, content string) error {
	return Trunc(src, func(fd *os.File) (err error) {
		_, err = fd.WriteString(content)
		return
	})
}

// FirstLine read first line from file
func FirstLine(src string) (line string, err error) {
	err = Filter(src, false, func(_ int, l []byte) ([]byte, error) {
		line = string(l)
		return nil, io.EOF
	})
	return
}

// Filter filter file content with given filter
// if trunc is true, content after filter will overwrite exist file
// if recieved io.EOF or other error, will stop read next line,
// and not truncate file, io.EOF means an early end.
func Filter(src string, trunc bool,
	filter func(int, []byte) ([]byte, error)) error {
	var openfunc func(string, FileOpFunc) error
	if trunc {
		openfunc = ReadWrite
	} else {
		openfunc = Read
	}
	return openfunc(src, func(fd *os.File) (err error) {
		var w io.Writer
		if trunc {
			w = fd
		}
		return io2.Filter(fd, w, false, filter)
	})
}

// FilterTo filter file content with given filter, then write result
// to dest file
func FilterTo(dst, src string, trunc bool, filter io2.LineFilterFunc) error {
	return Read(src, func(fd *os.File) (err error) {
		return Create(dst, trunc,
			func(dstFd *os.File) error {
				return io2.Filter(fd, dstFd, true, filter)
			})
	})
}

// FileOpFunc accept a file descriptor, return an error or nil
type FileOpFunc func(*os.File) error

// OpenOrCreate open or create file
func OpenOrCreate(fname string) (*os.File, error) {
	return os.OpenFile(fname, os.O_CREATE|os.O_RDWR, FilePerm)
}

// Open openfile use given flag
func Open(fname string, flags int, fn FileOpFunc) error {
	fd, err := os.OpenFile(fname, flags, FilePerm)
	if err == nil {
		if fn != nil {
			err = fn(fd)
		}
		fd.Close()
	}
	return err
}

func Read(fname string, fn FileOpFunc) error {
	return Open(fname, os.O_RDONLY, fn)
}

func ReadWrite(fname string, fn FileOpFunc) error {
	return Open(fname, os.O_RDWR, fn)
}

func Create(fname string, trunc bool, fn FileOpFunc) error {
	return Open(fname, os.O_WRONLY|os.O_CREATE|os.O_EXCL, fn)
}

func Trunc(fname string, fn FileOpFunc) error {
	return Open(fname, os.O_WRONLY|os.O_TRUNC, fn)
}

// WriteFlag return os.O_APPEND if not delete content, else os.O_TRUNC
func WriteFlag(trunc bool) int {
	flag := os.O_APPEND
	if trunc {
		flag = os.O_TRUNC
	}
	return flag
}
