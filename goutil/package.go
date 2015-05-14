package goutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosiner/gohper/os2/file"
	"github.com/cosiner/gohper/os2/path2"
)

// PackagePath find package absolute path use env variable GOPATH
func PackagePath(pkgName string) string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return ""
	}

	fn := func(c rune) bool {
		sep := path2.EnvSeperator()

		return c == sep && sep != path2.UNKNOWN
	}
	for _, path := range strings.FieldsFunc(gopath, fn) {
		path = filepath.Join(path, "src", pkgName)

		if file.IsExist(path) && file.IsDir(path) {
			return path
		}
	}

	return ""
}

// WriteImportPath write path to writer, automaticlly join with '/',
// quote it and start a new line
func WriteImportpath(w io.Writer, path ...string) (int, error) {
	return fmt.Fprintf(w, `"%s"`+"\n", strings.Join(path, "/"))
}
