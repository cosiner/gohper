package context

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cosiner/golib/errors"

	"github.com/cosiner/gomodule/memcache"
	"github.com/cosiner/gomodule/rediscache"
)

type Server struct {
	memcache.MemCache
	*Router
	tmpl         *template.Template
	sessionStore *rediscache.RedisCache
}

const DEF_TMPL_SUFFIX = ".tpl"

func (s *Server) AddTemplates(names ...string) (err error) {
	addTmpl := func(path string, info *os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), DEF_TMPL_SUFFIX) {
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

type HandlerFunc func(response *Response, request *Request)
type Handler interface {
	Init(server *Server) error
	Get(*Response, *Request)
	Post(*Response, *Request)
	Delete(*Response, *Request)
	Put(*Response, *Request)
	Finish()
}

type funcHandler struct {
	Get    HandlerFunc
	Post   HandlerFunc
	Delete HandlerFunc
	Put    HandlerFunc
}

func newFuncHandler() *funcHandler {
	return &funcHandler{
		Get:    ForbiddenHandle,
		Post:   ForbiddenHandle,
		Put:    ForbiddenHandle,
		Delete: ForbiddenHandle,
	}
}

func (fh *funcHandler) Init(s *Server) error {}
func (fh *funcHandler) Finish()              {}

func ForbiddenHandle(response *Response, _ *Request) {
	response.WriteString("403 Forbidden")
}

func (s *Server) Get(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncRoute(s, pattern, GET, handleFunc)
}

func (s *Server) Post(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncRoute(s, pattern, POST, handleFunc)
}

func (s *Server) Put(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncRoute(s, pattern, PUT, handleFunc)
}

func (s *Server) Delete(pattern string, handlerFunc HandlerFunc) {
	s.AddFuncRoute(s, pattern, DELETE, handleFunc)
}

func (s *Server) ServeHttp(w http.ResponseWriter, req *http.Request) {
	var (
		handleFunc = ForbiddenHandle
		response   = &Response{w}
		request    = &Request{Request: req, server: s}
		method     = parseRequestMethod(req.Method)
	)
	http.Handler
	req.Method = method
	handler, urlStories := s.Handler(req.URL.Path)
	if handler != nil {
		request.urlStories = urlStories
		switch method {
		case GET:
			handleFunc = handler.Get
		case POST:
			handleFunc = handler.Post
		case DELETE:
			handleFunc = handler.Delete
		case PUT:
			handleFunc = handler.Put
		}
	}
	go handleFunc(request, response)
}

func (s *Server) Start(addr string) {
	var hasErr bool
	err := s.CompileTemplate()
	if err != nil {
		hasErr = true
		log.Println(err)
	}
	strachDestroy()
	s.Router.InitHandler(func(handler Handler) bool {
		if err := handler.Init(s); err != nil {
			hasErr = true
			log.Println(err)
		}
	})
	if hasErr {
		os.Exit(1)
	}
	http.ListenAndServe(addr, s)
}

func NewServer() *Server {
	return new(Server)
}

type Session struct {
	id           string
	server       *Server
	sessionStore *rediscache.RedisCache
}

type Request struct {
	server     *Server
	session    *Session
	urlStories map[string]string
	*http.Request
}

func (req *Request) Session() *Session {
	return req.session
}

func (req *Request) setSession(s *Session) {
	req.session = s
}

func (req *Request) Cookie(name string) string {
}

type Response struct {
	s *Server
	http.ResponseWriter
}

func (r *Response) Render(tmplName string, val interface{}) error {
	return r.s.tmpl.ExecuteTemplate(r, tmplName, val)
}

func (r *Response) WriteString(str string) error {
	io.WriteString(r, s)
}

func (r *Response) WriteJson(val interface{}) error {
	val, err := json.Marshal(val)
	if err == nil {
		_, err = r.Write(val)
	}
	return err
}
