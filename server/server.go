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

	"github.com/cosiner/gomodule/websocket"

	. "github.com/cosiner/golib/errors"
)

// serverStart do only once
var serverStart = new(sync.Once)
var serverInit = new(sync.Once)

//==============================================================================
//                           Server Init
//==============================================================================
type (
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
		AttrContainer: NewLockedAttrContainer(),
		Router:        NewRouter(),
	}
}

func (sc *ServerConfig) init() {
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
	} else {
		sc.SessionManager = newPanicSessionManager()
	}
}

// Init init server with given config, it do only once, if the config contains
// router, it must be called before any operation of add Handler/Filter/WebsocketHandler
// otherwise, ann added Handler/Filter/WebsocketHandler will be added to the old router
func (s *Server) Init(conf *ServerConfig) {
	serverInit.Do(func() {
		if conf.Router != nil {
			s.Router = conf.Router
		}
		strach.setServerConfig(conf)
	})
}

// Start start server
func (s *Server) start() {
	srvConf := strach.serverConfig()
	if srvConf == nil {
		srvConf = new(ServerConfig)
	}
	srvConf.init()
	s.filterForward = srvConf.FilterForward
	manager := srvConf.SessionManager
	if !srvConf.SessionDisable {
		log.Println("Init Session Store and Manager")
		store := srvConf.SessionStore
		OnErrPanic(store.Init(srvConf.StoreConfig))
		OnErrPanic(manager.Init(store, srvConf.SessionLifetime))
	}
	s.SessionManager = manager

	log.Println("Init Handlers and Filters")
	s.Router.Init(func(handler Handler) bool {
		OnErrPanic(handler.Init(s))
		return true
	}, func(filter Filter) bool {
		OnErrPanic(filter.Init(s))
		return true
	}, func(websocketHandler WebSocketHandler) bool {
		OnErrPanic(websocketHandler.Init(s))
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
func (s *Server) Start(listenAddr string) error {
	serverStart.Do(s.start)
	return http.ListenAndServe(listenAddr, s)
}

// StartTLS start server as https server
func (s *Server) StartTLS(listenAddr, certFile, keyFile string) error {
	serverStart.Do(s.start)
	return http.ListenAndServeTLS(listenAddr, certFile, keyFile, s)
}

//==============================================================================
//                          Server Process
//==============================================================================
// ServHttp serve for http reuest
// find handler and resolve path, find filters, process
func (s *Server) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if websocket.IsWebSocketRequest(request) {
		s.serveWebSocket(w, request)
	} else {
		s.serveHttp(w, request)
	}
}

// serveHttp serve for http protocal
func (s *Server) serveHttp(w http.ResponseWriter, request *http.Request) {
	req, resp := s.setupContext(w, request)
	s.processHttpRequest(request.URL, req, resp, false) // process
	if req.hasSession() {                               // store session
		s.StoreSession(req.Session())
	}
	resp.destroy() // destroy request, response
	req.destroy()
}

// setupContext set up context for request and response
func (s *Server) setupContext(w http.ResponseWriter, request *http.Request) (
	*request, *response) {
	ctx := newContext(s, w, request)
	resp := newResponse(ctx, w)
	req := newRequest(ctx, request)
	ctx.init(req, resp)
	return req, resp
}

// processHttpRequest do process http request
func (s *Server) processHttpRequest(url *url.URL, req *request, resp *response, forward bool) {
	var handlerFunc HandlerFunc
	handler, indexer, values := s.MatchHandler(url)
	if handler != nil {
		req.setUrlVars(indexer, values)
		if handlerFunc = IndicateHandler(req.Method(), handler); handlerFunc == nil {
			handlerFunc = s.MethodNotAllowedHandler()
		}
	} else { // no handler means no resource there
		handlerFunc = s.NotFoundHandler()
	}
	s.handleWithFilter(req, resp, handlerFunc, url, forward)
}

// serveWebSocket serve for websocket protocal
func (s *Server) serveWebSocket(w http.ResponseWriter, request *http.Request) {
	handler, indexer, values := s.MatchWebSocketHandler(request.URL)
	if handler == nil {
		w.WriteHeader(http.StatusNotFound)
	} else if conn, err := websocket.UpgradeWebsocket(w, request, nil); err == nil {
		handler.Handle(newWebSocketConn(conn).setUrlVars(indexer, values))
	}
}

// handleWithFilter handle request and response
// if request is forward and server is configured to filter forward request
// filters will not be triggered
func (s *Server) handleWithFilter(req Request, resp Response, handlerFunc HandlerFunc,
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
	if tmpls := strach.tmpls(); len(tmpls) != 0 {
		tmpl, err = template.New("tmpl").
			Delims(strach.tmplDelims()).
			Funcs(globalTmplFuncs).
			ParseFiles(strach.tmpls()...)
		if err == nil {
			s.tmpl = tmpl
		}
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

//==============================================================================
//                           Server util
//==============================================================================
// PanicServer panic server by create a new goroutine then panic
func PanicServer(str string) {
	go panic(str)
}
