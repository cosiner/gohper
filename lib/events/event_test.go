package events

import (
	"fmt"
	"testing"
)

func TestEvent(t *testing.T) {
	// tt := test.WrapTest(t)
	On(fmt.Println).Do(func() {
		fmt.Println()
	})
}
