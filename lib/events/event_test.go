package events

import (
	"fmt"
	"testing"
)

func TestEvent(t *testing.T) {
	// tt := test.Wrap(t)
	On(fmt.Println).Do(func() {
		fmt.Println()
	})
}
