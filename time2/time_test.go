package time2

import "testing"

func TestFormatLayout(t *testing.T) {
	if FormatLayout("yyyymmdd-HHMMSS") != "20060102-150405" {
		t.Fail()
	}
}
