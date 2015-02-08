package server

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	. "github.com/cosiner/golib/errors"
)

type (
	// Router is responsible for route manage and match
	Router interface {
		// Init init handlers and filters, websocket handlers
		Init(func(Handler) bool, func(Filter) bool, func(WebSocketHandler) bool)
		// Destroy destroy router, also responsible for destroy all handlers and filters
		Destroy()

		// AddFuncHandler add a function handler, method are defined as constant string
		AddFuncHandler(pattern string, method string, handler HandlerFunc) error
		// AddHandler add a handler
		AddHandler(pattern string, handler Handler) error
		// AddFuncFilter add function filter
		AddFuncFilter(pattern string, filter FilterFunc) error
		// AddFilter add a filter
		AddFilter(pattern string, filter Filter) error
		// AddFuncWebSocketHandler add a websocket functionhandler
		AddFuncWebSocketHandler(pattern string, handler WebSocketHandlerFunc) error
		// AddWebSocketHandler add a websocket handler
		AddWebSocketHandler(pattern string, handler WebSocketHandler) error

		// MatchHandler match given url to find  final handler, filters must be nil
		MatchHandler(url *url.URL) (Handler, UrlVarIndexer, []Filter)
		// MatchHandlerFilters match given url to find all matched filters and final handler
		MatchHandlerFilters(url *url.URL) (Handler, UrlVarIndexer, []Filter)
		// MatchWebSocketHandler match given url to find a matched websocket handler
		MatchWebSocketHandler(url *url.URL) (WebSocketHandler, UrlVarIndexer)
	}

	// UrlVarIndexer is a indexer for name to value
	UrlVarIndexer interface {
		// ValuesOf return values of variable in given values
		UrlVar(name string) string
		// ScanUrlVars scan given values into variable addresses
		// if address is nil, skip it
		ScanUrlVars(vars ...*string)
		// UrlVars return all values of variable
		UrlVars() []string
	}

	// ValuesPool is a pool for url variable values
	VaulesPool interface {
		Acquire() []string
		Recover([]string)
	}

	// handlerProcessor keep handler and url variables of this route
	handlerProcessor struct {
		vars    map[string]int
		handler Handler
	}
	// wsHandlerProcessor keep websocket handler and url variables of this route
	wsHandlerProcessor struct {
		vars      map[string]int
		wsHandler WebSocketHandler
	}
	// routeProcessor is processor of a route, it can hold handler, filters, and websocket handler
	routeProcessor struct {
		wsHandlerProcessor *wsHandlerProcessor
		handlerProcessor   *handlerProcessor
		filters            []Filter
	}

	// router is a actual url router, it only process path of url, other section is
	// not mentioned
	router struct {
		str       string          // path section hold by current route node
		chars     []byte          // all possible first characters of next route node
		childs    []*router       // child routers
		processor *routeProcessor // processor for current route node
	}

	// pathVars is an implementation of UrlVarIndexer
	pathVars struct {
		vars   map[string]int // url variables and indexs of sections splited by '/'
		values []string       // all url variable values
	}

	pathValuesPool struct {
		pool [][]string
		*sync.RWMutex
	}
)

func newValuesPool(size int) VaulesPool {
	p := new(pathValuesPool)
	p.pool = make([][]string, 0, size)
	p.RWMutex = new(sync.RWMutex)
	return p
}

func (p *pathValuesPool) Acquire() (v []string) {
	p.Lock()
	pool := p.pool
	len := len(pool)
	if len == 0 {
		v = make([]string, 0, PathVarCount)
	} else {
		v = pool[len-1]
		p.pool = pool[:len-1]
	}
	p.Unlock()
	return
}

func (p *pathValuesPool) Recover(v []string) {
	if cap(v) != 0 {
		v = v[:0]
		p.Lock()
		pool := p.pool
		if cap(pool) != len(pool) { // if pool is full, do nothing
			p.pool = append(pool, v)
		}
		p.Unlock()
	}
}

var (
	// nilVars is empty variable map
	nilVars = make(map[string]int)
	// nilVarIndexer is a empty UrlVarIndexer
	nilVarIndexer = &pathVars{vars: nilVars, values: nil}
	// reserveChildsCount is route childs slice increment and init size for addPath
	reserveChildsCount = 1
	// PathVarCount is common url path variable count
	// match functions of router will create a slice use it as capcity to store
	// all path variable values
	// to get best performance, it should commonly set to the average, default, it's 2
	PathVarCount = 2

	valuesPool = newValuesPool(1024)
)

