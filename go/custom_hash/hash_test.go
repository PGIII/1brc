package main

import (
	"bufio"
	"os"
	"testing"
)

func TestInsert(t *testing.T) {
	m := make(hashMap[*string], 4096)
	cow, beef := []byte("cow"), "beef"
	pig, pork := []byte("pig"), "pork"
	chicken, poultry := []byte("chicken"), "poultry"
	m.insert(cow, &beef)
	m.insert(pig, &pork)
	m.insert(chicken, &poultry)

	found, _ := m.find(cow)
	if *found != beef {
		t.Fatalf("TestInsert: Got %s, Wanted %s", *found, beef)

	}

	found, _ = m.find(pig)
	if *found != pork {
		t.Fatalf("TestInsert: Got %s, Wanted %s", *found, pork)

	}

	found, _ = m.find(chicken)
	if *found != poultry {
		t.Fatalf("TestInsert: Got %s, Wanted %s", *found, poultry)

	}
}

func Test10kUniqueInsert(t *testing.T) {
	m := make(hashMap[FixedPoint], 1<<17)
	names := make([][]byte, 10000)
	file, err := os.Open("../../samples/measurements-10000-unique-keys.txt")
	//file, err := os.Open("../../samples/measurements-10.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	collisions := 0
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		name, num := parseLine(scanner.Bytes())

		if !m.contains(name) {
			names = append(names, name)
			collisions += m.insertChecked(name, num)
		}
	}

	t.Logf("Had: %d collisions\n", collisions)
	for _, n := range names {
		if !m.contains(n) {
			t.Fatalf("Test10kUniqueInsert: Missing %s in map", n)
		}
	}
}
