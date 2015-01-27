package server

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cosiner/golib/errors"
)

//==============================================================================
//                           Server Init
//==============================================================================
// Server represent a web server
type Server struct {
	AttrContainer
	*Router
	*sessionManager
	ErrorHandlers
	tmpl          *template.Template
	filterForward bool
}

// NewServer create a new server
func NewServer() *Server {
	return &Server{
		AttrContainer: NewAttrContainer(),
		Router:        NewRouter(),
		ErrorHandlers: NewErrorHandlers(),
	}
}

// Start start server
func (s *Server) Start(listenAddr string) {
	log.Println("Compile Templates")
	OnErrPanic(s.compileTemplates())
	log.Println("Initial I18N Locales")
	localeFiles := strach.localeFiles()
	if len(localeFiles) > 0 {
		for _, localeFile := range localeFiles {
			OnErrPanic(_tr.Load(localeFile))
		}
		defaultLocale := strach.defaultLocale()
		if defaultLocale == "" {
			defaultLocale = "en_US"
		}
		OnErrPanic(_tr.SetDefaultLocale(defaultLocale))
	}
	log.Println("Init Session Container")
	store, conf, expire := strach.sessionStore()
	if store == nil {
		store = new(memStore)
		conf = DEF_SESSION_MEMSTORE_CONF
		expire = DEF_SESSION_EXPIRE
	}
	OnErrPanic(store.Init(conf))
	s.sessionManager = newSessionManager(store, expire)
	log.Println("Init Handlers")
	s.initHandler(func(handler Handler) bool {
		OnErrPanic(handler.Init(s))
		return true
	})
	strach.destroy()
	http.ListenAndServe(listenAddr, s)
	s.destroy()
}

//==============================================================================
//                          Server Process
//==============================================================================
// ServHttp serve for http reuest
// find handler and resolve path, find filters, process
func (s *Server) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	resp, req := s.setupContext(w, request)
	s.serve(request.URL, resp, req, false)  // process
	if sess := req.Session(); sess != nil { // store session
		s.storeSession(sess)
	}
	resp.destroy() // destroy request, response
	req.destroy()
}

// setupContext set up context for request and response
func (s *Server) setupContext(w http.ResponseWriter, request *http.Request) (
	*Response, *Request) {
	ctx := newContext(s, w, request)
	resp := newResponse(ctx, w)
	req := newRequest(ctx, request)
	ctx.setup(resp, req)
	return resp, req
}

// serve do actually serve for request
func (s *Server) serve(url *url.URL, resp *Response, req *Request, forward bool) {
	var handlerFunc HandlerFunc
	path := url.Path
	method := parseRequestMethod(req.Method)
	req.Method = method
	handler, urlVars := s.handler(path)
	if handler != nil {
		req.setUrlVars(urlVars)
		handlerFunc = indicateMethod(handler, method)
		if handlerFunc != nil {
			handleWithFilter(resp, req, handlerFunc, path, forward)
			return // normal return
		} else { // find handler but no handle function for this method
			handlerFunc = s.MethodNotAllowedHandler()
		}
	} else { // no handler means no resource there
		handlerFunc = s.NotFoundHandler()
	}
	handlerFunc(resp, req) // error handle
}

// handleWithFilter handle request and response
// if request is forward and server is configured to filter forward request
// filters will not be triggered
func (s *Server) handleWithFilter(resp *Response, req *Request, handlerFunc HandlerFunc,
	path string, forward bool) {
	if forward && !s.FilterForward() {
		HandlerFunc(resp, req)
	} else {
		filters := s.filters(path)
		i, l := 0, len(filters)
		for ; i < l; i++ { // do before in added order
			if !filters[i].Before(resp, req) {
				break
			}
		}
		handlerFunc(resp, req)
		for i = l - 1; i >= 0; i-- { // do after in reverse order
			if !filters[i].After(resp, req) {
				break
			}
		}
	}
}

//==============================================================================
//                          Server Session
//==============================================================================
// SetSessionStore set session store for server with given conf and expire time
// nil store is not accepted
func (s *Server) SetSessionStore(store SessionStore, conf string, expire uint64) {
	if store == nil {
		return
	}
	strach.setSessionStore(store, conf, expire)
}

