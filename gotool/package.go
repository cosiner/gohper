package gotool

import (
	"mlib/util/sys"
	"os"
	"path/filepath"
	"strings"
)

// FindPackage find package absolute path use env variable GOPATH
func FindPackage(pkgName string) string {
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		for _, path := range strings.FieldsFunc(gopath, func(c rune) bool {
			sep := sys.EnvPathSeperator()
			return c == sep && sep != sys.UNKNOWN_PATH_SEPRATOR
		}) {
			path = filepath.Join(path, "src", pkgName)
			if sys.IsExist(path) && sys.IsDir(path) {
				return path
			}
		}
	}
	return ""
}
