package runtime2

import (
	"testing"

	"github.com/cosiner/gohper/strings2"
)

func TestCaller(t *testing.T) {
	if p := strings2.RemoveSpace(Caller(0)); p != "runtime_test.go:runtime2.TestCaller:10" {
		t.Fatalf("Error: expect runtime_test.go: runtime2.TestCaller: 10, but get %s", p)
	}
}
