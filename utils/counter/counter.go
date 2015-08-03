package counter

type Counter map[string]int

func New() Counter {
	return make(Counter)
}

func (c Counter) Add(key string) int {
	cnt := c[key]
	c[key] = cnt + 1
	return cnt + 1
}

func (c Counter) Remove(key string) int {
	val, has := c[key]
	if !has {
		return 0
	}

	if val <= 1 {
		delete(c, key)
		return 0
	}

	c[key] = val - 1
	return val - 1
}

func (c Counter) Clear(key string) int {
	cnt := c[key]
	delete(c, key)
	return cnt
}

func (c Counter) Count(key string) int {
	return c[key]
}

func (c Counter) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		if c[k] != 0 {
			keys = append(keys, k)
		}
	}
	return keys
}
