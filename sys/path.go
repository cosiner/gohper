package sys

import (
	"os"
	"os/user"
	"path/filepath"

	. "github.com/cosiner/golib/errors"
	"github.com/cosiner/golib/types"
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

// IsRelativePath check whether a path is relative
// In these condition: path is empty, start with '[.~][/\]', '/', "[a-z]:\"
func IsRelativePath(path string) bool {
	return !(types.StartWith(path, "./") ||
		types.StartWith(path, ".\\") ||
		types.StartWith(path, "~/") ||
		types.StartWith(path, "~\\") ||
		types.StartWith(path, "/") ||
		IsWindowsRootpath(path))
}

// IsWindowsRootPath check whether a path is windows absolute path with disk letter
func IsWindowsRootpath(path string) bool {
	if path == "" {
		return false
	}
	return types.IsLetter(path[0]) && types.StartWith(path[1:], ":\\")
}

// IsRootPath check wether or not path is root of filesystem
func IsRootPath(path string) bool {
	if l := len(path); l > 0 {
		switch OperateSystem() {
		case WINDOWS:
			return IsWindowsRootpath(path)
		case LINUX, DARWIN, FREEBSD, SOLARIS, ANDROID:
			return l == 1 && path[0] == '/'
		}
	}
	return false
}

// MkdirWithParent create dirs with parent dir
func MkdirWithParent(path string) error {
	if IsDir(path) {
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
