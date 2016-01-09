package runtime2

import (
	"testing"

	"github.com/cosiner/gohper/strings2"
)

func TestCaller(t *testing.T) {
	exp := "runtime2/runtime_test.go:11"
	if p := strings2.RemoveSpace(Caller(0)); p != exp {
		t.Fatalf("Error: expect %s, but get %s", exp, p)
	}
}
