package set

type SortedStrings struct {
	indexes map[string]int
	keys    []string
}

func NewSortedStrings() SortedStrings {
	return SortedStrings{
		indexes: make(map[string]int),
		keys:    make([]string, 0),
	}
}

func (s *SortedStrings) Put(key string) {
	if s.HasKey(key) {
		return
	}

	s.indexes[key] = len(s.keys)
	s.keys = append(s.keys, key)
}

func (s *SortedStrings) Remove(key string) {
	index, has := s.indexes[key]
	if !has {
		return
	}
	delete(s.indexes, key)

	for l := len(s.keys); index < l-1; index++ {
		k := s.keys[index+1]
		s.keys[index] = k
		s.indexes[k] = index
	}
	s.keys = s.keys[:index]
}

func (s *SortedStrings) HasKey(key string) bool {
	_, has := s.indexes[key]
	return has
}

func (s *SortedStrings) Keys() []string {
	return s.keys
}

func (s *SortedStrings) Clear() {
	for k := range s.indexes {
		delete(s.indexes, k)
	}

	s.keys = s.keys[:0]
}

type SortedInts struct {
	indexes map[int]int
	keys    []int
}

func (s *SortedInts) Put(key int) {
	if s.HasKey(key) {
		return
	}

	s.indexes[key] = len(s.keys)
	s.keys = append(s.keys, key)
}

func (s *SortedInts) Remove(key int) {
	index, has := s.indexes[key]
	if !has {
		return
	}
	delete(s.indexes, key)

	for l := len(s.keys); index < l-1; index++ {
		k := s.keys[index+1]
		s.keys[index] = k
		s.indexes[k] = index
	}
	s.keys = s.keys[:index]
}

func (s *SortedInts) HasKey(key int) bool {
	_, has := s.indexes[key]
	return has
}

func (s *SortedInts) Keys() []int {
	return s.keys
}

func NewSortedInts() SortedInts {
	return SortedInts{
		indexes: make(map[int]int),
		keys:    make([]int, 0),
	}
}

func (s *SortedInts) Clear() {
	for k := range s.indexes {
		delete(s.indexes, k)
	}

	s.keys = s.keys[:0]
}
