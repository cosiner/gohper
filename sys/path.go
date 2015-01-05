package sys

import (
	. "github.com/cosiner/golib/errors"
	"github.com/cosiner/golib/types"
	"os"
	"os/user"
	"path/filepath"
)

// HomeDir return current user's Home dir
func HomeDir() string {
	u, _ := user.Current()
	return u.HomeDir
}

// ExpandHome expand ~ to user's home dir
func ExpandHome(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}
	u, _ := user.Current()
	return u.HomeDir + path[1:]
}

// ExpandAbs expand path to absolute path
func ExpandAbs(path string) string {
	path, _ = filepath.Abs(ExpandHome(path))
	return path
}

// ProgramDir return dir of program use os.Args[0]
func ProgramDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// LastDir return last dir of path,
// if path is dir, return itself
// else return path's contain dir name
func LastDir(path string) (string, error) {
	var dir string
	absPath, err := filepath.Abs(path)
	if err == nil {
		info, err := os.Stat(absPath)
		if err == nil {
			if info.IsDir() {
				_, dir = filepath.Split(absPath)
			} else {
				dir = filepath.Dir(absPath)
				_, dir = filepath.Split(dir)
			}
		}
	}
	return dir, err
}

// IsRootPath check wether or not path is root of filesystem
func IsRootPath(path string) bool {
	if l := len(path); l > 0 {
		switch OperateSystem() {
		case WINDOWS:
			return types.IsLetter(path[0]) && path[1:] == ":\\"
		case LINUX, DARWIN, FREEBSD, SOLARIS, ANDROID:
			return l == 1 && path[0] == '/'
		}
	}
	return false
}

// MkdirWithParent create dirs with parent dir
func MkdirWithParent(path string) error {
	if IsExist(path) {
		return nil
	}
	var err error
	dirs := make([]string, 0, 4)
	path = ExpandAbs(path)
	for {
		if os.Mkdir(path, DIR_PERM) != nil {
			dirs = append(dirs, path)
			path = filepath.Dir(path)
			if IsRootPath(path) {
				return Errorf("Can't make dir %s with parent, all subdirectory can't be created",
					path)
			}
		} else {
			break
		}
	}
	for i := len(dirs) - 1; i >= 0; i-- {
		err = os.Mkdir(dirs[i], DIR_PERM)
	}
	return err
}
