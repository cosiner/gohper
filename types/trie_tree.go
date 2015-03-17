// the TrieTree is port from package zerver's router
// but simplified the implementation
package types

import (
	"fmt"

	"io"
)

type TrieTree struct {
	str        string
	childChars []byte
	childs     []*TrieTree
	Value      interface{}
}

func (tt *TrieTree) AddPath(path string, value interface{}) {
	str := tt.str
	if str == "" && len(tt.childChars) == 0 {
		tt.str = path
	} else {
		diff, pathLen, strLen := 0, len(path), len(str)
		for diff != pathLen && diff != strLen && path[diff] == str[diff] {
			diff++
		}
		if diff < pathLen {
			first := path[diff]
			if diff == strLen {
				for i, c := range tt.childChars {
					if c == first {
						tt.childs[i].AddPath(path[diff:], value)
					}
				}
			} else { // diff < strLen
				tt.moveAllToChild(str[diff:], str[:diff])
			}
			newNode := &TrieTree{str: path[diff:]}
			tt.addChild(first, newNode)
			tt = newNode
		} else if diff < strLen {
			tt.moveAllToChild(str[diff:], path)
		}
	}
	tt.Value = value
}

// moveAllToChild move all attributes to a new node, and make this new node
//  as one of it's child
func (tt *TrieTree) moveAllToChild(childStr string, newStr string) {
	rnCopy := &TrieTree{
		str:        childStr,
		childChars: tt.childChars,
		childs:     tt.childs,
		Value:      tt.Value,
	}
	tt.childChars, tt.childs, tt.Value = nil, nil, nil
	tt.addChild(childStr[0], rnCopy)
	tt.str = newStr
}

// addChild add an child, all childs is sorted
func (tt *TrieTree) addChild(b byte, n *TrieTree) {
	childChars, childs := tt.childChars, tt.childs
	l := len(childChars)
	childChars, childs = make([]byte, l+1), make([]*TrieTree, l+1)
	copy(childChars, tt.childChars)
	copy(childs, tt.childs)
	for ; l > 0 && childChars[l-1] > b; l-- {
		childChars[l], childs[l] = childChars[l-1], childs[l-1]
	}
	childChars[l], childs[l] = b, n
	tt.childChars, tt.childs = childChars, childs
}

// Match match one longest route node and return values of path variable
func (tt *TrieTree) Match(path string) interface{} {
	var (
		str                string
		strLen             int
		pathIndex, pathLen = 0, len(path)
		node               = tt
		start              int
	)
	for node != nil {
		str = tt.str
		strLen = len(str)
		pathIndex += start
		for i := start; i < strLen; i++ {
			if pathIndex == pathLen || str[i] != path[pathIndex] {
				return nil
			}
			pathIndex++
		}
		node = nil
		if pathIndex != pathLen { // path not parse end, must find a child node to continue
			p := path[pathIndex]
			for i, c := range tt.childChars {
				if c == p {
					node = tt.childs[i] // child
					break
				}
			}
			if node == nil {
				return nil
			}
			tt = node // child to parse
		} /* else { path parse end, node is the last matched node }*/
		start = 1
	}
	return tt.Value
}

// Print print an route tree
// every level will be seperated by "-"
func (tt *TrieTree) Print(w io.Writer) {
	tt.print(w, "")
}

// print print route tree with given parent path
func (tt *TrieTree) print(w io.Writer, parentPath string) {
	if parentPath != "" {
		parentPath = parentPath + "-"
	}
	cur := parentPath + tt.str
	fmt.Fprintln(w, cur)
	for _, n := range tt.childs {
		n.print(w, cur)
	}
}
