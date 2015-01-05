package gotool

import (
	"go/ast"
	"mlib/util/types"
)

// IsExported return whether or not field is exported
// it's just a wrapper of ast.IsExported
func IsExported(name string) bool {
	return ast.IsExported(name)
}

// SameExportedCase return the same exported case with example string of a string
func SameExportedCase(example, str string) string {
	if IsExported(example) {
		return ExportedCase(str)
	} else {
		return UnexportedCase(str)
	}
}

// ExportedCase return the exported case of a string
func ExportedCase(str string) string {
	if str == "" {
		return ""
	}
	return string(types.UpperCase(str[0])) + str[1:]
}

// ExportedCase return the unexported case of a string
func UnexportedCase(str string) string {
	if str == "" {
		return ""
	}
	return string(types.LowerCase(str[0])) + str[1:]
}
