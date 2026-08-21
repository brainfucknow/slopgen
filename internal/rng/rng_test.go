package rng

import "testing"

func TestReproducibleAndSplit(t *testing.T) {
	a, b := New(42), New(42)
	for range 100 {
		if a.Uint64() != b.Uint64() {
			t.Fatal("same seed diverged")
		}
	}
	parent := New(9)
	x := parent.Split(3).Uint64()
	parent.Uint64()
	if x != parent.Split(3).Uint64() {
		t.Fatal("split depends on parent counter")
	}
}
