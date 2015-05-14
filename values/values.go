package values

import "sync"

type (
	Values map[string]interface{}

	LockedValues struct {
		Values
		*sync.RWMutex
	}
)

func NewLocked() *LockedValues {
	return &LockedValues{
		Values:  New(),
		RWMutex: new(sync.RWMutex),
	}
}

func NewLockedWith(v map[string]interface{}) *LockedValues {
	return &LockedValues{
		Values:  v,
		RWMutex: new(sync.RWMutex),
	}
}

func New() Values {
	return make(Values)
}

func NewWith(v map[string]interface{}) Values {
	return Values(v)
}

func (v Values) IsExist(key string) bool {
	_, has := v[key]
	return has
}

func (v Values) Size() int {
	return len(v)
}

func (v Values) Get(key string) interface{} {
	return v[key]
}

func (v Values) Remove(key string) {
	delete(v, key)
}

func (v Values) RandomRemove() {
	for k := range v {
		v.Remove(k)
		break
	}
}

func (v Values) Set(key string, val interface{}) {
	v[key] = val
}

func (v Values) Update(key string, val interface{}) (s bool) {
	if s = v.IsExist(key); s {
		v[key] = val
	}

	return
}

func (v Values) AccessAll(fn func(Values)) {
	fn(v)
}

func (lc *LockedValues) Size() int {
	lc.RLock()
	size := lc.Values.Size()
	lc.RUnlock()

	return size
}

func (lc *LockedValues) Get(key string) (val interface{}) {
	lc.RLock()
	val = lc.Values.Get(key)
	lc.RUnlock()

	return
}

func (lc *LockedValues) IsExist(key string) bool {
	lc.RLock()
	has := lc.Values.IsExist(key)
	lc.RUnlock()

	return has
}

func (lc *LockedValues) Remove(key string) {
	lc.Lock()
	lc.Values.Remove(key)
	lc.Unlock()
}

func (lc *LockedValues) RandomRemove() {
	lc.Lock()
	lc.Values.RandomRemove()
	lc.Unlock()
}

func (lc *LockedValues) Set(key string, val interface{}) {
	lc.Lock()
	lc.Values.Set(key, val)
	lc.Unlock()
}

func (lc *LockedValues) Update(key string, val interface{}) (s bool) {
	lc.Lock()
	s = lc.Values.Update(key, val)
	lc.Unlock()

	return
}

func (lc *LockedValues) AccessAll(fn func(Values)) {
	lc.RLock()
	lc.Values.AccessAll(fn)
	lc.RUnlock()
}
