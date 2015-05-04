package file

import "os"

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
