package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/strings2"
)

const (
	// END cause whole parsing end
	END = errors.Err("parsing end")

	// TYPE_END cause single type's parsing end
	TYPE_END = errors.Err("type parsing end")
)

type Attrs struct {
	// Package name
	Package string

	// helpful for access single type, if don't know, don't use it
	Accessed bool

	// Common
	TypeName string

	// Struct
	Field, Tag string

	// Const
	Name, Value string

	// Func, share Name attribute whith Const
	// Name string
	PtrRecv bool // whether a method's reciever is pointer, only valid for method
}

type Callback struct {
	// if Struct is not nil, StrucField should also be,
	// otherwize this type will be skipped
	Struct      func(*Attrs) error
	StructField func(*Attrs) error

	Const func(*Attrs) error

	Func func(*Attrs) error
}

func ParseFile(fname string, call Callback) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err == nil {
		err = Parse(f, call)
	}
	return err
}

func Parse(file *ast.File, call Callback) (err error) {
	attrs := &Attrs{
		Package: file.Name.Name,
	}

	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.TYPE {
				if call.Struct == nil || call.StructField == nil {
					continue
				}

				for i, l := 0, len(decl.Specs); i < l; i++ {
					err = call.callStruct(decl.Specs[i].(*ast.TypeSpec), attrs)
					if err != nil {
						goto END
					}
				}
			} else if decl.Tok == token.CONST && call.Const != nil {
				if err = call.callConsts(decl, attrs); err != nil {
					goto END
				}
			}
		case *ast.FuncDecl:
			if call.Func != nil {
				if err = call.callFunc(decl, attrs); err != nil {
					goto END
				}
			}
		}
	}

END:
	if err == END {
		err = nil
	}
	return err
}

func (call Callback) callStruct(spec *ast.TypeSpec, attrs *Attrs) error {
	st, is := spec.Type.(*ast.StructType)
	if !is {
		return nil
	}

	attrs.TypeName = spec.Name.Name
	err := call.Struct(attrs)
	if err != nil {
		goto END
	}

	for _, f := range st.Fields.List {
		for _, n := range f.Names {
			attrs.Field = n.Name
			attrs.Tag = ""
			if f.Tag != nil {
				attrs.Tag, _ = strings2.TrimQuote(f.Tag.Value)
			}
			if err = call.StructField(attrs); err != nil {
				goto END
			}
		}
	}

END:
	if err == TYPE_END {
		err = nil
	}
	return err
}

func (call Callback) callConsts(decl *ast.GenDecl, attrs *Attrs) error {
	attrs.TypeName = ""

	for i, l := 0, len(decl.Specs); i < l; i++ {
		spec := decl.Specs[i].(*ast.ValueSpec)

		if spec.Type != nil {
			attrs.TypeName = fmt.Sprint(spec.Type)
		} else if spec.Values != nil {
			attrs.TypeName = "" // iota is break out
		}

		vlen := len(spec.Values)
		for i, name := range spec.Names {
			attrs.Name = name.Name
			attrs.Value = ""
			if i < vlen {
				attrs.Value = fmt.Sprint(spec.Values[i])
			}
			if err := call.Const(attrs); err == TYPE_END {
				return nil
			} else if err != nil {
				return err
			}
		}
	}
	return nil
}

func (call Callback) callFunc(decl *ast.FuncDecl, attrs *Attrs) error {
	attrs.TypeName = ""
	attrs.PtrRecv = false
	attrs.Name = decl.Name.Name

	if decl.Recv != nil {
		switch recv := decl.Recv.List[0].Type.(type) {
		case *ast.Ident:
			attrs.TypeName = recv.Name
		case *ast.StarExpr:
			attrs.TypeName = fmt.Sprint(recv.X)
			attrs.PtrRecv = true
		}
	}

	err := call.Func(attrs)
	if err == TYPE_END {
		err = nil
	}
	return err
}
