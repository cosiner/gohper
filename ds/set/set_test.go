package set

import (
	"sort"
	"testing"

	"github.com/cosiner/gohper/testing2"
)

var stringKeys = []string{
	"A",
	"B",
	"C",
	"D",
	"E",
	"F",
	"A",
}
var intKeys = []int{
	'A',
	'B',
	'C',
	'D',
	'E',
	'F',
	'A',
}

func TestStrings(t *testing.T) {
	tt := testing2.Wrap(t)

	strings := make(Strings)

	for _, k := range stringKeys {
		strings.Put(k)
	}

	tt.Eq(len(stringKeys)-1, len(strings))

	for _, k := range stringKeys {
		tt.True(strings.HasKey(k))
	}

	strings.Remove("A")

	keys := strings.Keys()
	sort.Strings(keys)
	tt.DeepEq(stringKeys[1:len(stringKeys)-1], keys)

}

func TestInts(t *testing.T) {
	tt := testing2.Wrap(t)
	ints := make(Ints)

	for _, k := range intKeys {
		ints.Put(k)
	}

	tt.Eq(len(intKeys)-1, len(ints))

	for _, k := range intKeys {
		tt.True(ints.HasKey(k))
	}

	ints.Remove('A')
	keys := ints.Keys()
	sort.Ints(keys)
	tt.DeepEq(intKeys[1:len(intKeys)-1], keys)
}

func TestSortedStrings(t *testing.T) {
	tt := testing2.Wrap(t)

	strings := NewSortedStrings()
	for _, k := range stringKeys {
		strings.Put(k)
	}
	tt.Eq(len(intKeys)-1, len(strings.Keys()))
	strings.Remove("G")
	strings.Remove("A")
	tt.DeepEq(stringKeys[1:len(stringKeys)-1], strings.Keys())
}

func TestSortedInts(t *testing.T) {
	tt := testing2.Wrap(t)

	ints := NewSortedInts()
	for _, k := range intKeys {
		ints.Put(k)
	}
	tt.Eq(len(intKeys)-1, len(ints.Keys()))
	ints.Remove('G')
	ints.Remove('A')
	tt.DeepEq(intKeys[1:len(intKeys)-1], ints.Keys())
}
