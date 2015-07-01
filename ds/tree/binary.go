package tree

type Binary struct {
	score int
	value interface{}

	left, right *Binary
}

func (t *Binary) Search(score int) interface{} {
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

func (t *Binary) Add(score int, value interface{}, replace bool) {
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
				root.left = &Binary{score: score, value: value}

				return
			} else {
				root = root.left
			}
		case root.score < score:
			if root.right == nil {
				root.right = &Binary{score: score, value: value}

				return
			} else {
				root = root.right
			}
		}
	}
}
