package net2

import "testing"

func TestAddr(t *testing.T) {
	t.Log(ReplaceHost(":1080", Localhost()))
}
