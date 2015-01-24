package context

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cosiner/golib/errors"

	"github.com/cosiner/gomodule/cache"
	"github.com/cosiner/gomodule/log"
)

//==============================================================================
//                        Request Method
//==============================================================================
const (
	GET            = "GET"
	POST           = "POST"
	DELETE         = "DELETE"
	PUT            = "PUT"
	UNKNOWN_METHOD = "UNKNOWN"
)

func parseRequestMethod(s string) string {
	if s == "" {
		return GET
	}
	return strings.ToUpper(s)
}

type MethodIndicator interface {
	IndicateMethod(method string) HandlerFunc
}

type methodIndicator struct {
	Handler
}

func (s methodIndicator) IndicateMethod(method string) (handleFunc HandlerFunc) {
	switch method {
	case GET:
		handleFunc = s.Get
	case POST:
		handleFunc = s.Post
	case DELETE:
		handleFunc = s.Delete
	case PUT:
		handleFunc = s.Put
	}
	return
}

//==============================================================================
//                         Handler
//==============================================================================
type HandlerFunc func(*Response, *Request)

func errorHandlerBuilder(header int) HandlerFunc {
	return func(resp *Response, req *Request) {
		resp.WriteHeader(header)
	}
}

var (
	ForbiddenHandler = errorHandlerBuilder(http.StatusForbidden)
	NotFoundHandler  = errorHandlerBuilder(http.StatusNotFound)
)

type Handler interface {
	Init(*Server) error
	Get(*Response, *Request)
	Post(*Response, *Request)
	Delete(*Response, *Request)
	Put(*Response, *Request)
	Finish()
}

type funcHandler struct {
	get    HandlerFunc
	post   HandlerFunc
	delete HandlerFunc
	put    HandlerFunc
}

func newFuncHandler() *funcHandler {
	return &funcHandler{
		get:    ForbiddenHandler,
		post:   ForbiddenHandler,
		put:    ForbiddenHandler,
		delete: ForbiddenHandler,
	}
}

func (fh *funcHandler) Get(response *Response, request *Request) {
	fh.get(response, request)
}
func (fh *funcHandler) Post(response *Response, request *Request) {
	fh.post(response, request)
}
func (fh *funcHandler) Put(response *Response, request *Request) {
	fh.Put(response, request)
}
func (fh *funcHandler) Delete(response *Response, request *Request) {
	fh.delete(response, request)
}
func (fh *funcHandler) setMethod(method string, handleFunc HandlerFunc) (err error) {
	switch method {
	case GET:
		fh.get = handleFunc
	case POST:
		fh.post = handleFunc
	case PUT:
		fh.put = handleFunc
	case DELETE:
		fh.delete = handleFunc
	default:
		err = Err("Not supported request method")
	}
	return
}
func (fh *funcHandler) Init(s *Server) error { return nil }
func (fh *funcHandler) Finish()              {}

//==============================================================================
//                           Server
//==============================================================================
const DEF_TMPL_SUFFIX = ".tmpl"

func isTemplate(name string) bool {
	return strings.HasSuffix(name, DEF_TMPL_SUFFIX) ||
		strings.HasSuffix(name, "html")
}

type Server struct {
	cache.Cache
	*Router
	tmpl             *template.Template
	NotFoundHandler  HandlerFunc
	ForbiddenHandler HandlerFunc
	SessionStore     SessionStore
}

func NewServer() *Server {
	s := new(Server)
	s.ForbiddenHandler = ForbiddenHandler
	s.NotFoundHandler = NotFoundHandler
	return s
}

func (s *Server) AddTemplates(names ...string) (err error) {
	addTmpl := func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && isTemplate(path) {
			strachAddTmpl(path)
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
		if tmpl, err = template.ParseFiles(strachTmpls()...); err == nil {
			s.tmpl = tmpl
		}
	}
	return
}

func (s *Server) Get(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncRoute(pattern, GET, handlerFunc)
}

func (s *Server) Post(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncRoute(pattern, POST, handlerFunc)
}

func (s *Server) Put(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncRoute(pattern, PUT, handlerFunc)
}

func (s *Server) Delete(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncRoute(pattern, DELETE, handlerFunc)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	var (
		handleFunc = s.NotFoundHandler
		resp       = &Response{ResponseWriter: w}
		req        = &Request{Request: request, server: s}
		method     = parseRequestMethod(request.Method)
	)
	request.Method = method
	handler, urlStories := s.Handler(request.URL.Path)
	if handler != nil {
		req.urlStories = urlStories
		var mi MethodIndicator
		switch handler := handler.(type) {
		case MethodIndicator:
			mi = handler
		default:
			mi = methodIndicator{Handler: handler}
		}
		if handleFunc = mi.IndicateMethod(method); handleFunc == nil {
			handleFunc = s.ForbiddenHandler
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
	strachDestroy()
	log.Debugln("Init Session Container")
	if s.SessionStore != nil {
		if err = s.SessionStore.Init(sessionConf); err != nil {
			hasErr = true
			log.Errorln(err)
		}
	}
	log.Debugln("Init Handlers")
	s.InitHandler(func(handler Handler) bool {
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
//                           Session
//==============================================================================

type SessionStore interface {
	Init(conf string) error
	IsExist(key string) (bool, error)
	IsHExist(h, key string) (bool, error)
	Get(key string) (interface{}, error)
	HGet(h, key string) (interface{}, error)
	Set(key string, val interface{}) error
	HSet(h, key string, val interface{}) error
	SetExpire(key string, timeout uint64) error
	SetWithExpire(key string, val interface{}, timeout uint64) error
	Remove(key string) error
	Incr(key string) error
	Decr(key string) error
}

type Session struct {
	id     string
	server *Server
	store  SessionStore
}

//==============================================================================
//                           Request
//==============================================================================
type Request struct {
	server     *Server
	session    *Session
	urlStories map[string]string
	*http.Request
}

func (req *Request) ResolveJson() (data map[string]string, err error) {
	var body []byte
	if body, err = ioutil.ReadAll(req.Body); err == nil {
		data = make(map[string]string)
		err = json.Unmarshal(body, data)
	}
	return
}

func (req *Request) Session() *Session {
	return req.session
}

func (req *Request) setSession(s *Session) {
	req.session = s
}

func (req *Request) Cookie(name string) string {
	return ""
}

//==============================================================================
//                           Respone
//==============================================================================
type Response struct {
	s *Server
	http.ResponseWriter
}

func (r *Response) Render(tmplName string, val interface{}) error {
	return r.s.tmpl.ExecuteTemplate(r, tmplName, val)
}

func (r *Response) WriteString(str string) error {
	_, err := io.WriteString(r, str)
	return err
}

func (r *Response) WriteJson(val interface{}) error {
	bs, err := json.Marshal(val)
	if err == nil {
		_, err = r.Write(bs)
	}
	return err
}

func (r *Response) WriteXml(val interface{}) error {
	bs, err := xml.Marshal(val)
	if err == nil {
		_, err = r.Write(bs)
	}
	return err
}
