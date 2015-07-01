package tree

import (
	"fmt"
	"io"
)

type Trie struct {
	Str        string
	ChildChars []byte
	Childs     []*Trie
	Value      interface{}
}

func (t *Trie) AddPath(path string, value interface{}) {
	t.AddPathFor(path, func(t *Trie) {
		if value != nil {
			t.Value = value
		}
	})
}

func (t *Trie) AddPathFor(path string, fn func(*Trie)) {
	if !t.HasElement() {
		t.Str = path
		fn(t)
		return
	}

	str := t.Str
	diff, pathLen, strLen := 0, len(path), len(str)
	for diff != pathLen && diff != strLen && path[diff] == str[diff] {
		diff++
	}

	if diff < pathLen {
		first := path[diff]
		if diff == strLen {
			for i, c := range t.ChildChars {
				if c == first {
					t.Childs[i].AddPathFor(path[diff:], fn)

					return
				}
			}
		} else { // diff < strLen
			t.moveAllToChild(str[diff:], str[:diff])
		}

		newNode := &Trie{Str: path[diff:]}
		t.addChild(first, newNode)
		t = newNode
	} else if diff < strLen {
		t.moveAllToChild(str[diff:], path)
	}

	fn(t)
}

// moveAllToChild move all atreeributes to a new node, and make this new node
//  as one of it's child
func (t *Trie) moveAllToChild(childStr string, newStr string) {
	rnCopy := &Trie{
		Str:        childStr,
		ChildChars: t.ChildChars,
		Childs:     t.Childs,
		Value:      t.Value,
	}

	t.ChildChars, t.Childs, t.Value = nil, nil, nil
	t.addChild(childStr[0], rnCopy)
	t.Str = newStr
}

// addChild add an child, all Childs is sorted
func (t *Trie) addChild(b byte, n *Trie) {
	chars, childs := t.ChildChars, t.Childs
	l := len(chars)
	chars, childs = make([]byte, l+1), make([]*Trie, l+1)
	copy(chars, t.ChildChars)
	copy(childs, t.Childs)

	for ; l > 0 && chars[l-1] > b; l-- {
		chars[l], childs[l] = chars[l-1], childs[l-1]
	}
	chars[l], childs[l] = b, n
	t.ChildChars, t.Childs = chars, childs
}

const (
	TRIE_NO     = iota // TRIE_NO means there is a chaacter don't match
	TRIE_PREFIX        // TRIE_PREFIX means last node's `Str` is only match the begining part
	TRIE_FULL          // TRIE_FULL means last node's `Str` is full matched
)

func (t *Trie) MatchFrom(nodestart int, path string) (tr *Trie, index int, typ int) {
	var (
		str                string
		strLen             int
		pathIndex, pathLen     = 0, len(path)
		node                   = t
		start              int = nodestart
	)

	for node != nil {
		str = t.Str
		strLen = len(str)
		pathIndex += start
		for i := start; i < strLen; i++ {
			if pathIndex == pathLen {
				return t, i, TRIE_PREFIX
			} else if str[i] != path[pathIndex] {
				return nil, 0, TRIE_NO
			}

			pathIndex++
		}

		node = nil
		if pathIndex != pathLen { // path not parse end, must find a child node to continue
			p := path[pathIndex]
			for i, c := range t.ChildChars {
				if c == p {
					node = t.Childs[i] // child
					break
				}
			}

			if node == nil {
				return nil, 0, TRIE_NO
			}

			t = node // child to parse
			start = 1
		} /* else { path parse end, node is the last matched node }*/
	}

	return t, strLen, TRIE_FULL
}

func (t *Trie) Match(path string) (*Trie, int, int) {
	return t.MatchFrom(0, path)
}

// Match one longest route node and return values of path variable
func (t *Trie) MatchValue(path string) interface{} {
	t, _, m := t.Match(path)
	if m != TRIE_FULL {
		return nil
	}

	return t.Value
}

func (t *Trie) prefixMatch(parent *Trie, path string) *Trie {
	var (
		pathIndex, pathLen = 0, len(path)
	)

	str := t.Str
	strLen := len(str)
	for i := 0; i < strLen; i++ {
		if pathIndex == pathLen || str[i] != path[pathIndex] {
			return parent
		}

		pathIndex++
	}

	if pathIndex != pathLen {
		p := path[pathIndex]
		for i, c := range t.ChildChars {
			if c == p {
				return t.Childs[i].prefixMatch(t, path[pathIndex:])
			}
		}
	}

	return t
}

// PrefixMatchValue assumes each node as a prefix, it will match the longest prefix
// and return it's node value or nil
func (t *Trie) PrefixMatchValue(path string) interface{} {
	if t = t.prefixMatch(nil, path); t == nil {
		return nil
	}

	return t.Value
}

func (t *Trie) HasElement() bool {
	return t.Str != "" || len(t.Childs) != 0
}

func NopHook(interface{}) string {
	return ""
}

func (t *Trie) Print(w io.Writer, withCurr bool, parentPath, sep string, hook func(value interface{}) string) {
	if parentPath != "" {
		parentPath = parentPath + sep
	}

	if withCurr {
		parentPath += t.Str
		if t.Value != nil {
			fmt.Fprintln(w, parentPath+hook(t.Value))
		}
	}

	for _, n := range t.Childs {
		n.Print(w, true, parentPath, sep, hook)
	}
}

func (t *Trie) Visit(visitor func(string, interface{})) {
	if visitor != nil {
		t.visit("", visitor)
	}
}

func (t *Trie) visit(parentPath string, visitor func(string, interface{})) {
	path := parentPath + t.Str
	if t.Value != nil {
		visitor(path, t.Value)
	}

	for _, c := range t.Childs {
		c.visit(path, visitor)
	}
}
