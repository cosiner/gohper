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
	"sync"

	. "github.com/cosiner/golib/errors"
)

// serverStart do only once
var serverStart = new(sync.Once)

//==============================================================================
//                           Server Init
//==============================================================================
type (
	// Destroyer is a common interface for resource need to destroy on server stoped
	Destroyer interface {
		Destroy()
	}

	// ServerConfig is all config of server
	ServerConfig struct {
		Router          Router
		ErrorHandlers   ErrorHandlers
		SessionDisable  bool
		SessionManager  SessionManager
		SessionStore    SessionStore
		StoreConfig     string
		SessionLifetime int64
		FilterForward   bool
	}

	// Server represent a web server
	Server struct {
		AttrContainer
		Router
		SessionManager
		ErrorHandlers
		tmpl          *template.Template
		filterForward bool
	}
)

// NewServer create a new server
func NewServer() *Server {
	return &Server{
		AttrContainer: NewAttrContainer(),
		Router:        NewRouter(),
		ErrorHandlers: NewErrorHandlers(),
	}
}

func (sc *ServerConfig) init() {
	if sc.Router == nil {
		sc.Router = NewRouter()
	}
	if sc.ErrorHandlers == nil {
		sc.ErrorHandlers = NewErrorHandlers()
	}
	if !sc.SessionDisable {
		if sc.SessionStore == nil {
			sc.StoreConfig = DEF_SESSION_MEMSTORE_CONF
			sc.SessionStore = NewMemStore()
		}
		if sc.SessionManager == nil {
			sc.SessionManager = NewSessionManager()
		}

		if sc.SessionLifetime == 0 {
			sc.SessionLifetime = DEF_SESSION_LIFETIME
		}
	}
}

// Init init server with given config
func (s *Server) Init(conf *ServerConfig) {
	conf.init()
	strach.setServerConfig(conf)
}

// Start start server
func (s *Server) start() {
	srvConf := strach.serverConfig()
	if srvConf == nil {
		srvConf = new(ServerConfig)
		srvConf.init()
	}

	if srvConf.SessionDisable {
		s.SessionManager = newEmptySessionManager()
	} else {
		log.Println("Init Session Store and Manager")
		store, manager := srvConf.SessionStore, srvConf.SessionManager
		OnErrPanic(store.Init(srvConf.StoreConfig))
		OnErrPanic(manager.Init(store, srvConf.SessionLifetime))
		s.SessionManager = manager
	}

	log.Println("Init Handlers and Filters")
	s.Router = srvConf.Router
	s.Router.Init(func(handler Handler) bool {
		OnErrPanic(handler.Init(s))
		return true
	}, func(filter Filter) bool {
		OnErrPanic(filter.Init(s))
		return true
	})

	log.Println("Init Error Handlers")
	s.ErrorHandlers = srvConf.ErrorHandlers

	log.Println("Compile Templates")
	OnErrPanic(s.compileTemplates())

	log.Println("Initial I18N Locales")
	localeFiles := strach.localeFiles()
	if len(localeFiles) > 0 {
		for _, localeFile := range localeFiles {
			OnErrPanic(_tr.load(localeFile))
		}
		defaultLocale := strach.defaultLocale()
		if defaultLocale == "" {
			defaultLocale = "en_US"
		}
		OnErrPanic(_tr.setDefaultLocale(defaultLocale))
	}

	strach.destroy()
	log.Println("Server Start")
}

// Start start server as http server
func (s *Server) Start(listenAddr string) {
	serverStart.Do(s.start)
	http.ListenAndServe(listenAddr, s)
}

// StartTLS start server as https server
func (s *Server) StartTLS(listenAddr, certFile, keyFile string) {
	serverStart.Do(s.start)
	http.ListenAndServeTLS(listenAddr, certFile, keyFile, s)
}

//==============================================================================
//                          Server Process
//==============================================================================
// ServHttp serve for http reuest
// find handler and resolve path, find filters, process
func (s *Server) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	req, resp := s.setupContext(w, request)
	s.serve(request.URL, req, resp, false)  // process
	if sess := req.Session(); sess != nil { // store session
		s.StoreSession(sess)
	}
	resp.destroy() // destroy request, response
	req.destroy()
}

