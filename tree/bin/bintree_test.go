package bin

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestBinaryTree(t *testing.T) {
	tt := testing2.Wrap(t)

	bt := BinTree{}
	bt.Add(1, 'A', false)
	bt.Add(2, 'B', false)
	bt.Add(3, 'C', false)
	bt.Add(4, 'D', false)
	bt.Add(5, 'E', false)
	bt.Add(6, 'F', false)
	bt.Add(7, 'G', false)

	tt.Eq('A', bt.Search(1).(int32))
	tt.Eq('B', bt.Search(2).(int32))
	tt.Eq('C', bt.Search(3).(int32))
	tt.Eq('D', bt.Search(4).(int32))
	tt.Eq('E', bt.Search(5).(int32))
	tt.Eq('F', bt.Search(6).(int32))
	tt.Eq('G', bt.Search(7).(int32))

	tt.Nil(bt.Search(8))
	tt.Nil(bt.Search(9))

	bt.Add(7, 'H', true)
	tt.Eq('H', bt.Search(7).(int32))

	bt.Add(7, 'I', false)
	tt.Eq('H', bt.Search(7).(int32))
}
