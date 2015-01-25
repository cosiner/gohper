package server

import (
	"github.com/cosiner/gomodule/cache"
)

func _cache(ac *AttrContainer) *cache.OrdinaryCache {
	return (*cache.OrdinaryCache)(ac)
}

type AttrContainer cache.OrdinaryCache

func NewAttrContainer() *AttrContainer {
	return (*AttrContainer)(cache.NewOrdinaryCache())
}

func NewAttrContainerVals(values map[string]interface{}) *AttrContainer {
	return (*AttrContainer)(cache.NewOrdinaryCacheVals(values))
}

func (ac *AttrContainer) Attr(name string) interface{} {
	return _cache(ac).Get(name)
}

func (ac *AttrContainer) SetAttr(name string, value interface{}) {
	_cache(ac).Set(name, value)
}

func (ac *AttrContainer) RemoveAttr(name string) {
	_cache(ac).Remove(name)
}

func (ac *AttrContainer) IsAttrExist(name string) bool {
	return _cache(ac).IsExist(name)
}

func (ac *AttrContainer) AccessAllAttrs(fn func(map[string]interface{})) {
	_cache(ac).AccessAllValues(fn)
}
