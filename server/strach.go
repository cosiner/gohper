// this file is strach before server start, and will be destroied
// at the moment of server init
package server

import (
	"github.com/cosiner/golib/regexp/urlmatcher"
)

type _strach struct {
	_tmplNames        []string                       // template names
	_tmplFuncs        map[string]interface{}         // template functions
	_funcHandlers     map[string]*funcHandler        // function handlers
	_funcFilters      map[string]*funcFilter         // function filters
	_routeMatchers    map[string]*urlmatcher.Matcher // route matchers
	_sessionStoreConf string
}

var strach = &_strach{
	_tmplNames:     make([]string, 10),
	_tmplFuncs:     make(map[string]interface{}),
	_funcHandlers:  make(map[string]*funcHandler),
	_funcFilters:   make(map[string]*funcFilter),
	_routeMatchers: make(map[string]*urlmatcher.Matcher),
}

func (s *_strach) addTmpl(name string) {
	s._tmplNames = append(s._tmplNames, name)
}

func (s *_strach) tmpls() []string {
	return s._tmplNames
}

func (s *_strach) setTmplFunc(name string, fn interface{}) {
	s._tmplFuncs[name] = fn
}

func (s *_strach) setTmplFuncs(funcs map[string]interface{}) {
	if len(s._tmplFuncs) == 0 {
		s._tmplFuncs = funcs
	} else {
		for k, v := range funcs {
			s._tmplFuncs[k] = v
		}
	}
}

func (s *_strach) tmplFuncs() map[string]interface{} {
	return s._tmplFuncs
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

func (s *_strach) sessionStoreConf() string {
	return s._sessionStoreConf
}

func (s *_strach) setSessionStoreConf(conf string) {
	s._sessionStoreConf = conf
}

func (s *_strach) destroy() {
	s._tmplNames = nil
	s._tmplFuncs = nil
	s._funcHandlers = nil
	s._funcFilters = nil
	s._routeMatchers = nil
}
