package main

import (
	"testing"
)

func TestInsert(t *testing.T) {
	m := make(hashMap[string], 4096)
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
