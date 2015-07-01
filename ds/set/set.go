package set

var exist = struct{}{}

type Strings map[string]struct{}

func (s Strings) Put(key string) {
	s[key] = exist
}

func (s Strings) Remove(key string) {
	delete(s, key)
}

func (s Strings) HasKey(key string) bool {
	_, has := s[key]
	return has
}

func (s Strings) Keys() []string {
	keys := make([]string, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i++
	}

	return keys
}

type Ints map[int]struct{}

func (s Ints) Put(key int) {
	s[key] = exist
}

func (s Ints) Remove(key int) {
	delete(s, key)
}

func (s Ints) HasKey(key int) bool {
	_, has := s[key]
	return has
}

func (s Ints) Keys() []int {
	keys := make([]int, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i++
	}

	return keys
}