const (
	// _WILDCARD is the replacement of named variable in compiled path
	_WILDCARD = '|' // MUST BE:other character < _WILDCARD < _REMAINSALL
	// _REMAINSALL is the replacement of catch remains all variable in compiled path
	_REMAINSALL = '~'
	// _PRINT_SEP is the seperator of tree level when print route tree
	_PRINT_SEP = "-"
)

// pathSections divide path by '/', trim end '/'and the first '/'
func pathSections(path string) []string {
	if l := len(path); l > 0 {
		if path[l-1] == '/' {
			l = l - 1
		}
		path = path[1:l] // trim first and last '/'
	}
	return strings.Split(path, "/")
}

// isInvalidSection check whether section has the predefined _WILDCARD and match
// all character
func isInvalidSection(s string) bool {
	for _, c := range s {
		if bc := byte(c); bc == _WILDCARD || bc == _REMAINSALL {
			return true
		}
	}
	return false
}

// compile compile a url path to a clean path that replace all named variable
// to _WILDCARD or _REMAINSALL and extract all variable names
// if just want to match and don't need variable value, only use ':' or '*'
// for ':', it will catch the single section of url path seperated by '/'
// for '*', it will catch all remains url path, it should appear in the last
// of pattern for variables behind it will all be ignored
func compile(path string) (newPath string, names map[string]int, err error) {
	new := make([]byte, 0, len(newPath))
	sections := pathSections(path)
	fmt.Println(sections)
	nameIndex := 0
	for _, s := range sections {
		new = append(new, '/')
		if s != "" {
			switch s[0] {
			case ':':
				new = append(new, _WILDCARD)
				if name := s[1:]; len(name) > 0 {
					if isInvalidSection(name) {
						goto ERROR
					}
					if names == nil {
						names = make(map[string]int)
					}
					names[name] = nameIndex
				}
				nameIndex++
			case '*':
				new = append(new, _REMAINSALL)
				if name := s[1:]; len(name) > 0 {
					if isInvalidSection(name) {
						goto ERROR
					}
					if names == nil {
						names = make(map[string]int)
					}
					names[name] = nameIndex // -i means from section i to end
				}
				nameIndex++
				break // if read '*', other section is ignored
			default:
				new = append(new, []byte(s)...)
			}
		}
	}
	newPath = string(new)
	if names == nil {
		names = nilVars
	}
	return
ERROR:
	return "", nil,
		Errorf("path %s has pre-defined characters %c or %c",
			path, _WILDCARD, _REMAINSALL)
}

// newVarIndexer create a new VarIndexer with variable map and values
// if variables is empty then use default empty var indexer
func newVarIndexer(vars map[string]int, values []string) UrlVarIndexer {
	if len(vars) == 0 {
		return nilVarIndexer
	}
	v := new(pathVars)
	v.vars = vars
	v.values = values
	return v
}

// UrlVar return values of variable
func (v *pathVars) UrlVar(name string) string {
	if index, has := v.vars[name]; has {
		return v.values[index]
	}
	return ""
}

// UrlVars return all values of variable
func (v *pathVars) UrlVars() []string {
	return v.values
}

// ScanUrlVars scan values into variable addresses
// if address is nil, skip it
func (v *pathVars) ScanUrlVars(vars ...*string) {
	values := v.values
	l1, l2 := len(values), len(vars)
	for i := 0; i < l1 && i < l2; i++ {
		if vars[i] != nil {
			*vars[i] = values[i]
		}
	}
}

// init init handler and filters hold by routeProcessor
func (rp *routeProcessor) init(initHandler func(Handler) bool,
	initFilter func(Filter) bool,
	initWebSocketHandler func(WebSocketHandler) bool) bool {
	continu := true
	if rp.handlerProcessor != nil {
		continu = initHandler(rp.handlerProcessor.handler)
	}
	for _, f := range rp.filters {
		if !continu {
			break
		}
		continu = initFilter(f)
	}
	if continu && rp.wsHandlerProcessor != nil {
		continu = initWebSocketHandler(rp.wsHandlerProcessor.wsHandler)
	}
	return continu
}

