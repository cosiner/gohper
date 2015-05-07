package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/cosiner/gohper/errors"
)

const EOF = errors.Err("parsing end")

type Attrs struct {
	Package string

	// helpful for access single type, if don't know, don't use it
	Accessed bool

	// Common
	TypeName string

	// Struct
	Field, Tag string

	// Const
	Name, Value string

	// Func, share Const's Name attribute
	// Name string
}

type Callback struct {
	Struct func(*Attrs) error
	Const  func(*Attrs) error
	Func   func(*Attrs) error
}

func ParseFile(fname string, call Callback) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err == nil {
		err = Parse(f, call)
	}
	return err
}

func Parse(file *ast.File, call Callback) error {
	attrs := &Attrs{
		Package: file.Name.Name,
	}

	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.TYPE {
				if call.Struct == nil {
					continue
				}

				for i, l := 0, len(decl.Specs); i < l; i++ {
					spec := decl.Specs[i].(*ast.TypeSpec)

					if st, is := spec.Type.(*ast.StructType); is {
						attrs.TypeName = spec.Name.Name
						if err := call.callStruct(st, attrs); err != nil {
							return nonEOF(err)
						}
					}
				}
			} else if decl.Tok == token.CONST && call.Const != nil {
				attrs.TypeName = ""
				for i, l := 0, len(decl.Specs); i < l; i++ {
					spec := decl.Specs[i].(*ast.ValueSpec)
					if spec.Type != nil {
						attrs.TypeName = fmt.Sprint(spec.Type)
					}
					if err := call.callConsts(spec, attrs); err != nil {
						return nonEOF(err)
					}
				}
			}
		case *ast.FuncDecl:
			if call.Func != nil {
				attrs.TypeName = ""
				if decl.Recv != nil {
					attrs.TypeName = fmt.Sprint(decl.Recv.List[0].Type)
				}
				if err := call.callFunc(decl, attrs); err != nil {
					return nonEOF(err)
				}
			}
		}
	}
	return nil
}

func (call Callback) callStruct(t *ast.StructType, attrs *Attrs) error {
	for _, f := range t.Fields.List {
		for _, n := range f.Names {
			attrs.Field = n.Name
			attrs.Tag = ""
			if f.Tag != nil {
				attrs.Tag = f.Tag.Value
			}
			if err := call.Struct(attrs); err != nil {
				return err
			}
		}
	}
	return nil
}

func (call Callback) callConsts(spec *ast.ValueSpec, attrs *Attrs) error {
	vl := len(spec.Values)
	for i, name := range spec.Names {
		attrs.Name = name.Name
		attrs.Value = ""
		if i < vl {
			attrs.Value = fmt.Sprint(spec.Values[i])
		}
		if err := call.Const(attrs); err != nil {
			return err
		}
	}
	return nil
}

func (call Callback) callFunc(decl *ast.FuncDecl, attrs *Attrs) error {
	attrs.Name = decl.Name.Name
	return call.Func(attrs)
}

func nonEOF(err error) error {
	if err == EOF {
		err = nil
	}
	return err
}
