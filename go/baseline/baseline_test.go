package main

import (
	"testing"
)

func TestParseTemperature(t *testing.T) {
	in := []byte("-100.11")
	want := FixedPoint(-100110)
	out := parseBytesToFixedPoint(in)
	if out != want {
		t.Fatalf(`parseBytesToFixedPoint(%s) = %d, want %d`, in, out, want)
	}
}
