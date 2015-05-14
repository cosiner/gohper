package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"

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
	S struct {
		Field, Type string // if type is empty, means anonymous field
		Tag         reflect.StructTag
	}

	// Const
	C struct {
		Name, Value string
	}

	// Interface
	I struct {
		Method string
	}

	// Func
	F struct {
		Name    string
		PtrRecv bool // whether a method's reciever is pointer, only valid for method
	}
}

type Parser struct {
	Const     func(*Attrs) error
	Interface func(*Attrs) error
	Struct    func(*Attrs) error
	Func      func(*Attrs) error
}

func (p Parser) ParseFile(fname string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	return p.Parse(f)
}

func (p Parser) Parse(file *ast.File) (err error) {
	attrs := &Attrs{
		Package: file.Name.Name,
	}

	defer func(err *error) {
		if *err == END {
			*err = nil
		}
	}(&err)

	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.TYPE {
				if p.Struct == nil || p.Interface == nil {
					continue
				}

				for _, spec := range decl.Specs {
					if err = p.parseType(spec.(*ast.TypeSpec), attrs); err != nil {
						return
					}
				}
			} else if decl.Tok == token.CONST && p.Const != nil {
				if err = p.parseConsts(decl, attrs); err != nil {
					return
				}
			}

		case *ast.FuncDecl:
			if p.Func != nil {
				if err = p.parseFunc(decl, attrs); err != nil {
					return
				}
			}
		}
	}

	return
}

func (p Parser) parseType(spec *ast.TypeSpec, attrs *Attrs) error {
	attrs.TypeName = ""
	switch typ := spec.Type.(type) {
	case *ast.StructType:
		if p.Struct != nil {
			attrs.TypeName = spec.Name.Name
			return p.parseStruct(typ, attrs)
		}

	case *ast.InterfaceType:
		if p.Interface != nil {
			attrs.TypeName = spec.Name.Name
			return p.parseInterface(typ, attrs)
		}
	}

	return nil
}

func (p Parser) parseStruct(spec *ast.StructType, attrs *Attrs) (err error) {
	defer nonTypeEnd(&err)

	for _, f := range spec.Fields.List {
		attrs.S.Type = ""
		if f.Type != nil {
			attrs.S.Type = fmt.Sprint(f.Type)
		}

		for _, n := range f.Names {
			attrs.S.Field = n.Name
			attrs.S.Tag = ""

			if f.Tag != nil {
				tag, _ := strings2.TrimQuote(f.Tag.Value)
				attrs.S.Tag = reflect.StructTag(tag)
			}

			if err = p.Struct(attrs); err != nil {
				return
			}
		}
	}

	return
}

func (p Parser) parseInterface(spec *ast.InterfaceType, attrs *Attrs) (err error) {
	defer nonTypeEnd(&err)

	for _, m := range spec.Methods.List {
		for _, n := range m.Names {
			attrs.I.Method = n.Name
			if err = p.Interface(attrs); err != nil {
				return
			}
		}
	}

	return
}

func (p Parser) parseConsts(decl *ast.GenDecl, attrs *Attrs) (err error) {
	defer nonTypeEnd(&err)

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
			attrs.C.Name = name.Name
			attrs.C.Value = ""
			if i < vlen {
				attrs.C.Value = fmt.Sprint(spec.Values[i])
			}

			if err = p.Const(attrs); err != nil {
				return
			}
		}
	}

	return
}

func (p Parser) parseFunc(decl *ast.FuncDecl, attrs *Attrs) (err error) {
	attrs.TypeName = ""
	attrs.F.PtrRecv = false
	attrs.F.Name = decl.Name.Name

	if decl.Recv != nil {
		switch recv := decl.Recv.List[0].Type.(type) {
		case *ast.Ident:
			attrs.TypeName = recv.Name
		case *ast.StarExpr:
			attrs.TypeName = fmt.Sprint(recv.X)
			attrs.F.PtrRecv = true
		}
	}

	err = p.Func(attrs)
	nonTypeEnd(&err)

	return
}

func nonTypeEnd(err *error) {
	if *err == TYPE_END {
		*err = nil
	}
}