// destroy destroy handler and filters hold by routeProcessor
func (rp *routeProcessor) destroy() {
	if rp.handlerProcessor != nil {
		rp.handlerProcessor.handler.Destroy()
	}
	for _, f := range rp.filters {
		f.Destroy()
	}
	if rp.wsHandlerProcessor != nil {
		rp.wsHandlerProcessor.wsHandler.Destroy()
	}
}

// NewRouter create a new Router
func NewRouter() Router {
	return new(router)
}

// Init init all handlers, filters, websocket handlers in route tree
func (rt *router) Init(initHandler func(Handler) bool,
	initFilter func(Filter) bool,
	initWebSocketHandler func(WebSocketHandler) bool) {
	rt.init(initHandler, initFilter, initWebSocketHandler)
}

// Init init all handlers, filters, websocket handlers in route tree
func (rt *router) init(initHandler func(Handler) bool,
	initFilter func(Filter) bool,
	initWebSocketHandler func(WebSocketHandler) bool) bool {
	continu := true
	if rt.processor != nil {
		continu = rt.processor.init(initHandler, initFilter, initWebSocketHandler)
	}
	for _, c := range rt.childs {
		if !continu {
			break
		}
		continu = c.init(initHandler, initFilter, initWebSocketHandler)
	}
	return continu
}

// Destroy destroy router and all handlers, filters, websocket handlers
func (rt *router) Destroy() {
	if rt.processor != nil {
		rt.processor.destroy()
	}
	for _, c := range rt.childs {
		c.Destroy()
	}
}

// routeProcessor return processor of current route node, if not exist
// then create a new one
func (rt *router) routeProcessor() *routeProcessor {
	if rt.processor == nil {
		rt.processor = new(routeProcessor)
	}
	return rt.processor
}

// AddFuncHandler add function handler to router for given pattern and method
func (rt *router) AddFuncHandler(pattern, method string, handler HandlerFunc) (err error) {
	if fHandler := strach.funcHandler(pattern); fHandler == nil {
		fHandler = new(funcHandler)
		if err = fHandler.setMethod(method, handler); err == nil {
			if err = rt.AddHandler(pattern, fHandler); err == nil {
				strach.setFuncHandler(pattern, fHandler)
			}
		}
	} else {
		err = fHandler.setMethod(method, handler)
	}
	return
}

// AddHandler add handler to router for given pattern
func (rt *router) AddHandler(pattern string, handler Handler) error {
	return rt.addPattern(pattern, func(rp *routeProcessor, pathVars map[string]int) error {
		if rp.handlerProcessor != nil {
			return Errorf("handler already exist")
		}
		rp.handlerProcessor = &handlerProcessor{vars: pathVars, handler: handler}
		return nil
	})
}

// AddFuncWebSocketHandler add funciton websocket handler to router for given pattern
func (rt *router) AddFuncWebSocketHandler(pattern string, handler WebSocketHandlerFunc) error {
	return rt.AddWebSocketHandler(pattern, handler)
}

// AddWebSocetHandler add websocket handler to router for given pattern
func (rt *router) AddWebSocketHandler(pattern string, handler WebSocketHandler) error {
	return rt.addPattern(pattern, func(rp *routeProcessor, pathVars map[string]int) error {
		if rp.wsHandlerProcessor != nil {
			return Err("websocket handler already exist")
		}
		rp.wsHandlerProcessor = &wsHandlerProcessor{vars: pathVars, wsHandler: handler}
		return nil
	})
}

// AddFuncFilter add function filter to router
func (rt *router) AddFuncFilter(pattern string, filter FilterFunc) error {
	return rt.AddFilter(pattern, filter)
}

// AddFuncFilter add filter to router
func (rt *router) AddFilter(pattern string, filter Filter) error {
	return rt.addPattern(pattern, func(rp *routeProcessor, _ map[string]int) error {
		rp.filters = append(rp.filters, filter)
		return nil
	})
}

