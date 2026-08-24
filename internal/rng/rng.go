// Package rng provides deterministic, splittable random streams.
package rng

// Stream is a counter-based SplitMix64 stream. A stream's output depends only
// on its key and counter, so independently derived streams remain reproducible
// regardless of scheduling.
type Stream struct {
	key     uint64
	counter uint64
}

func New(seed uint64) *Stream { return &Stream{key: mix(seed)} }

// Split derives an independent stream without consuming the parent.
func (s *Stream) Split(label uint64) *Stream { return &Stream{key: mix(s.key ^ mix(label))} }

// Derive maps a root seed and a label (such as a file index) to an
// independent seed. Unlike root+label arithmetic, overlapping label ranges
// under different roots do not produce the same derived seed, so output from
// separate runs can safely share a package. Only the seed is pre-mixed,
// keeping the combination asymmetric: mixing both operands would make
// Derive(a, b) equal Derive(b, a).
func Derive(seed, label uint64) uint64 { return mix(mix(seed) ^ label) }

func (s *Stream) Uint64() uint64 {
	v := mix(s.key + s.counter*0x9e3779b97f4a7c15)
	s.counter++
	return v
}

func (s *Stream) Intn(n int) int {
	if n <= 0 {
		panic("rng: non-positive bound")
	}
	return int(s.Uint64() % uint64(n))
}

func (s *Stream) Chance(numerator, denominator uint64) bool {
	if denominator == 0 || numerator > denominator {
		panic("rng: invalid probability")
	}
	return s.Uint64()%denominator < numerator
}

func mix(x uint64) uint64 {
	x += 0x9e3779b97f4a7c15
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	return x ^ (x >> 31)
}
