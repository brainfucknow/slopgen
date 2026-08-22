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

func TestDeriveDistinctAcrossOverlappingRanges(t *testing.T) {
	seen := make(map[uint64][2]uint64)
	for root := uint64(0); root < 64; root++ {
		for label := uint64(0); label < 64; label++ {
			d := Derive(root, label)
			if prev, ok := seen[d]; ok {
				t.Fatalf("Derive(%d, %d) == Derive(%d, %d)", root, label, prev[0], prev[1])
			}
			seen[d] = [2]uint64{root, label}
		}
	}
}
