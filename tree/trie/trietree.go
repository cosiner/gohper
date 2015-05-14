package trie

import (
	"fmt"
	"io"
)

type TrieTree struct {
	Str        string
	ChildChars []byte
	Childs     []*TrieTree
	Value      interface{}
}

func (tt *TrieTree) AddPath(path string, value interface{}) {
	tt.AddPathFor(path, func(t *TrieTree) {
		if value != nil {
			t.Value = value
		}
	})
}

func (tt *TrieTree) AddPathFor(path string, fn func(*TrieTree)) {
	if !tt.HasElement() {
		tt.Str = path
		fn(tt)
		return
	}

	str := tt.Str
	diff, pathLen, strLen := 0, len(path), len(str)
	for diff != pathLen && diff != strLen && path[diff] == str[diff] {
		diff++
	}

	if diff < pathLen {
		first := path[diff]
		if diff == strLen {
			for i, c := range tt.ChildChars {
				if c == first {
					tt.Childs[i].AddPathFor(path[diff:], fn)

					return
				}
			}
		} else { // diff < strLen
			tt.moveAllToChild(str[diff:], str[:diff])
		}

		newNode := &TrieTree{Str: path[diff:]}
		tt.addChild(first, newNode)
		tt = newNode
	} else if diff < strLen {
		tt.moveAllToChild(str[diff:], path)
	}

	fn(tt)
}

// moveAllToChild move all attributes to a new node, and make this new node
//  as one of it's child
func (tt *TrieTree) moveAllToChild(childStr string, newStr string) {
	rnCopy := &TrieTree{
		Str:        childStr,
		ChildChars: tt.ChildChars,
		Childs:     tt.Childs,
		Value:      tt.Value,
	}

	tt.ChildChars, tt.Childs, tt.Value = nil, nil, nil
	tt.addChild(childStr[0], rnCopy)
	tt.Str = newStr
}

// addChild add an child, all Childs is sorted
func (tt *TrieTree) addChild(b byte, n *TrieTree) {
	ChildChars, Childs := tt.ChildChars, tt.Childs
	l := len(ChildChars)
	ChildChars, Childs = make([]byte, l+1), make([]*TrieTree, l+1)
	copy(ChildChars, tt.ChildChars)
	copy(Childs, tt.Childs)

	for ; l > 0 && ChildChars[l-1] > b; l-- {
		ChildChars[l], Childs[l] = ChildChars[l-1], Childs[l-1]
	}
	ChildChars[l], Childs[l] = b, n
	tt.ChildChars, tt.Childs = ChildChars, Childs
}

const (
	NO     = iota // NO means there is a chaacter don't match
	PREFIX        // PREFIX means last node's `Str` is only match the begining part
	FULL          // FULL means last node's `Str` is full matched
)

func (tt *TrieTree) MatchFrom(nodestart int, path string) (t *TrieTree, index int, typ int) {
	var (
		str                string
		strLen             int
		pathIndex, pathLen     = 0, len(path)
		node                   = tt
		start              int = nodestart
	)

	for node != nil {
		str = tt.Str
		strLen = len(str)
		pathIndex += start
		for i := start; i < strLen; i++ {
			if pathIndex == pathLen {
				return tt, i, PREFIX
			} else if str[i] != path[pathIndex] {
				return nil, 0, NO
			}

			pathIndex++
		}

		node = nil
		if pathIndex != pathLen { // path not parse end, must find a child node to continue
			p := path[pathIndex]
			for i, c := range tt.ChildChars {
				if c == p {
					node = tt.Childs[i] // child
					break
				}
			}

			if node == nil {
				return nil, 0, NO
			}

			tt = node // child to parse
			start = 1
		} /* else { path parse end, node is the last matched node }*/
	}

	return tt, strLen, FULL
}

func (tt *TrieTree) Match(path string) (t *TrieTree, index int, typ int) {
	return tt.MatchFrom(0, path)
}

// Match one longest route node and return values of path variable
func (tt *TrieTree) MatchValue(path string) interface{} {
	t, _, m := tt.Match(path)
	if m != FULL {
		return nil
	}

	return t.Value
}

func (tt *TrieTree) HasElement() bool {
	return tt.Str != "" || len(tt.Childs) != 0
}

func NopHook(interface{}) string {
	return ""
}

func (tt *TrieTree) Print(w io.Writer, withCurr bool, parentPath, sep string, hook func(value interface{}) string) {
	if parentPath != "" {
		parentPath = parentPath + sep
	}

	if withCurr {
		parentPath += tt.Str
		if tt.Value != nil {
			fmt.Fprintln(w, parentPath+hook(tt.Value))
		}
	}

	for _, n := range tt.Childs {
		n.Print(w, true, parentPath, sep, hook)
	}
}

func (tt *TrieTree) Visit(visitor func(string, interface{})) {
	if visitor != nil {
		tt.visit("", visitor)
	}
}

func (tt *TrieTree) visit(parentPath string, visitor func(string, interface{})) {
	path := parentPath + tt.Str
	if tt.Value != nil {
		visitor(path, tt.Value)
	}

	for _, c := range tt.Childs {
		c.visit(path, visitor)
	}
}
