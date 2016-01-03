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

func (s Strings) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s Strings) Size() int {
	return len(s)
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

func (s Ints) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s Ints) Size() int {
	return len(s)
}

type Bytes map[byte]struct{}

func (s Bytes) Put(key byte) {
	s[key] = exist
}

func (s Bytes) Remove(key byte) {
	delete(s, key)
}

func (s Bytes) HasKey(key byte) bool {
	_, has := s[key]
	return has
}

func (s Bytes) Keys() []byte {
	keys := make([]byte, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i++
	}

	return keys
}

func (s Bytes) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s Bytes) Size() int {
	return len(s)
}