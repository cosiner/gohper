package server

import (
	. "github.com/cosiner/gomodule/cache"
)

type (
	// AttrContainer is a common container store attribute
	AttrContainer interface {
		Attr(name string) interface{}
		SetAttr(name string, value interface{})
		IsAttrExist(name string) bool
		AccessAllAttrs(fn func(Values))
	}

	// Values is stored type of attribute container
	Values map[string]interface{}
	// lockedValues peformed as an ordinary cache with lock
	lockedValues OrdinaryCache
)

// NewAttrContainer return an new AttrContainer with lock
func NewAttrContainer() AttrContainer {
	return (*lockedValues)(NewOrdinaryCache())
}

// newAttrContainerVals return an new AttrContainer initial with given values
// and lock
func NewAttrContainerVals(values Values) AttrContainer {
	return (*lockedValues)(NewOrdinaryCacheVals(values))
}

// Attr return exist attribute value by name
func (v Values) Attr(name string) interface{} {
	return v[name]
}

// SetAttr store name-value pair to container
func (v Values) SetAttr(name string, value interface{}) {
	v[name] = value
}

// RemoveAttr remove an attribute by name
func (v Values) RemoveAttr(name string) {
	delete(v, name)
}

// IsAttrExist check whether given attribute is exist
func (v Values) IsAttrExist(name string) bool {
	_, has := v[name]
	return has
}

// AccessAllAttrs access all attributes exist in container
func (v Values) AccessAllAttrs(fn func(Values)) {
	fn(v)
}

// Attr return exist attribute value by name
func (lc *lockedValues) Attr(name string) interface{} {
	return (*OrdinaryCache)(lc).Get(name)
}

// SetAttr store name-value pair to container
func (lc *lockedValues) SetAttr(name string, value interface{}) {
	(*OrdinaryCache)(lc).Set(name, value)
}

// RemoveAttr remove an attribute by name
func (lc *lockedValues) RemoveAttr(name string) {
	(*OrdinaryCache)(lc).Remove(name)
}

// IsAttrExist check whether given attribute is exist
func (lc *lockedValues) IsAttrExist(name string) bool {
	return (*OrdinaryCache)(lc).IsExist(name)
}

// AccessAllAttrs access all attributes exist in container
func (lc *lockedValues) AccessAllAttrs(fn func(Values)) {
	(*OrdinaryCache)(lc).AccessAllValues(func(values map[string]interface{}) { fn(values) })
}
