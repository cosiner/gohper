// Package events implements a events register
package events

type On struct {
	val interface{}
}

func When(o interface{}) On {
	return On{val: o}
}

func (o On) Do(fn func()) {
	if test(o.val) && fn != nil {
		fn()
	}
}

func (o On) WithDo(fn func(o interface{})) {
	if test(o) && fn != nil {
		fn(interface{}(o))
	}
}

func test(o interface{}) bool {
	var val bool
	switch o := o.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64, complex64, complex128:
		val = (o != 0)
	case string:
		val = (o != "")
	case bool:
		val = o
	default:
		val = (o != nil)
	}
	return val
}
