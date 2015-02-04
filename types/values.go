package types

import "sync"

type (
	// AttrContainer is a common container store attribute
	AttrContainer interface {
		Attr(name string) interface{}
		SetAttr(name string, value interface{})
		RemoveAttr(name string)
		IsAttrExist(name string) bool
		AccessAllAttrs(fn func(Values))
	}

	Values map[string]interface{}

	LockedValues struct {
		Values
		*sync.RWMutex
	}
)

func NewLockedValues() *LockedValues {
	return &LockedValues{
		Values:  NewValues(),
		RWMutex: new(sync.RWMutex),
	}
}

func NewLockedValuesWith(v map[string]interface{}) *LockedValues {
	return &LockedValues{
		Values:  v,
		RWMutex: new(sync.RWMutex),
	}
}

func NewValues() Values {
	return make(Values)
}

func NewValuesWith(v map[string]interface{}) Values {
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

func (v Values) IsAttrExist(key string) bool {
	return v.IsExist(key)
}

func (v Values) Attr(key string) interface{} {
	return v.Get(key)
}

func (v Values) RemoveAttr(key string) {
	v.Remove(key)
}

func (v Values) SetAttr(key string, val interface{}) {
	v.Set(key, val)
}

func (v Values) UpdateAttr(key string, val interface{}) bool {
	return v.Update(key, val)
}

func (v Values) AccessAllAttrs(fn func(Values)) {
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

func (lc *LockedValues) IsAttrExist(key string) bool {
	return lc.IsExist(key)
}

func (lc *LockedValues) Attr(key string) interface{} {
	return lc.Get(key)
}

func (lc *LockedValues) RemoveAttr(key string) {
	lc.Remove(key)
}

func (lc *LockedValues) SetAttr(key string, val interface{}) {
	lc.Set(key, val)
}

func (lc *LockedValues) UpdateAttr(key string, val interface{}) bool {
	return lc.Update(key, val)
}

func (lc *LockedValues) AccessAllAttrs(fn func(Values)) {
	lc.AccessAll(fn)
}
