package server

import (
	. "github.com/cosiner/gomodule/cache"
)

type (
	// Values is stored type of attribute container
	Values map[string]interface{}
	// AttrContainer peformed as an ordinary cache
	AttrContainer OrdinaryCache
)

// NewAttrContainer return an new AttrContainer
func NewAttrContainer() *AttrContainer {
	return (*AttrContainer)(NewOrdinaryCache())
}

// newAttrContainerVals return an new AttrContainer initial with given values
func NewAttrContainerVals(values Values) *AttrContainer {
	return (*AttrContainer)(NewOrdinaryCacheVals(values))
}

// Attr return exist attribute value by name
func (ac *AttrContainer) Attr(name string) interface{} {
	return (*OrdinaryCache)(ac).Get(name)
}

// SetAttr store name-value pair to container
func (ac *AttrContainer) SetAttr(name string, value interface{}) {
	(*OrdinaryCache)(ac).Set(name, value)
}

// RemoveAttr remove an attribute by name
func (ac *AttrContainer) RemoveAttr(name string) {
	(*OrdinaryCache)(ac).Remove(name)
}

// IsAttrExist check whether given attribute is exist
func (ac *AttrContainer) IsAttrExist(name string) bool {
	return (*OrdinaryCache)(ac).IsExist(name)
}

// AccessAllAttrs access all attributes exist in container
func (ac *AttrContainer) AccessAllAttrs(fn func(Values)) {
	(*OrdinaryCache)(ac).AccessAllValues(func(values map[string]interface{}) { fn(values) })
}
