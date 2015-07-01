package slices

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestStrings(t *testing.T) {
	tt := testing2.Wrap(t)

	slice := []string{}
	strings := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	for _, s := range strings {
		slice = IncrAppendString(slice, s)
		tt.Eq(len(slice), cap(slice))
	}
	slice = FitCapToLenString(slice)
	tt.Eq(len(slice), cap(slice))

	slice = append(slice, "9", "10")
	tt.NE(len(slice), cap(slice))
	slice = FitCapToLenString(slice)
	tt.Eq(len(slice), cap(slice))
}

func TestInts(t *testing.T) {
	tt := testing2.Wrap(t)

	slice := []int{}
	strings := []int{1, 2, 3, 4, 5, 6, 7, 8}
	for _, s := range strings {
		slice = IncrAppendInt(slice, s)
		tt.Eq(len(slice), cap(slice))
	}

	slice = FitCapToLenInt(slice)
	tt.Eq(len(slice), cap(slice))
	slice = append(slice, 9, 10)
	tt.NE(len(slice), cap(slice))
	slice = FitCapToLenInt(slice)
	tt.Eq(len(slice), cap(slice))
}

func TestUints(t *testing.T) {
	tt := testing2.Wrap(t)

	slice := []uint{}
	strings := []uint{1, 2, 3, 4, 5, 6, 7, 8}
	for _, s := range strings {
		slice = IncrAppendUint(slice, s)
		tt.Eq(len(slice), cap(slice))
	}

	slice = FitCapToLenUint(slice)
	tt.Eq(len(slice), cap(slice))
	slice = append(slice, 9, 10)
	tt.NE(len(slice), cap(slice))
	slice = FitCapToLenUint(slice)
	tt.Eq(len(slice), cap(slice))
}
