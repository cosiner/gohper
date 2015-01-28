package server

import (
	"os"
	"path/filepath"

	. "github.com/cosiner/golib/errors"

	"github.com/cosiner/gomodule/config"
)

// translation represent locales data
type translation struct {
	defLocale map[string]string
	locales   map[string]map[string]string
}

// _tr is the global i18n translator
var _tr *translation = new(translation)

// I18N translate a message to locale-specified string
func I18N(locale, message string) string {
	return _tr.translate(locale, message)
}

// Load load locale data from file or dir, use base name of file as locale name
func (tr *translation) load(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			if name := info.Name(); filepath.Ext(name) == ".locale" {
				c := config.NewConfig(config.INI)
				if err = c.ParseFile(path); err == nil {
					if values := c.SectionVals(c.DefSec()); len(values) != 0 {
						locale := filepath.Base(name)
						tr.locales[locale] = values
					}
				}
			}
		}
		return err
	})
}

// SetDefaultLocale setup default locale
func (tr *translation) setDefaultLocale(locale string) (err error) {
	if tr.defLocale = tr.locales[locale]; tr.defLocale == nil {
		err = Errorf("Default locale %s has not been loaded", locale)
	}
	return
}

// locale return locale-specified data, if locale is empty, use default locale
func (tr *translation) locale(locale string) (l map[string]string) {
	if l = tr.locales[locale]; l == nil {
		l = tr.defLocale
	}
	return
}

// Translate translate a message to locale-specified string
func (tr *translation) translate(locale, message string) string {
	return tr.locale(locale)[message]
}
