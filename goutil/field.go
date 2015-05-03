package goutil

import (
	"go/ast"

	"github.com/cosiner/gohper/unibyte"
)

// IsExported return whether or not field is exported
// it's just a wrapper of ast.IsExported
func IsExported(name string) bool {
	return ast.IsExported(name)
}

// ToSameExported return the same exported case with example string of a string
func ToSameExported(example, str string) (res string) {
	if IsExported(example) {
		res = ToExported(str)
	} else {
		res = ToUnexported(str)
	}
	return
}

// ToExported return the exported case of a string
func ToExported(str string) string {
	if str == "" {
		return ""
	}
	return unibyte.ToUpperString(str[0]) + str[1:]
}

// ToUnexported return the unexported case of a string
func ToUnexported(str string) string {
	if str == "" {
		return ""
	}
	return unibyte.ToLowerString(str[0]) + str[1:]
}
