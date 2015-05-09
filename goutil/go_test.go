package goutil

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cosiner/gohper/strings2"
	"github.com/cosiner/gohper/testing2"
)

func TestFileType(t *testing.T) {
	testing2.Eq(t, true, IsGoFile("aa.go"))
	testing2.Eq(t, true, IsSrcFile("aa.go"))
	testing2.Eq(t, true, IsTestFile("aa_test.go"))
	testing2.Eq(t, false, IsSrcFile("aa_test.go"))
	testing2.Eq(t, "aa_test.go", SrcTestFile("aa.go"))
}

func TestExport(t *testing.T) {
	tt := testing2.Wrap(t)
	tt.True(IsExported("Name"))
	tt.True(ToSameExported("Abcd", "abc") == "Abc")
	tt.True(ToSameExported("dbcd", "Abc") == "abc")
	tt.True(ToExported("aaa") == "Aaa")
	tt.True(ToExported("") == "")
	tt.True(ToUnexported("Aaa") == "aaa")
	tt.True(ToUnexported("") == "")
}

func TestFormat(t *testing.T) {
	tt := testing2.Wrap(t)
	code := `package main

        func TestExport(t *testing.T) {
            tt := testing2.Wrap(t)
        tt.True(IsExported("Name"))
                tt.True(ToSameExported("Abcd", "abc") == "Abc")
tt.True(ToSameExported("dbcd", "Abc") == "abc");tt.True(ToExported("aaa") == "Aaa")
    tt.True(ToExported("") == "")
    tt.True(ToUnexported("Aaa") == "aaa")
    tt.True(ToUnexported("") == "")
}
    `
	out := bytes.NewBuffer(make([]byte, 0, 1024))
	tt.True(Format("test", strings.NewReader(code), out) == nil)

	tt.Eq(strings2.RemoveSpace(`package main

func TestExport(t *testing.T) {
        tt := testing2.Wrap(t)
        tt.True(IsExported("Name"))
        tt.True(ToSameExported("Abcd", "abc") == "Abc")
        tt.True(ToSameExported("dbcd", "Abc") == "abc")
        tt.True(ToExported("aaa") == "Aaa")
        tt.True(ToExported("") == "")
        tt.True(ToUnexported("Aaa") == "aaa")
        tt.True(ToUnexported("") == "")
}
`), strings2.RemoveSpace(out.String()))
}

func TestPath(t *testing.T) {
	tt := testing2.Wrap(t)
	out := bytes.NewBuffer(make([]byte, 0, 1024))
	WriteImportpath(out, "github.com", "cosiner", "gohper")
	tt.Eq(`"github.com/cosiner/gohper"`+"\n", out.String())

	PackagePath("bufio") // difficult to test
}
