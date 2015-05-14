package goutil

import (
	"go/format"
	"go/parser"
	"go/token"
	"io"
)

// Format source from reader to writer
func Format(fname string, r io.Reader, w io.Writer) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, r, parser.ParseComments)
	if err != nil {
		return err
	}

	return format.Node(w, fset, f)
}