// addPattern compile pattern, extract all variables, and add it to route tree
// setup by given function
func (rt *router) addPattern(pattern string, fn func(*routeProcessor, map[string]int) error) error {
	routePath, pathVars, err := compile(pattern)
	if err == nil {
		rt.addPath(routePath, func(n *router) {
			err = fn(n.routeProcessor(), pathVars)
		})
	}
	return err
}

// MatchWebSockethandler match url to find final websocket handler
func (rt *router) MatchWebSocketHandler(url *url.URL) (WebSocketHandler, UrlVarIndexer) {
	path := url.Path
	rt, values := rt.matchOne(path)
	if rt != nil {
		if processor := rt.processor; processor != nil {
			if wsp := processor.wsHandlerProcessor; wsp != nil {
				return wsp.wsHandler, newVarIndexer(wsp.vars, values)
			}
		}
	}
	return nil, nilVarIndexer
}

// MatchHandler match url to find final handler, last returned value will only
// be nil, it appeared only for compate with MatchHandlerFilters
func (rt *router) MatchHandler(url *url.URL) (Handler, UrlVarIndexer, []Filter) {
	path := url.Path
	rt, values := rt.matchOne(path)
	if rt != nil {
		if processor := rt.processor; processor != nil {
			if hp := processor.handlerProcessor; hp != nil {
				return hp.handler, newVarIndexer(hp.vars, values), nil
			}
		}
	}
	return nil, nilVarIndexer, nil
}

// MatchHandlerFilters match url to fin final handler and each filters
func (rt *router) MatchHandlerFilters(url *url.URL) (handler Handler,
	indexer UrlVarIndexer, filters []Filter) {
	var (
		pathIndex int
		values    []string
		continu   = true
		path      = url.Path
	)
	for continu {
		pathIndex, values, rt, continu = rt.matchMulti(path, pathIndex, values)
		if rt != nil {
			if rp := rt.processor; rp != nil {
				filters = append(filters, rp.filters...)
				if !continu {
					if hp := rp.handlerProcessor; hp != nil {
						handler, indexer = hp.handler, newVarIndexer(hp.vars, values)
					}
				}
			}
		}
	}
	if indexer == nil {
		indexer = nilVarIndexer
	}
	return
}

// addPath add an new path to route, use given function to operate the final
// route node for this path
func (rt *router) addPath(path string, fn func(*router)) {
	str := rt.str
	if str == "" {
		rt.str = path
	} else {
		diff, pathLen, strLen := 0, len(path), len(str)
		for diff != pathLen && diff != strLen && path[diff] == str[diff] {
			diff++
		}
		if diff < pathLen {
			first := path[diff]
			if diff == strLen {
				for i, c := range rt.chars {
					if c == first {
						rt.childs[i].addPath(path[diff:], fn)
						return
					}
				}
			} else { // diff < strLen
				rt.moveAllToChild(str[diff:], str[:diff])
			}
			newNode := &router{str: path[diff:]}
			rt.addChild(first, newNode)
			rt = newNode
		} else if diff < strLen {
			rt.moveAllToChild(str[diff:], path)
		}
	}
	fn(rt)
}

// moveAllToChild move all attributes to a new node, and make this new node
//  as one of it's child
func (rt *router) moveAllToChild(childStr string, newStr string) {
	rnCopy := &router{
		str:       childStr,
		chars:     rt.chars,
		childs:    rt.childs,
		processor: rt.processor,
	}
	rt.chars, rt.childs, rt.processor = nil, nil, nil
	rt.addChild(childStr[0], rnCopy)
	rt.str = newStr
}

// insertChild insert a child at given index, element since this index
// will all back an offset of 1 to make room for new element, the last
// element will be overrided it
func (rt *router) insertChild(index int, b byte, n *router) {
	chars, childs := rt.chars, rt.childs
	for i := len(chars) - 2; i >= index; i-- {
		chars[i+1], childs[i+1] = chars[i], childs[i]
	}
	chars[index], childs[index] = b, n
	rt.chars, rt.childs = chars, childs
}

