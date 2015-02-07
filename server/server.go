package server

import (
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/cosiner/gomodule/websocket"

	. "github.com/cosiner/golib/errors"
)

// serverStart do only once
var (
	serverStart = new(sync.Once)
	serverInit  = new(sync.Once)
)

//==============================================================================
//                           Server Init
//==============================================================================
type (
	// ServerConfig is all config of server
	ServerConfig struct {
		Router             Router             // router
		ErrorHandlers      ErrorHandlers      // error handlers
		SessionDisable     bool               // session disble
		SessionManager     SessionManager     // session manager
		SessionStore       SessionStore       // session store
		StoreConfig        string             // session store config
		SessionLifetime    int64              // session lifetime
		TemplateEngine     TemplateEngine     // template engine
		FilterForward      bool               // fliter forward request
		XsrfEnable         bool               // Enable xsrf cookie
		XsrfFlushInterval  int                // xsrf cookie value flush interval
		XsrfLifetime       int                // xsrf cookie expire
		XsrfTokenGenerator XsrfTokenGenerator // xsrf token generator
	}

	// Server represent a web server
	Server struct {
		AttrContainer
		Router
		SessionManager
		ErrorHandlers
		TemplateEngine
		xsrf          Xsrf
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

// init
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
	if sc.TemplateEngine == nil {
		sc.TemplateEngine = NewTemplateEngine()
	}
}

// Init init server with given config, it do only once, if the config contains
// router, it must be called before any operation of add Handler/Filter/WebsocketHandler
// otherwise, ann added Handler/Filter/WebsocketHandler will be added to the old router
// templates and locales are inited last, so Handler/Filter/WebSockethandler's Init
// can use server to add templates and other resources
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
	if srvConf.XsrfEnable {
		log.Println("Init xsrf")
		xsrfFlush := srvConf.XsrfFlushInterval
		if xsrfFlush <= 0 {
			xsrfFlush = XSRF_DEFFLUSH
		}
		xsrfLifetime := srvConf.XsrfLifetime
		if xsrfLifetime <= 0 {
			xsrfLifetime = XSRF_DEFLIFETIME
		}
		xsrfGen := srvConf.XsrfTokenGenerator
		s.xsrf = NewXsrf(xsrfGen, xsrfLifetime)
	} else {
		s.xsrf = emptyXsrf{}
	}

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
	s.Router.Init(func(handler Handler) {
		OnErrPanic(handler.Init(s))
	}, func(filter Filter) {
		OnErrPanic(filter.Init(s))
	}, func(websocketHandler WebSocketHandler) {
		OnErrPanic(websocketHandler.Init(s))
	})

	log.Println("Init Error Handlers")
	s.ErrorHandlers = srvConf.ErrorHandlers
	s.ErrorHandlers.Init(s)

	log.Println("Compile Templates")
	s.TemplateEngine = srvConf.TemplateEngine
	OnErrPanic(s.CompileTemplates())

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

// serveWebSocket serve for websocket protocal
func (s *Server) serveWebSocket(w http.ResponseWriter, request *http.Request) {
	handler, indexer := s.MatchWebSocketHandler(request.URL)
	if handler == nil {
		w.WriteHeader(http.StatusNotFound)
	} else if conn, err := websocket.UpgradeWebsocket(w, request, nil); err == nil {
		handler.Handle(newWebSocketConn(conn).setUrlVarIndexer(indexer))
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
	var (
		xsrfError   bool
		handlerFunc HandlerFunc
		method      = req.Method()
		matchFunc   = s.MatchHandlerFilters
	)
	if !forward { // check xsrf error
		if method == GET {
			resp.setXsrfToken(s.xsrf.Set(req, resp))
		} else {
			xsrfError = !s.xsrf.IsValid(req)
		}
	} else if !s.filterForward {
		matchFunc = s.MatchHandler
	}
	handler, indexer, filters := matchFunc(url)
	handlerFunc = s.indicateHandlerFunc(method, xsrfError, handler)
	if filters == nil {
		handlerFunc(req.setVarIndexer(indexer), resp)
	} else {
		NewFilterChain(filters, handlerFunc).Filter(req.setVarIndexer(indexer), resp)
	}
}

// indicateHandlerFunc indicate handler function from method, has xsrf error and handler
func (s *Server) indicateHandlerFunc(method string, xsrfError bool, handler Handler) (
	handlerFunc HandlerFunc) {
	if handler != nil {
		if handlerFunc = IndicateHandler(method, handler); handlerFunc == nil {
			handlerFunc = s.MethodNotAllowedHandler()
		} else if xsrfError {
			if handler, is := handler.(XsrfErrorHandler); is {
				handlerFunc = handler.HandleXsrfError
			} else {
				handlerFunc = s.XsrfErrorHandler()
			}
		}
	} else { // no handler means no resource there
		handlerFunc = s.NotFoundHandler()
	}
	return
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
