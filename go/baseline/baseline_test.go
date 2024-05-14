package main

import (
	"os"
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

func BenchmarkParse(b *testing.B) {
	nullFs, error := os.Create("/dev/null")
	if error != nil {
		panic(error)
	}

	defer nullFs.Close()
	for i := 0; i < b.N; i++ {
		parse(nullFs)
	}
}
