// this file is strach before server start, and will be destroied
// at the moment of server init
package server

import (
	"github.com/cosiner/golib/regexp/urlmatcher"
)

type _strach struct {
	_tmplNames        []string                       // template names
	_tmplDelims       [2]string                      // template delimeters
	_funcHandlers     map[string]*funcHandler        // function handlers
	_funcFilters      map[string]*funcFilter         // function filters
	_routeMatchers    map[string]*urlmatcher.Matcher // route matchers
	_sessionStore     SessionStore                   // sessionStore
	_sessionStoreConf string                         // session store config
	_sessionLifetime  uint64                         // session lifetime
	_localeFiles      []string                       // locale files
	_defaultLocale    string                         // default locale
}

var strach = &_strach{
	_tmplNames:     make([]string, 10),
	_tmplDelims:    [2]string{"{{", "}}"},
	_funcHandlers:  make(map[string]*funcHandler),
	_funcFilters:   make(map[string]*funcFilter),
	_routeMatchers: make(map[string]*urlmatcher.Matcher),
	_localeFiles:   make([]string, 2),
}

func (s *_strach) destroy() {
	s._tmplNames = nil
	s._funcHandlers = nil
	s._funcFilters = nil
	s._routeMatchers = nil
	s._localeFiles = nil
	s._sessionStore = nil
}

func (s *_strach) addTmpl(name string) {
	s._tmplNames = append(s._tmplNames, name)
}

func (s *_strach) tmpls() []string {
	return s._tmplNames
}

func (s *_strach) setTmplDelims(left, right string) {
	s._tmplDelims[0] = left
	s._tmplDelims[1] = right
}

func (s *_strach) tmplDelims() (string, string) {
	return s._tmplDelims[0], s._tmplDelims[1]
}

func (s *_strach) setFuncHandler(pattern string, handler *funcHandler) {
	s._funcHandlers[pattern] = handler
}

func (s *_strach) funcHandler(pattern string) *funcHandler {
	return s._funcHandlers[pattern]
}

func (s *_strach) setFuncFilter(pattern string, filter *funcFilter) {
	s._funcFilters[pattern] = filter
}

func (s *_strach) funcFilter(pattern string) *funcFilter {
	return s._funcFilters[pattern]
}

func (s *_strach) setRouteMatcher(pattern string, matcher *urlmatcher.Matcher) {
	s._routeMatchers[pattern] = matcher
}

func (s *_strach) routeMatcher(pattern string) *urlmatcher.Matcher {
	return s._routeMatchers[pattern]
}

func (s *_strach) setSessionStore(store SessionStore, conf string, lifetime uint64) {
	s._sessionStore = store
	s._sessionStoreConf = conf
	s._sessionLifetime = lifetime
}

func (s *_strach) sessionStore() (SessionStore, string, uint64) {
	return s._sessionStore, s._sessionStoreConf, s._sessionLifetime
}

func (s *_strach) addLocaleFile(path string) {
	s._localeFiles = append(s._localeFiles, path)
}
func (s *_strach) localeFiles() []string {
	return s._localeFiles
}

func (s *_strach) setDefaultLocale(locale string) {
	s._defaultLocale = locale
}
func (s *_strach) defaultLocale() string {
	return s._defaultLocale
}
