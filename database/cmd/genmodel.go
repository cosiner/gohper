package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cosiner/golib/sys"

	. "github.com/cosiner/golib/errors"
	"github.com/cosiner/golib/goutil"
	"github.com/cosiner/golib/types"

	"flag"
)

var (
	infile       string
	outfile      string
	models       string
	tmpl         string
	useTmpl      string
	printDefTmpl bool
)

func cliArgs() {
	flag.StringVar(&infile, "i", "", "input file")
	flag.StringVar(&outfile, "o", "", "output file")
	flag.StringVar(&models, "m", "", "models to parse, seperate by comma")
	flag.StringVar(&tmpl, "t", "", "template file")
	flag.BoolVar(&printDefTmpl, "deftmpl", false, "print default path of template file")
	flag.StringVar(&useTmpl, "ut", "", "use this file as default template file")
	flag.Parse()
}

//go:generate cp model.tmpl ~/.config/go/model.tmpl
func main() {
	cliArgs()
	defTmplPath := filepath.Join(sys.HomeDir(), ".config", "go", "model.tmpl")
	if printDefTmpl {
		ExitWith(defTmplPath)
	}

	if useTmpl != "" {
		if sys.IsExist(defTmplPath) {
			ExitErrorln("default template file already exist or is a directory")
		}
		OnErrDo(sys.CopyFile(defTmplPath, useTmpl), ExitErrln)
		return
	}

	if infile == "" {
		ExitErrorln("No input file specified.")
	}
	models := types.TrimSplit(models, ",")
	tree, err := parser.ParseFile(token.NewFileSet(), infile, nil, 0)
	OnErrDo(err, ExitErrln)
	mv := new(modelVisitor)
	mv.addModels(models)
	mv.walk(tree)
	if len(mv.models) == 0 {
		return
	}
	if outfile == "" {
		outfile = infile
	}
	OnErrExit(sys.OpenOrCreateFor(outfile, false, func(outfd *os.File) error {
		modelFields := buildModelFields(mv.models)
		if tmpl == "" {
			tmpl = defTmplPath
		}
		var t *template.Template
		if t, err = template.ParseFiles(tmpl); err == nil {
			err = t.Execute(outfd, modelFields)
		}
		return nil
	}))
}

func excludeField(field *ast.Field) bool {
	if t, is := field.Type.(*ast.Ident); is {
		if t.Name == "ColumnParser" {
			return true
		}
	}
	return false
}

type StructName struct {
	Name           string // struct's normal name
	Self           string
	UnexportedName string
	LowerName      string // lower case name
	UpperName      string // upper case name
}

type FieldName struct {
	Name       string // field's normal name
	ColumnName string // field's column name, in snake_case
	ConstName  string // field's const name is in STRUCTNAME_FIELDNAME case
}

func NewFieldName(model *StructName, field string) *FieldName {
	return &FieldName{Name: field,
		ColumnName: strings.ToLower(types.SnakeString(field)),
		ConstName:  model.UpperName + "_" + strings.ToUpper(field)}
}

func NewStructName(name string) *StructName {
	return &StructName{Name: name,
		Self:           types.AbridgeStringToLower(name),
		UnexportedName: goutil.UnexportedCase(name),
		LowerName:      strings.ToLower(name),
		UpperName:      strings.ToUpper(name)}
}

// buildModelFields build model map from parse result
func buildModelFields(models map[string][]string) map[*StructName][]*FieldName {
	names := make(map[*StructName][]*FieldName, len(models))
	for model, fields := range models {
		modelStruct := NewStructName(model)
		for _, name := range fields {
			names[modelStruct] = append(names[modelStruct], NewFieldName(modelStruct, name))
		}
	}
	return names
}

type modelVisitor struct {
	models map[string][]string
}

// initContainer init result container
func (mv *modelVisitor) initContainer() {
	if mv.models == nil {
		mv.models = make(map[string][]string, 10)
	}
}

// addModels add models that need parse
func (mv *modelVisitor) addModels(models []string) {
	if len(models) == 0 {
		return
	}
	mv.initContainer()
	for _, m := range models {
		if m != "" {
			mv.models[m] = nil
		}
	}
}

// add add an model and it's field to parse result
func (mv *modelVisitor) add(model, field string) {
	mv.initContainer()
	mv.models[model] = append(mv.models[model], field)
}

// needParse check whether a model should be parsed
// unexporeted model don't parse
// if visitor's model list is not empty, only parse model exist in list
// otherwise parse all
func (mv *modelVisitor) needParse(model string) bool {
	if !goutil.IsExported(model) {
		return false
	}
	if mv.models != nil && len(mv.models) > 0 {
		if _, has := mv.models[model]; !has {
			return false
		}
	}
	return true
}

// walk parse ast tree to find exported struct and it's fields
func (mv *modelVisitor) walk(tree *ast.File) {
	for _, decl := range tree.Decls { // Top Declare
		if decl, is := decl.(*ast.GenDecl); is { // General Declare
			if decl.Tok == token.TYPE { // Type Keyword
				for _, spec := range decl.Specs {
					spec, _ := spec.(*ast.TypeSpec)
					if t, is := spec.Type.(*ast.StructType); is { // type struct
						model := spec.Name.Name // model name
						if !mv.needParse(model) {
							continue
						}
						for _, f := range t.Fields.List { // model field
							if excludeField(f) {
								continue
							}
							for _, ident := range f.Names {
								name := ident.Name
								if goutil.IsExported(name) {
									mv.add(model, name)
								}
							}
						}
					}
				}
			}
		}
	}
}
