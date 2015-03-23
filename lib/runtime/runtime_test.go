package runtime

import "testing"

func TestCallerPosition(t *testing.T) {
	if p := CallerPosition(0); p != "runtime_test.go: runtime.TestCallerPosition: 6" {
		t.Fatalf("Error: expect runtime_test.go: runtime.TestCallerPosition: 6, but get %s", p)
	}
}
