package sortedmap

type Map struct {
	Indexes  map[string]int
	Elements []interface{}
}

func (m *Map) Set(key string, value interface{}) {
	if m.Indexes == nil {
		m.Indexes = make(map[string]int)
	}

	index, has := m.Indexes[key]
	if has {
		m.Elements[index] = value
	} else {
		m.Indexes[key] = len(m.Elements)
		m.Elements = append(m.Elements, value)
	}
}

// func (m *Map) Delete(key string) {
// 	delete(m.Indexes, key)
// }

func (m *Map) Has(key string) bool {
	if m.Indexes == nil {
		return false
	}

	_, has := m.Indexes[key]
	return has
}

func (m *Map) Get(key string) interface{} {
	if m.Indexes == nil {
		return nil
	}

	index, has := m.Indexes[key]
	if !has {
		return nil
	}

	return m.Elements[index]
}

func (m *Map) DefGet(key string, def interface{}) interface{} {
	if m.Indexes == nil {
		return def
	}

	index, has := m.Indexes[key]
	if !has {
		return def
	}

	return m.Elements[index]
}
