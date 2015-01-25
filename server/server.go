package server

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosiner/gomodule/log"
)

//==============================================================================
//                           Server
//==============================================================================
const HEADER_CONTENTTYPE = "Content-Type"
const HEADER_CONTENTLENGTH = "Content-Length"
const HEADER_SETCOOKIE = "Set-Cookie"

type Server struct {
	*AttrContainer
	Router
	tmpl                    *template.Template
	NotFoundHandler         HandlerFunc
	ForbiddenHandler        HandlerFunc
	MethodNotAllowedHandler HandlerFunc
	*serverSession
}

func NewServer() *Server {
	return &Server{
		AttrContainer:           NewAttrContainer(),
		Router:                  NewRouter(),
		tmpl:                    nil,
		ForbiddenHandler:        forbiddenHandler,
		NotFoundHandler:         notFoundHandler,
		MethodNotAllowedHandler: methodNotAllowedHandler,
		serverSession:           nil,
	}
}

func (s *Server) IsSessionEnabled() bool {
	return s.serverSession == nil
}

func (s *Server) SetSessionStore(store SessionStore, conf string, expire uint64) {
	strach.setSessionStoreConf(conf)
	s.serverSession = newServerSession(store, expire)
}

func withFilterHandle(resp *Response, req *Request, handleFunc HandlerFunc,
	filters []Filter) {
	for _, f := range filters {
		if !f.Before(resp, req) {
			break
		}
	}
	handleFunc(resp, req)
	for _, f := range filters {
		if !f.After(resp, req) {
			break
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	var (
		handleFunc = s.NotFoundHandler
		resp       = newResponse(s, w)
		req        = newRequest(s, request, resp)
		method     = parseRequestMethod(request.Method)
	)
	request.Method = method
	path := request.URL.Path
	handler, urlStories := s.handler(path)
	filters := s.filters(path)
	if handler != nil {
		req.setUrlStories(urlStories)
		var mi MethodIndicator
		switch handler := handler.(type) {
		case MethodIndicator:
			mi = handler
		default:
			mi = methodIndicator{Handler: handler}
		}
		if handleFunc = mi.Method(method); handleFunc != nil {
			if s.IsSessionEnabled() {
				if id := req.sessionId(); id != "" {
					req.setSession(s.session(id))
				}
			}
			withFilterHandle(resp, req, handleFunc, filters)
			return
		} else {
			handleFunc = s.MethodNotAllowedHandler
		}
	}
	handleFunc(resp, req)
}

func (s *Server) Start(listenAddr, sessionConf string) {
	log.Init(log.DEF_FLUSHINTERVAL, log.LEVEL_DEBUG)
	log.AddConsoleWriter("")
	var hasErr bool
	log.Debugln("Compile Templates")
	err := s.CompileTemplate()
	if err != nil {
		hasErr = true
		log.Errorln(err)
	}
	strach.destroy()
	log.Debugln("Init Session Container")
	if s.serverSession != nil {
		if err = s.initStore(strach.sessionStoreConf()); err != nil {
			hasErr = true
			log.Errorln(err)
		}
	}
	log.Debugln("Init Handlers")
	s.initHandler(func(handler Handler) bool {
		if err := handler.Init(s); err != nil {
			hasErr = true
			log.Errorln(err)
		}
		return true
	})
	if hasErr {
		log.Fatal()
	}
	http.ListenAndServe(listenAddr, s)
}

//==============================================================================
//                           Server Templates
//==============================================================================
var tmplSuffix = map[string]bool{"tmpl": true, "html": true}

func isTemplate(name string) (is bool) {
	index := strings.LastIndex(name, ".")
	if is = (index >= 0); is {
		is = tmplSuffix[name[index+1:]]
	}
	return
}

func (s *Server) AddTemplateSuffix(suffix string) {
	if suffix != "" {
		if suffix[0] == '.' {
			suffix = suffix[1:]
		}
		tmplSuffix[suffix] = true
	}
}

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

func (s *Server) CompileTemplate() (err error) {
	if s.tmpl == nil {
		var tmpl *template.Template
		tmpl, err = template.New("tmpl").
			Funcs(strach.tmplFuncs()).ParseFiles(strach.tmpls()...)
		if err == nil {
			s.tmpl = tmpl
		}
	}
	return
}

func (s *Server) RegisterTemplateFunc(name string, fn interface{}) {
	strach.setTmplFunc(name, fn)
}

func (s *Server) RegisterTemplateFuncs(funcs map[string]interface{}) {
	strach.setTmplFuncs(funcs)
}

func (s *Server) RenderTemplate(wr io.Writer, name string, val interface{}) error {
	return s.tmpl.ExecuteTemplate(wr, name, val)
}

//==============================================================================
//                           Server Handler
//==============================================================================
func (s *Server) Get(pattern string, handlerFunc HandlerFunc) {
	s.addFuncHandler(pattern, GET, handlerFunc)
}

func (s *Server) Post(pattern string, handlerFunc HandlerFunc) {
	s.addFuncHandler(pattern, POST, handlerFunc)
}

func (s *Server) Put(pattern string, handlerFunc HandlerFunc) {
	s.addFuncHandler(pattern, PUT, handlerFunc)
}

func (s *Server) Delete(pattern string, handlerFunc HandlerFunc) {
	s.addFuncHandler(pattern, DELETE, handlerFunc)
}

//==============================================================================
//                           Server Filter
//==============================================================================
func (s *Server) Before(pattern string, filterFunc FilterFunc) {
	s.addFuncFilter(pattern, _BEFORE, filterFunc)
}

func (s *Server) After(pattern string, filterFunc FilterFunc) {
	s.addFuncFilter(pattern, _AFTER, filterFunc)
}
