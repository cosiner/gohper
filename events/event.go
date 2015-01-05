// Package events implements a events register
package events

// On call fn when event is true
func On(event bool, fn func()) {
	if event {
		fn()
	}
}

type Event uint

type Eventer struct {
	eventMapping map[string]func(event Event) error
}
