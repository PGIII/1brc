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

func TestParseLine(t *testing.T) {
	in := []byte("Hamburg;12.0")
	wantName := "Hamburg"
	wantNum := FixedPoint(1200)
	outName, outNum := parseLine(in)
	if string(outName) != wantName || outNum != wantNum {
		t.Fatalf(`parseLine(%s) = %s, %d, want %s, %d`, in, string(outName), outNum, string(wantName), wantNum)
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

func BenchmarkOriginalParseLine(b *testing.B) {
	str := []byte("Hamburg;12.0")
	for i := 0; i < b.N; i++ {
		originalParseLine(str)
	}
}

func BenchmarkLineParse(b *testing.B) {
	str := []byte("Hamburg;12.0")
	for i := 0; i < b.N; i++ {
		parseLine(str)
	}
}
