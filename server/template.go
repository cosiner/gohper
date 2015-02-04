package server

import (
	htmpl "html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type (
	// TemplateEngine is a template engine which support set delimiters, add file suffix name,
	// add template file and dirs, add template functions, compile template, and
	// render template
	TemplateEngine interface {
		SetTemplateDelims(left, right string)
		AddTemplateSuffix(s []string)
		AddTemplates(path ...string) error
		AddTemplateFunc(name string, fn interface{})
		AddTemplateFuncs(funcs map[string]interface{})
		CompileTemplates() error
		RenderTemplate(wr io.Writer, name string, value interface{}) error
	}

	// template implements TemplateEngine interface use standard html/template package
	template struct {
		tmpl *htmpl.Template
	}
)

var (
	// GlobalTmplFuncs is the default template functions
	GlobalTmplFuncs = map[string]interface{}{
	// "I18N": I18N,
	}
	// tmplSuffixes is all template file's suffix
	tmplSuffixes = map[string]bool{"tmpl": true, "html": true}
)

// NewTemplateEngine create a new template engine
func NewTemplateEngine() TemplateEngine {
	return new(template)
}

// isTemplate check whether a file name is recognized template file
func (*template) isTemplate(name string) (is bool) {
	index := strings.LastIndex(name, ".")
	if is = (index >= 0); is {
		is = tmplSuffixes[name[index+1:]]
	}
	return
}

// AddTemplateSuffix add suffix for template
func (*template) AddTemplateSuffix(suffixes []string) {
	for _, suffix := range suffixes {
		if suffix != "" {
			if suffix[0] == '.' {
				suffix = suffix[1:]
			}
			tmplSuffixes[suffix] = true
		}
	}
}

// SetTemplateDelims set default template delimeters
func (*template) SetTemplateDelims(left, right string) {
	strach.setTmplDelims(left, right)
}

// AddTemplates add templates to server, all templates will be
// parsed on server start
func (t *template) AddTemplates(names ...string) (err error) {
	addTmpl := func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && t.isTemplate(path) {
			strach.addTmpl(path)
		}
		return err
	}
	for _, name := range names {
		if err = filepath.Walk(name, addTmpl); err != nil {
			break
		}
	}
	return
}

// CompileTemplates compile all added templates
func (t *template) CompileTemplates() (err error) {
	var tmpl *htmpl.Template
	if tmpls := strach.tmpls(); len(tmpls) != 0 {
		tmpl, err = tmpl.New("tmpl").
			Delims(strach.tmplDelims()).
			Funcs(GlobalTmplFuncs).
			ParseFiles(strach.tmpls()...)
		if err == nil {
			t.tmpl = tmpl
		}
	}
	return
}

// AddTemplateFunc register a function used in templates
func (*template) AddTemplateFunc(name string, fn interface{}) {
	GlobalTmplFuncs[name] = fn
}

// AddTemplateFuncs register some functions used in templates
func (*template) AddTemplateFuncs(funcs map[string]interface{}) {
	for name, fn := range funcs {
		GlobalTmplFuncs[name] = fn
	}
}

// RenderTemplate render a template with given name use given
// value to given writer
func (t *template) RenderTemplate(wr io.Writer, name string, val interface{}) error {
	return t.tmpl.ExecuteTemplate(wr, name, val)
}
