package server

import (
	"fmt"
	"strings"

	. "github.com/cosiner/golib/errors"
)

const (
	wildcard   = '|' // MUST BE:other character < wildcard < remainsAll
	remainsAll = '~'
)

// pathSections divide path by '/', trim end '/', the first '/' will not be processed
func pathSections(path string) []string {
	if l := len(path); l > 0 {
		if path[l-1] == '/' {
			l = l - 1
		}
		path = path[0:l] // trim first and last '/'
	}
	return strings.Split(path, "/")
}

// isInvalidSection check whether section has the predefined wildcard and match
// all character
func isInvalidSection(s string) bool {
	for _, c := range s {
		if bc := byte(c); bc == wildcard || bc == remainsAll {
			return true
		}
	}
	return false
}

// compile compile a url path to a clean path that replace all named variable
// to '*' and extract all variable names
// only format like "/:name/" is supported
func compile(path string) (newPath string, names map[string]int, err error) {
	new := make([]byte, 0, len(newPath))
	new[0] = '/'
	names = make(map[string]int)
	sections := pathSections(path)
	for i, s := range sections {
		new = append(new, '/')
		if s != "" {
			switch s[0] {
			case ':':
				new = append(new, wildcard)
				if name := s[1:]; len(name) > 0 {
					if isInvalidSection(name) {
						goto ERROR
					}
					names[name] = i
				}
			case '*':
				new = append(new, remainsAll)
				if name := s[1:]; len(name) > 0 {
					if isInvalidSection(name) {
						goto ERROR
					}
					names[name] = -i // -i means from section i to end
				}
				break // if read '*', other section is ignored
			default:
				new = append(new, []byte(s)...)
			}
		}
	}
	newPath = string(new)
	return
ERROR:
	return "", nil,
		Errorf("path %s has pre-defined characters %c or %c",
			path, wildcard, remainsAll)
}

// reserveChildsCount is route childs slice increment and init size
var reserveChildsCount = 1

// insertChild insert a child at given index, element since this index
// will all back an offset of 1 to make room for new element, the last
// element will be overrided
func (rn *routeNode) insertChild(index int, b byte, n *routeNode) {
	chars, childs := rn.chars, rn.childs
	for i := len(chars) - 2; i >= index; i-- {
		chars[i+1], childs[i+1] = chars[i], childs[i]
	}
	chars[index], childs[index] = b, n
	rn.chars, rn.childs = chars, childs
}

// addChild add an child, all childs is sorted
func (rn *routeNode) addChild(b byte, n *routeNode) {
	chars, childs := rn.chars, rn.childs
	l := len(chars)
	if l == 0 {
		chars, childs = make([]byte, 0, reserveChildsCount),
			make([]*routeNode, 0, reserveChildsCount)
	} else if cap := cap(chars); l == cap {
		cap += reserveChildsCount
		chars, childs = make([]byte, cap), make([]*routeNode, cap)
		copy(chars, rn.chars)
		copy(childs, rn.childs)
	}
	chars, childs = chars[:l+1], childs[:l+1]
	var i int
	for i = l; i > 0 && chars[i] > b; i-- {
		chars[i], childs[i] = chars[i-1], childs[i-1]
	}
	chars[i], childs[i] = b, n
	rn.chars, rn.childs = chars, childs
}

type handlerRoute struct {
	routeVar map[string]int
	handler  Handler
}

type wsHandlerRoute struct {
	routeVar  map[string]int
	wsHandler WebSocketHandler
}

type routeHandler struct {
	handler   handlerRoute
	filters   []Filter
	wsHandler wsHandlerRoute
}

type routeNode struct {
	str     string
	chars   []byte
	childs  []*routeNode
	handler *routeHandler
}

func (rn *routeNode) routeHandler() *routeHandler {
	if rn.handler == nil {
		rn.handler = new(routeHandler)
	}
	return rn.handler
}

func (rn *routeNode) AddHandler(pattern string, handler Handler) error {
	return rn.addPattern(pattern, func(rh *routeHandler, rv *map[string]int) error {
		if rh.handler != nil {
			return Errorf("handler of path %s already exist", routePath)
		}
		rh.handler = handlerRoute{routeVars: routeVars, handler: handler}
		return nil
	})
}

func (rn *routeNode) AddFuncWebSocketHandler(pattern string, handler WebSocketHandlerFunc) error {
	return rn.AddWebSocketHandler(pattern, handler)
}

func (rn *routeNode) AddWebSocketHandler(pattern string, handler WebSocketHandler) error {
	return rn.addPattern(pattern, func(rh *routeHandler, rv *map[string]int) error {
		if rh.wsHandler != nil {
			return Errorf("websocket handler of path %s already exist", routePath)
		}
		rh.wsHandler = handlerRoute{routeVars: routeVars, handler: handler}
		return nil
	})
}