// setupContext set up context for request and response
func (s *Server) setupContext(w http.ResponseWriter, request *http.Request) (
	*Request, *Response) {
	ctx := newContext(s, w, request)
	resp := newResponse(ctx, w)
	req := newRequest(ctx, request)
	ctx.init(req, resp)
	return req, resp
}

// serve do actually serve for request
func (s *Server) serve(url *url.URL, req *Request, resp *Response, forward bool) {
	var handlerFunc HandlerFunc
	method := parseRequestMethod(req.Method())
	req.setMethod(method)
	handler, urlVars := s.MatchHandler(url)
	if handler != nil {
		req.setUrlVars(urlVars)
		handlerFunc = IndicateHandler(method, handler)
		if handlerFunc != nil {
			s.handleWithFilter(req, resp, handlerFunc, url, forward)
			return // normal return
		} else { // find handler but no handle function for this method
			handlerFunc = s.MethodNotAllowedHandler()
		}
	} else { // no handler means no resource there
		handlerFunc = s.NotFoundHandler()
	}
	handlerFunc(req, resp) // error handle
}

// handleWithFilter handle request and response
// if request is forward and server is configured to filter forward request
// filters will not be triggered
func (s *Server) handleWithFilter(req *Request, resp *Response, handlerFunc HandlerFunc,
	url *url.URL, forward bool) {
	var filters []Filter
	if !forward || s.filterForward {
		filters = s.MatchFilters(url)
	}
	newFilterChain(filters, handlerFunc).Filter(req, resp)
}

//==============================================================================
//                           Server I18N
//==============================================================================
// AddLocale add an locale file or dir
func (*Server) AddLocale(path string) {
	strach.addLocaleFile(path)
}

// SetDefaultLocale set default locale for i18n
func (*Server) SetDefaultLocale(locale string) {
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
func (*Server) isTemplate(name string) (is bool) {
	index := strings.LastIndex(name, ".")
	if is = (index >= 0); is {
		is = tmplSuffixes[name[index+1:]]
	}
	return
}

// AddTemplateSuffix add an suffix for template
func (*Server) AddTemplateSuffix(suffix string) {
	if suffix != "" {
		if suffix[0] == '.' {
			suffix = suffix[1:]
		}
		tmplSuffixes[suffix] = true
	}
}

// SetTemplateDelims set default template delimeters
func (*Server) SetTemplateDelims(left, right string) {
	strach.setTmplDelims(left, right)
}

// AddTemplates add templates to server, all templates will be
// parsed on server start
func (s *Server) AddTemplates(names ...string) (err error) {
	addTmpl := func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && s.isTemplate(path) {
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
func (*Server) RegisterTemplateFunc(name string, fn interface{}) {
	globalTmplFuncs[name] = fn
}

// RegisterTemplateFuncs register some functions used in templates
func (*Server) RegisterTemplateFuncs(funcs map[string]interface{}) {
	for name, fn := range funcs {
		globalTmplFuncs[name] = fn
	}
}

// RenderTemplate render a template with given name use given
// value to given writer
func (s *Server) renderTemplate(wr io.Writer, name string, val interface{}) error {
	return s.tmpl.ExecuteTemplate(wr, name, val)
}

//==============================================================================
//                           Server Handler
//==============================================================================
// Get register a function handler process GET request for given pattern
func (s *Server) Get(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncHandler(pattern, GET, handlerFunc)
}

// Post register a function handler process POST request for given pattern
func (s *Server) Post(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncHandler(pattern, POST, handlerFunc)
}

// Put register a function handler process PUT request for given pattern
func (s *Server) Put(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncHandler(pattern, PUT, handlerFunc)
}

// Delete register a function handler process DELETE request for given pattern
func (s *Server) Delete(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncHandler(pattern, DELETE, handlerFunc)
}
