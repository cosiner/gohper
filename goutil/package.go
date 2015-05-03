package goutil

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cosiner/gohper/os2/file"
	"github.com/cosiner/gohper/os2/path2"
)

// PackagePath find package absolute path use env variable GOPATH
func PackagePath(pkgName string) string {
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		for _, path := range strings.FieldsFunc(gopath, func(c rune) bool {
			sep := path2.EnvSeperator()
			return c == sep && sep != path2.UNKNOWN
		}) {
			path = filepath.Join(path, "src", pkgName)
			if file.IsExist(path) && file.IsDir(path) {
				return path
			}
		}
	}
	return ""
}
