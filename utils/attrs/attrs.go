// Package attrs provide two attributes containers, designed for embed to other types
package attrs

import "sync"

type (
	// Attrs is a common container store attribute
	Attrs interface {
		Attr(name string) interface{}
		AttrDef(name string, defVal interface{}) interface{}
		// if value is nil, remove it
		SetAttr(name string, value interface{})
		GetSetAttr(name string, value interface{}) interface{}
		IsAttrExist(name string) bool

		AllAttrs() Values
		Clear()
	}

	Values map[string]interface{}

	LockedValues struct {
		Values
		sync.RWMutex
	}
)

func New() Attrs {
	return make(Values)
}

func NewLocked() Attrs {
	return &LockedValues{
		Values: make(Values),
	}
}

func (v Values) Attr(key string) interface{} {
	return v[key]
}

func (v Values) AttrDef(key string, def interface{}) interface{} {
	val := v[key]
	if val == nil {
		return def
	}

	return val
}

func (v Values) SetAttr(key string, val interface{}) {
	if val != nil {
		v[key] = val
	} else {
		delete(v, key)
	}
}

func (v Values) GetSetAttr(key string, val interface{}) interface{} {
	value := v.Attr(key)
	v.SetAttr(key, val)
	return value
}

func (v Values) IsAttrExist(key string) bool {
	_, has := v[key]
	return has
}

func (v Values) AllAttrs() Values {
	return v
}

func (v Values) Clear() {
	for k := range v {
		delete(v, k)
	}
}

func (v *LockedValues) Attr(key string) interface{} {
	v.RLock()
	val := v.Values.Attr(key)
	v.RUnlock()
	return val
}

func (v *LockedValues) AttrDef(key string, def interface{}) interface{} {
	v.RLock()
	val := v.Values.AttrDef(key, def)
	v.RUnlock()
	return val
}

func (v *LockedValues) SetAttr(key string, val interface{}) {
	v.Lock()
	v.Values.SetAttr(key, val)
	v.Unlock()
}

func (v *LockedValues) GetSetAttr(key string, val interface{}) interface{} {
	v.Lock()
	val = v.Values.GetSetAttr(key, val)
	v.Unlock()
	return val
}

func (v *LockedValues) IsAttrExist(key string) bool {
	v.RLock()
	has := v.Values.IsAttrExist(key)
	v.RUnlock()
	return has
}
func (v *LockedValues) AllAttrs() Values {
	return v.Values
}
func (v *LockedValues) Clear() {
	v.Lock()
	v.Values.Clear()
	v.Unlock()
}
