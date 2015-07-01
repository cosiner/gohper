package sortedmap

type Element struct {
	Key   string
	Value interface{}
}

type Map struct {
	Indexes map[string]int
	Values  []Element
}

func New() Map {
	return Map{
		Indexes: make(map[string]int),
	}
}

func (m *Map) Set(key string, value interface{}) {
	index, has := m.Indexes[key]
	if has {
		m.Values[index] = Element{
			Key:   key,
			Value: value,
		}
	} else {
		m.Indexes[key] = len(m.Values)
		m.Values = append(m.Values, Element{
			Key:   key,
			Value: value,
		})
	}
}

func (m *Map) Delete(key string) {
	index, has := m.Indexes[key]
	if !has {
		return
	}
	delete(m.Indexes, key)

	for l := len(m.Values); index < l-1; index++ {
		v1 := &m.Values[index]
		v2 := &m.Values[index+1]
		v1.Key = v2.Key
		v1.Value = v2.Value
		m.Indexes[v2.Key] = index
	}

	m.Values = m.Values[:index]
}

func (m *Map) HasKey(key string) bool {
	_, has := m.Indexes[key]
	return has
}

func (m *Map) Get(key string) interface{} {
	index, has := m.Indexes[key]
	if !has {
		return nil
	}

	return m.Values[index].Value
}

func (m *Map) DefGet(key string, def interface{}) interface{} {
	index, has := m.Indexes[key]
	if !has {
		return def
	}

	return m.Values[index].Value
}

func (m *Map) Clear() {
	for k := range m.Indexes {
		delete(m.Indexes, k)
	}

	m.Values = m.Values[:0]
}