//==============================================================================
//                           Server I18N
//==============================================================================
// AddLocale add an locale file or dir
func (s *Server) AddLocale(path string) {
	strach.addLocaleFile(path)
}

// SetDefaultLocale set default locale for i18n
func (s *Server) SetDefaultLocale(locale string) {
	strach.setDefaultLocale(locale)
}

//==============================================================================
//                           Server Templates
//==============================================================================
var (
	// globalTmplFuncs is the default template functions
	globalTmplFuncs = map[string]interface{}{
		"I18N": I18N,
	}
	// tmplSuffixes is all template file's suffix
	tmplSuffixes = map[string]bool{"tmpl": true, "html": true}
)

// isTemplate check whether a file name is recognized template file
func isTemplate(name string) (is bool) {
	index := strings.LastIndex(name, ".")
	if is = (index >= 0); is {
		is = tmplSuffixes[name[index+1:]]
	}
	return
}

// AddTemplateSuffix add an suffix for template
func (s *Server) AddTemplateSuffix(suffix string) {
	if suffix != "" {
		if suffix[0] == '.' {
			suffix = suffix[1:]
		}
		tmplSuffixes[suffix] = true
	}
}

// SetTemplateDelims set default template delimeters
func (s *Server) SetTemplateDelims(left, right string) {
	strach.setTmplDelims(left, right)
}

// AddTemplates add templates to server, all templates will be
// parsed on server start
func (s *Server) AddTemplates(names ...string) (err error) {
	addTmpl := func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && isTemplate(path) {
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
func (s *Server) compileTemplates() (err error) {
	var tmpl *template.Template
	tmpl, err = template.New("tmpl").
		Delims(strach.tmplDelims()).
		Funcs(globalTmplFuncs).
		ParseFiles(strach.tmpls()...)
	if err == nil {
		s.tmpl = tmpl
	}
	return
}

// RegisterTemplateFunc register a function used in templates
func (s *Server) RegisterTemplateFunc(name string, fn interface{}) {
	globalTmplFuncs[name] = fn
}

// RegisterTemplateFuncs register some functions used in templates
func (s *Server) RegisterTemplateFuncs(funcs map[string]interface{}) {
	for name, fn := range funcs {
		s.RegisterTemplateFunc(name, fn)
	}
}

// RenderTemplate render a template with given name use given
// value to given writer
func (s *Server) RenderTemplate(wr io.Writer, name string, val interface{}) error {
	return s.tmpl.ExecuteTemplate(wr, name, val)
}

//==============================================================================
//                           Server Handler
//==============================================================================
// Get register a function handler process GET request for given pattern
func (s *Server) Get(pattern string, handlerFunc HandlerFunc) {
	s.addFuncHandler(pattern, GET, handlerFunc)
}

// Post register a function handler process POST request for given pattern
func (s *Server) Post(pattern string, handlerFunc HandlerFunc) {
	s.addFuncHandler(pattern, POST, handlerFunc)
}

// Put register a function handler process PUT request for given pattern
func (s *Server) Put(pattern string, handlerFunc HandlerFunc) {
	s.addFuncHandler(pattern, PUT, handlerFunc)
}

// Delete register a function handler process DELETE request for given pattern
func (s *Server) Delete(pattern string, handlerFunc HandlerFunc) {
	s.addFuncHandler(pattern, DELETE, handlerFunc)
}

//==============================================================================
//                           Server Filter
//==============================================================================
// globalFilters is default filters, slice type is in order to keep filter order
var globalFilters = []map[string]Filter{
	map[string]Filter{},
}

// Before register a function filter executed before handler for given url pattern
func (s *Server) Before(pattern string, filterFunc FilterFunc) {
	s.addFuncFilter(pattern, _FILTER_BEFORE, filterFunc)
}

// After register a function filter executed after handler for given url pattern
func (s *Server) After(pattern string, filterFunc FilterFunc) {
	s.addFuncFilter(pattern, _FILTER_AFTER, filterFunc)
}

// FilterForward also filter forward request, default false
func (s *Server) FilterForward() {
	s.filterForward = true
}
