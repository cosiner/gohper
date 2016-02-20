package file

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/io2"
)

const ErrDestIsFile = errors.Err("destination is a file")

// FileOpFunc accept a file descriptor, return an error or nil
type FileOpFunc func(*os.File) error

// Open file use given flag
func Open(fname string, flags int, fn FileOpFunc) error {
	fd, err := os.OpenFile(fname, flags, FilePerm)
	if err != nil {
		return err
	}

	if fn != nil {
		err = fn(fd)
	}
	if e := fd.Close(); e != nil && err == nil {
		err = e
	}

	return io2.NonEOF(err)
}

// WriteFlag return os.O_APPEND if not delete content, else os.O_TRUNC
func WriteFlag(trunc bool) int {
	if trunc {
		return os.O_TRUNC
	}

	return os.O_APPEND
}

// FirstLine read first line from file
func FirstLine(src string) (line string, err error) {
	err = Filter(src, func(_ int, l []byte) (error) {
		line = string(l)

		return io.EOF
	})

	return
}

// Filter file content with given filter, file is in ReadOnly mode
func Filter(src string, filter func(int, []byte) error) error {
	return Read(src, func(fd *os.File) (err error) {
		return io2.FilterRead(fd, filter)
	})
}

// FilterTo filter file content with given filter, then write result
// to dest file
func FilterTo(dst, src string, trunc bool, filter io2.LineFilterFunc) error {
	return Read(src, func(sfd *os.File) (err error) {
		return OpenOrCreate(dst, trunc, func(dfd *os.File) error {
			return io2.Filter(sfd, dfd, true, filter)
		})
	})
}

// Copy src file to dest file
func Copy(dst, src string) error {
	return FilterTo(dst, src, true, nil)
}

// CopyDir copy directory from source to destination
func CopyDir(dst, src string) error {
	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(src)

	for i := 0; i < len(files) && err == nil; i++ {
		file := files[i].Name()
		df := filepath.Join(dst, file)
		sf := filepath.Join(src, file)
		if IsFile(sf) {
			err = Copy(df, sf)
		} else {
			err = CopyDir(df, sf)
		}
	}

	return err
}

// Overwrite delete all content in file, and write new content to it
func Overwrite(src string, content string) error {
	return Trunc(src, func(fd *os.File) error {
		_, err := fd.WriteString(content)

		return err
	})
}

// for Read

func Read(fname string, fn FileOpFunc) error {
	return Open(fname, os.O_RDONLY, fn)
}

func ReadWrite(fname string, fn FileOpFunc) error {
	return Open(fname, os.O_RDWR, fn)
}

//  for Write

func Write(fname string, fn FileOpFunc) error {
	return Open(fname, os.O_WRONLY, fn)
}

func Trunc(fname string, fn FileOpFunc) error {
	return Open(fname, os.O_WRONLY|os.O_TRUNC, fn)
}

func Create(fname string, fn FileOpFunc) error {
	return Open(fname, os.O_WRONLY|os.O_CREATE|os.O_EXCL, fn)
}

func Append(fname string, fn FileOpFunc) error {
	return Open(fname, os.O_WRONLY|os.O_APPEND, fn)
}

func OpenOrCreate(fname string, trunc bool, fn FileOpFunc) error {
	return Open(fname, os.O_CREATE|os.O_WRONLY|WriteFlag(trunc), fn)
}