// addChild add an child, all childs is sorted
func (rt *router) addChild(b byte, n *router) {
	chars, childs := rt.chars, rt.childs
	l := len(chars)
	if l == 0 {
		chars, childs = make([]byte, 0, reserveChildsCount),
			make([]*router, 0, reserveChildsCount)
	} else if cap := cap(chars); l == cap {
		cap += reserveChildsCount
		chars, childs = make([]byte, cap), make([]*router, cap)
		copy(chars, rt.chars)
		copy(childs, rt.childs)
	}
	chars, childs = chars[:l+1], childs[:l+1]
	var i int
	for i = l; i > 0 && chars[i] > b; i-- {
		chars[i], childs[i] = chars[i-1], childs[i-1]
	}
	chars[i], childs[i] = b, n
	rt.chars, rt.childs = chars, childs
}

// matchMulti match multi route node
// returned value:(first:next path start index, second:if continue, it's next node to match,
// else it's final match node, last:whether continu match)
func (rt *router) matchMulti(path string, pathIndex int, vars []string) (int,
	[]string, *router, bool) {
	str, strIndex := rt.str, 1
	pathIndex++
	strLen, pathLen := len(str), len(path)
	for strIndex != strLen {
		if pathIndex != pathLen {
			c := str[strIndex]
			strIndex++
			switch c {
			case path[pathIndex]: // else check character MatchPath or not
				pathIndex++
			case '*':
				// if read '*', MatchPath until next '/'
				start := pathIndex
				for pathIndex < pathLen && path[pathIndex] != '/' {
					pathIndex++
				}
				vars = append(vars, path[start:pathIndex])
			case '-': // parse end, full matched
				vars = append(vars, path[pathIndex:pathLen])
				pathIndex = pathLen
				strIndex = strLen
			default:
				return -1, nil, nil, false // not matched
			}
		} else {
			return -1, nil, nil, false // path parse end
		}
	}
	var continu bool
	if pathIndex != pathLen { // path not parse end, to find a child node to continue
		var node *router
		p := path[pathIndex]
		chars := rt.chars
		for i := range chars {
			if chars[i] == p {
				node = rt.childs[i] // child
				break
			}
		}
		rt = node
		if node != nil {
			continu = true
		}
	}
	return pathIndex, vars, rt, continu
}

// matchOne match one longest route node and return values of path variable
func (rt *router) matchOne(path string) (*router, []string) {
	var (
		str                string
		strIndex, strLen   int
		pathIndex, pathLen = 0, len(path)
		node               = rt
		values             = make([]string, 0, PathVarCount)
	)
	for node != nil {
		// skip first character, if it's root node, first char '/' can be safety skipped
		// otherwise, current node is selected by it's first character in parent node
		// so it also can be skipped
		str, strIndex = rt.str, 1
		pathIndex++
		strLen = len(str)
		for strIndex != strLen {
			if pathIndex != pathLen {
				c := str[strIndex]
				strIndex++
				switch c {
				case path[pathIndex]: // else check character MatchPath or not
					pathIndex++
				case '*':
					// if read '*', MatchPath until next '/'
					start := pathIndex
					for pathIndex < pathLen && path[pathIndex] != '/' {
						pathIndex++
					}
					values = append(values, path[start:pathIndex])
				case '-': // parse end, full matched
					values = append(values, path[pathIndex:pathLen])
					pathIndex = pathLen
					strIndex = strLen
				default:
					return nil, nil // not matched
				}
			} else {
				return nil, nil // path parse end
			}
		}
		node = nil
		if pathIndex != pathLen { // path not parse end, must find a child node to continue
			p := path[pathIndex]
			chars := rt.chars
			for i := range chars {
				if chars[i] == p {
					node = rt.childs[i] // child
					break
				}
			}
			rt = node // child to parse
		} /* else { path parse end, node is the last matched node }*/
	}
	return rt, values
}

// accessAllChilds access all childs of node
func (rt *router) accessAllChilds(fn func(*router) bool) {
	for _, n := range rt.childs {
		if !fn(n) {
			break
		}
	}
}

// PrintRouteTree print an route tree
// every level will be seperated by "-"
func PrintRouteTree(rt Router) {
	printRouteTree(rt.(*router), "")
}

// printRouteTree print route tree with given parent path
func printRouteTree(root *router, parentPath string) {
	if parentPath != "" {
		parentPath = parentPath + _PRINT_SEP
	}
	cur := parentPath + root.str
	fmt.Println(cur)
	root.accessAllChilds(func(n *router) bool {
		printRouteTree(n, cur)
		return true
	})
}
