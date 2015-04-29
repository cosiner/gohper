package runtime

import (
	"testing"

    "github.com/cosiner/gohper/lib/types"
)

func TestCallerPosition(t *testing.T) {
	if p := types.RemoveSpace(CallerPosition(0)); p != "runtime_test.go:runtime.TestCallerPosition:10" {
		t.Fatalf("Error: expect runtime_test.go: runtime.TestCallerPosition: 10, but get %s", p)
	}
}
