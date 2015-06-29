package utils

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestPageStart(t *testing.T) {
	testing2.
		Expect(0).Arg("abcde", 10).
		Expect(0).Arg("-1", 10).
		Expect(0).Arg("0", 10).
		Expect(0).Arg("1", 10).
		Expect(10).Arg("2", 10).
		Expect(20).Arg("3", 10).
		Run(t, PageStart)
}
