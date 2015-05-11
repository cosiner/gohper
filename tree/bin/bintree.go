// Package bin implements a simple binary-search tree
package bin

type BinTree struct {
	score int
	value interface{}

	left, right *BinTree
}

func (t *BinTree) Search(score int) interface{} {
	root := t
	for root != nil {
		switch {
		case root.score == score:
			return root.value
		case root.score > score:
			root = root.left
		case root.score < score:
			root = root.right
		}
	}
	return nil
}

func (t *BinTree) Add(score int, value interface{}, replace bool) {
	root := t
	for {
		switch {
		case root.score == score:
			if replace {
				root.value = value
			}
			return
		case root.score > score:
			if root.left == nil {
				root.left = &BinTree{score: score, value: value}
				return
			} else {
				root = root.left
			}
		case root.score < score:
			if root.right == nil {
				root.right = &BinTree{score: score, value: value}
				return
			} else {
				root = root.right
			}
		}
	}
}
