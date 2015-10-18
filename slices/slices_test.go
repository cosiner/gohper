package slices

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestStrings(t *testing.T) {
	tt := testing2.Wrap(t)

	slice := Strings{}
	strings := Strings{"1", "2", "3", "4", "5", "6", "7", "8"}
	for _, s := range strings {
		slice = slice.IncrAppend(s)
		tt.Eq(len(slice), cap(slice))
	}
	slice = slice.FitCapToLen()
	tt.Eq(len(slice), cap(slice))

	slice = append(slice, "9", "10")
	tt.NE(len(slice), cap(slice))
	slice = slice.FitCapToLen()
	tt.Eq(len(slice), cap(slice))

	slice = slice.Remove(0)
	tt.DeepEq(Strings{"2", "3", "4", "5", "6", "7", "8", "9", "10"}, slice)
	slice = slice.Remove(slice.Len() - 1)
	tt.DeepEq(Strings{"2", "3", "4", "5", "6", "7", "8", "9"}, slice.Append("2").RmDups())

	tt.Eq("23456789", slice.Join("", ""))
	tt.Eq("2=?, 3=?, 4=?, 5=?, 6=?, 7=?, 8=?, 9=?", slice.Join("=?", ", "))
}

type IdNode struct {
	Id  string
	Num int
}

type IdNodes []IdNode

func (nodes IdNodes) Len() int {
	return len(nodes)
}

func (nodes IdNodes) Swap(i, j int) {
	nodes[i].Id, nodes[j].Id = nodes[j].Id, nodes[i].Id
	nodes[i].Num, nodes[j].Num = nodes[j].Num, nodes[i].Num
}

func (nodes IdNodes) Less(i, j int) bool {
	return nodes[i].Id < nodes[j].Id
}
func (nodes IdNodes) IsSame(i, j int) bool {
	return nodes[i].Id == nodes[j].Id
}

func (nodes IdNodes) Merge(dst, src int) {
	nodes[dst].Num += nodes[src].Num
}
func (nodes IdNodes) Move(dst, src int) {
	nodes[dst].Id = nodes[src].Id
	nodes[dst].Num = nodes[src].Num
}

func TestMergeNode(t *testing.T) {
	tt := testing2.Wrap(t)
	nodes := []IdNode{
		IdNode{"1", 1},
		IdNode{"2", 2},
		IdNode{"3", 3},
		IdNode{"1", 1},
		IdNode{"2", 2},
		IdNode{"3", 3},
		IdNode{"1", 1},
		IdNode{"2", 2},
		IdNode{"2", -1},
	}
	m := make(map[string]int)
	for _, n := range nodes {
		m[n.Id] += n.Num
	}
	len := MergeNodes(IdNodes(nodes), false)
	nodes = nodes[:len]

	for _, n := range nodes {
		tt.Eq(m[n.Id], n.Num)
	}
}
