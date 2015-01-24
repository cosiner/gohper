package context

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/cosiner/golib/errors"
	"github.com/cosiner/gomodule/cache"
	"github.com/cosiner/gomodule/config"
	"github.com/cosiner/gomodule/log"
	"github.com/cosiner/gomodule/redis"
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
	forbiddenHandler        = errorHandlerBuilder(http.StatusForbidden)
	notFoundHandler         = errorHandlerBuilder(http.StatusNotFound)
	methodNotAllowedHandler = errorHandlerBuilder(http.StatusMethodNotAllowed)
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
		get:    methodNotAllowedHandler,
		post:   methodNotAllowedHandler,
		put:    methodNotAllowedHandler,
		delete: methodNotAllowedHandler,
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
var tmplSuffix = map[string]bool{"tmpl": true, "html": true}

func isTemplate(name string) (is bool) {
	index := strings.LastIndex(name, ".")
	if is = (index >= 0); is {
		is = tmplSuffix[name[index+1:]]
	}
	return
}

type sessionNode struct {
	refs int
	sess *Session
}
type Server struct {
	cache.Cache
	*Router
	tmpl                    *template.Template
	NotFoundHandler         HandlerFunc
	ForbiddenHandler        HandlerFunc
	MethodNotAllowedHandler HandlerFunc
	SessionStore            SessionStore
	sessionExpire           uint64
	sessions                map[string]*sessionNode
	sessionLock             *sync.Mutex
}

func NewServer() *Server {
	return &Server{
		ForbiddenHandler:        forbiddenHandler,
		NotFoundHandler:         notFoundHandler,
		MethodNotAllowedHandler: methodNotAllowedHandler,
		sessions:                make(map[string]*sessionNode),
		sessionLock:             new(sync.Mutex),
	}
}

func generateId() string {
	return strconv.Itoa(int(time.Now().UnixNano()) / 10)
}

func (s *Server) session(id string) *Session {
	s.sessionLock.Lock()
	sn := s.sessions[id]
	sn.refs++
	s.sessionLock.Unlock()
	return sn.sess
}

func (s *Server) addSession(sess *Session) *Session {
	sn := &sessionNode{refs: 1, sess: sess}
	s.sessionLock.Lock()
	s.sessions[sess.id] = sn
	s.sessionLock.Unlock()
	return sess
}

func (s *Server) getSession(id string) *Session {
	return s.addSession(getSession(s, id, s.SessionStore))
}

func (s *Server) newSession() *Session {
	return s.addSession(newSession(s, generateId()))
}

func (s *Server) storeSession(sess *Session) {
	id := sess.id
	s.sessionLock.Lock()
	sn := s.sessions[id]
	if sn == nil {
		panic("Unexpected: session haven't stored in server.sessions")
	}
	sn.refs--
	if sn.refs == 0 {
		delete(s.sessions, id)
		sn.sess.store(s.SessionStore)
	}
	s.sessionLock.Unlock()
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

// TODO:response add header cookie to store session id
func (s *Server) readySession(resp *Response, req *Request) {
	var sess *Session
	if id := req.SessionId(); id == "" {
		sess = s.newSession()
		resp.setSessionCookie(id)
	} else {
		if sess = s.session(id); sess == nil {
			sess = s.getSession(id)
			resp.setSessionCookie(id)
		}
	}
	req.session = sess
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
			handleFunc = s.MethodNotAllowedHandler
		} else if s.SessionStore != nil {
			s.readySession(resp, req)
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
		c := config.NewConfig(config.LINE)
		s.sessionExpire = uint64(c.IntValDef("expire", 600))
		if err = s.SessionStore.Init(c.SectionVals(c.DefSec())); err != nil {
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
	Init(conf map[string]string) error
	Save(id string, values map[string]interface{}, expire uint64) error
	Get(id string) (map[string]interface{}, error)
}

type redisStore struct {
	store *redis.RedisStore
}

func (rstore *redisStore) Init(conf map[string]string) error {
	return rstore.store.InitWith(conf)
}

func (rstore *redisStore) Save(id string, values map[string]interface{}, expire uint64) (err error) {
	var (
		buffer  = bytes.NewBuffer([]byte{})
		encoder = gob.NewEncoder(buffer)
	)
	if err = encoder.Encode(values); err == nil {
		err = rstore.store.SetWithExpire(id, buffer.Bytes(), expire)
	}
	return
}

func (rstore *redisStore) Get(id string) (vals map[string]interface{}, err error) {
	if bs, err := redis.ToBytes(rstore.store.Get(id)); err == nil {
		vals = make(map[string]interface{})
		if len(bs) != 0 {
			decoder := gob.NewDecoder(bytes.NewBuffer(bs))
			err = decoder.Decode(vals)
		}
	}
	return
}

type Session struct {
	id     string
	server *Server
	lock   *sync.RWMutex
	values map[string]interface{}
}

func newSession(s *Server, id string) *Session {
	return &Session{
		id:     id,
		server: s,
		lock:   new(sync.RWMutex),
		values: make(map[string]interface{}),
	}
}

func getSession(s *Server, id string, store SessionStore) *Session {
	values, err := store.Get(id)
	if err == nil {
		if values == nil {
			values = make(map[string]interface{})
		}
		return &Session{
			id:     id,
			server: s,
			lock:   new(sync.RWMutex),
			values: values,
		}
	}
	return nil
}

func (sess *Session) store(store SessionStore) {
	sess.lock.RLock()
	store.Save(sess.id, sess.values, sess.server.sessionExpire)
	sess.lock.RUnlock()
}

func (sess *Session) Set(key string, val interface{}) {
	sess.lock.Lock()
	sess.values[key] = val
	sess.lock.Unlock()
}

func (sess *Session) Get(key string) interface{} {
	sess.lock.RLock()
	val := sess.values[key]
	sess.lock.RUnlock()
	return val
}

const _COOKIE_SESSION = "session"

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
	if c, err := req.Request.Cookie(name); err == nil {
		return c.Value
	}
	return ""
}

func (req *Request) SessionId() string {
	return req.Cookie(_COOKIE_SESSION)
}

//==============================================================================
//                           Respone
//==============================================================================
type Response struct {
	s *Server
	http.ResponseWriter
}

func (resp *Response) SetCookie(name, value string) {
	http.SetCookie(resp, &http.Cookie{Name: name, Value: value})
}

func (resp *Response) setSessionCookie(id string) {
	resp.SetCookie(_COOKIE_SESSION, id)
}

func (resp *Response) Redirect(req *Request, url string) {
	http.Redirect(resp, req.Request, url, http.StatusTemporaryRedirect)
}

func (resp *Response) Render(tmplName string, val interface{}) error {
	return resp.s.tmpl.ExecuteTemplate(resp, tmplName, val)
}

func (resp *Response) WriteString(str string) error {
	_, err := io.WriteString(resp, str)
	return err
}

func (resp *Response) WriteJson(val interface{}) error {
	bs, err := json.Marshal(val)
	if err == nil {
		_, err = resp.Write(bs)
	}
	return err
}

func (resp *Response) WriteXml(val interface{}) error {
	bs, err := xml.Marshal(val)
	if err == nil {
		_, err = resp.Write(bs)
	}
	return err
}
