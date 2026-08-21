// Package scalar defines the M0 type universe.
package scalar

import "strconv"

type ID uint8

const (
	Int ID = iota
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float64
	String
	Bool
)

var all = []ID{Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Float64, String, Bool}

func All() []ID { return all }

func (t ID) String() string {
	return [...]string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float64", "string", "bool"}[t]
}
func (t ID) Numeric() bool { return t <= Float64 }
func (t ID) Integer() bool { return t <= Uint64 }

// Literal returns a representable literal for t from arbitrary bits.
func (t ID) Literal(bits uint64, nonzero bool) string {
	switch t {
	case Bool:
		if bits&1 == 0 {
			return "false"
		}
		return "true"
	case String:
		return strconv.Quote([]string{"alpha", "value", "markov", "source", "node"}[bits%5])
	case Float64:
		v := bits % 20001
		if nonzero && v == 10000 {
			v++
		}
		return strconv.FormatFloat(float64(int64(v)-10000)/10, 'f', 1, 64)
	}
	width, signed := 64, t <= Int64
	switch t {
	case Int8, Uint8:
		width = 8
	case Int16, Uint16:
		width = 16
	case Int32, Uint32:
		width = 32
	}
	if signed {
		// Keep literals deliberately modest for architecture-independent int.
		limit := uint64(1) << min(width-1, 30)
		v := int64(bits%(2*limit)) - int64(limit)
		if nonzero && v == 0 {
			v = 1
		}
		return strconv.FormatInt(v, 10)
	}
	limit := uint64(1) << min(width, 31)
	v := bits % limit
	if nonzero && v == 0 {
		v = 1
	}
	return strconv.FormatUint(v, 10)
}
