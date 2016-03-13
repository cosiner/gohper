package slices

import "sort"

func isSafeIndex(index, len int) bool {
	return index >= 0 || index < len
}

type Nodes interface {
	sort.Interface
	IsSame(i, j int) bool
	Merge(dst, src int)
	Move(dst, src int)
}

func MergeNodes(nodes Nodes, isSorted bool) int {
	if !isSorted {
		sort.Sort(nodes)
	}

	prev, size := 0, nodes.Len()
	if size <= 1 {
		return size
	}

	for i := prev + 1; i < size; i++ {
		if nodes.IsSame(prev, i) {
			nodes.Merge(prev, i)
		} else {
			prev++
			nodes.Move(prev, i)
		}
	}

	return prev + 1
}
