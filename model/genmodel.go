package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/cosiner/golib/sys"

	. "github.com/cosiner/golib/errors"
	"github.com/cosiner/golib/goutil"
	"github.com/cosiner/golib/types"

	"flag"
)

var (
	infile  string
	models  string
	outfile string
)

func cliArgs() {
	flag.StringVar(&infile, "i", "", "input file")
	flag.StringVar(&outfile, "o", "", "output file")
	flag.StringVar(&models, "m", "", "models to parse, seperate by comma")
	flag.Parse()
}

func main() {
	cliArgs()

	if infile == "" {
		ErrorlnExit("No input file specified.")
	}
	models := types.TrimSplit(models, ",")
	tree, err := parser.ParseFile(token.NewFileSet(), infile, nil, 0)
	if err != nil {
		ErrorlnExit(err)
	}
	mv := new(modelVisitor)
	mv.addModels(models)
	mv.walk(tree)
	if outfile == "" {
		outfile = infile
	}
	outfd, err := os.OpenFile(outfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		ErrorlnExit(err)
	}
	out := sys.WrapWriter(outfd)
	modelFields := buildModelFields(mv.models)
	for model, fields := range modelFields {
		output(out, model, fields)
	}
	out.Flush()
	outfd.Close()
}

func output(bw sys.BufVWriter, model string, fields []*field) {
	lowerModel := strings.ToLower(model)
	upperModel := strings.ToUpper(model)
	bw.WriteString("var (\n")
	bw.WriteLString("\t", lowerModel, "Columns = [...]string{")
	for _, f := range fields {
		bw.WriteLString("\"", f.columnName, "\", ")
	}
	bw.WriteString("}\n")
	bw.WriteLString("\t", lowerModel, "Fields = [...]model.Field{")
	for _, f := range fields {
		bw.WriteLString(f.constName, ", ")
	}
	bw.WriteString("}\n")
	bw.WriteString(")\n")
	bw.WriteString("const (\n")
	for _, f := range fields {
		bw.WriteLString("\t", f.constName, " = iota\n")
	}
	bw.WriteLString("\t", upperModel, "_TABLE = \"", lowerModel, "\"")
	bw.WriteString("\n)\n")

	firstStr := lowerModel[:1]
	bw.WriteString(fmt.Sprintf(`func (%s *%s) Table() string {
    return %s_TABLE
}
`, firstStr, model, upperModel))
	bw.WriteString(fmt.Sprintf(`func (%s *%s) Fields() []Field {
    return %sFields
}
`, firstStr, model, lowerModel))

	bw.WriteString(fmt.Sprintf(`func (%s *%s) Columns() []string {
    return %sColumns
}
`, firstStr, model, lowerModel))

	bw.WriteString(fmt.Sprintf(`func (%s *%s) ColumnName(field Field) string {
    %s.MustEffectiveField(field)
    return %sColumns[field.Num()]
}
`, firstStr, model, firstStr, lowerModel))

	bw.WriteString(fmt.Sprintf(`func (%s *%s) Init() *%s {
    cp := NewColumnParser()
    cp.Bind(%s)
    %s.ColumnParser = cp
    return %s
}
`, firstStr, model, model, lowerModel, firstStr, firstStr))

	bw.WriteString(fmt.Sprintf(`func (%s *%s) FieldVal(field Field) (val interface{}) {
    %s.MustEffectiveField(field)
    switch field {
`, firstStr, model, firstStr))
	for _, f := range fields {
		bw.WriteString(fmt.Sprintf("\t\tcase %s:val = %s.%s\n", f.constName, firstStr, f.name))
	}
	bw.WriteString("\t}\n\treturn\n}\n")
}

func excludeField(field *ast.Field) bool {
	if t, is := field.Type.(*ast.Ident); is {
		if t.Name == "ColumnParser" {
			return true
		}
	}
	return false
}

type modelVisitor struct {
	models map[string][]string
}

type field struct {
	name       string
	columnName string
	constName  string
}

func buildModelFields(models map[string][]string) map[string][]*field {
	mfs := make(map[string][]*field, len(models))
	for model, fields := range models {
		for _, name := range fields {
			mfs[model] = append(mfs[model],
				&field{name,
					types.SnakeString(name),
					strings.ToUpper(model + "_" + name)})
		}
	}
	return mfs
}

func (mv *modelVisitor) initContainer() {
	if mv.models == nil {
		mv.models = make(map[string][]string, 10)
	}
}

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

func (mv *modelVisitor) add(model, field string) {
	mv.initContainer()
	mv.models[model] = append(mv.models[model], field)
}

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