func (rn *routeNode) AddFuncFilter(pattern string, filter FilterFunc) error {
	return rn.AddFilter(pattern, filter)
}

func (rn *routeNode) AddFilter(pattern string, filter Filter) error {
	return rn.addPattern(pattern, func(rh *routeHandler, _ map[string]int) error {
		rh.filters = append(rh.filters, filter)
	})
}

func (rn *routeNode) isHandlerNode() bool {
	return rn.handler != nil
}

// moveToChild move all attributes to a new node, and make new node as one of it's
// child
func (rn *routeNode) moveToChild(childStr string, newStr string) {
	rnCopy := &routeNode{
		str:     childStr,
		chars:   rn.chars,
		childs:  rn.childs,
		handler: rn.handler,
	}
	rn.chars, rn.childs, rn.handler = nil, nil, nil
	rn.addChild(childStr[0], rnCopy)
	rn.str = newStr
}

// addPattern compile pattern, and add it to route tree
func (rn *routeNode) addPattern(pattern string, fn func(*routeHandler, string, map[string]int) error) error {
	routePath, routeVars, err := compile(pattern)
	if err == nil {
		rn.addPath(routePath, func(n *routeNode) {
			err = fn(n.routeHandler(), routeVars)
		})
	}
	return err
}

// addPath add an new path to route
func (rn *routeNode) addPath(path string, fn func(*routeNode)) {
	str := rn.str
	diff, pathLen, strLen := 0, len(path), len(str)
	for diff != pathLen && diff != strLen && path[diff] == str[diff] {
		diff++
	}
	if diff < pathLen {
		first := path[diff]
		if diff == strLen {
			for i, c := range rn.chars {
				if c == first {
					rn.childs[i].addPath(path[diff:], fn)
					return
				}
			}
		} else { // diff < strLen
			rn.moveToChild(str[diff:], str[:diff])
		}
		newNode := &routeNode{str: path[diff:]}
		rn.addChild(first, newNode)
		rn = newNode
	} else if diff < strLen {
		rn.moveToChild(str[diff:], path)
	}
	fn(rn)
}

// matchMulti match multi route node
// returned value:(first:next path start index, second:if continue, it's next node to match,
// else it's final match node, last:whether continu match)
func (rn *routeNode) matchMulti(path string, pathIndex int) (int, *routeNode, bool) {
	str, strIndex := rn.str, 1
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
				for pathIndex < pathLen && path[pathIndex] != '/' {
					pathIndex++
				}
			case '-': // parse end, full matched
				pathIndex = pathLen
				strIndex = strLen
			default:
				return -1, nil, false // not matched
			}
		} else {
			return -1, nil, false // path parse end
		}
	}
	var continu bool
	if pathIndex != pathLen { // path not parse end, to find a child node to continue
		var node *routeNode
		p := path[pathIndex]
		for i, c := range rn.chars {
			if c == p {
				node = rn.childs[i] // child
				break
			}
		}
		if node != nil {
			continu = true
			rn = node
		}
	}
	return pathIndex, rn, continu
}

// matchOne match one longest route node
func (rn *routeNode) matchOne(path string, pathIndex int) *routeNode {
	var (
		str              string
		strIndex, strLen int
		pathLen                     = len(path)
		node             *routeNode = rn
	)
	for rn != nil {
		str, strIndex = rn.str, 1
		strLen = len(str)
		pathIndex++
		for strIndex != strLen {
			if pathIndex != pathLen {
				c := str[strIndex]
				strIndex++
				switch c {
				case path[pathIndex]: // else check character MatchPath or not
					pathIndex++
				case '*':
					// if read '*', MatchPath until next '/'
					for pathIndex < pathLen && path[pathIndex] != '/' {
						pathIndex++
					}
				case '-': // parse end, full matched
					pathIndex = pathLen
					strIndex = strLen
				default:
					return nil // not matched
				}
			} else {
				return nil // path parse end
			}
		}
		rn = nil
		if pathIndex != pathLen { // path not parse end, must find a child node to continue
			p := path[pathIndex]
			for i, c := range node.chars {
				if c == p {
					rn = node.childs[i] // child
					break
				}
			}
			node = rn // child to parse
		} /* else { path parse end, node is the last matched node }*/
	}
	return node
}

// PrintRouteTree print an route tree
func PrintRouteTree(root *routeNode) {
	printRouteTree(root, "")
}

func printRouteTree(root *routeNode, parentPath string) {
	if parentPath != "" {
		parentPath = parentPath + "-"
	}
	cur := parentPath + root.str
	fmt.Println(cur)
	root.accessAllChilds(func(n *routeNode) bool {
		printRouteTree(n, cur)
		return true
	})
}

// accessAllChilds access all childs of node
func (rn *routeNode) accessAllChilds(fn func(*routeNode) bool) {
	for _, n := range rn.childs {
		if !fn(n) {
			break
		}
	}
}
