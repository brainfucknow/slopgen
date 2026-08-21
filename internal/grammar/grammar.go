// Package grammar selects legal productions with depth-damped weights.
package grammar

import "github.com/openai/slopgen/internal/rng"

type Production uint8

const (
	Literal Production = iota
	Binding
	Binary
	Convert
	Call
)

// ChooseExpr chooses an expression production. Recursive choices decay with
// depth, and only terminating productions remain at the configured cap.
func ChooseExpr(r *rng.Stream, depth, cap int, bindings, calls bool) Production {
	type candidate struct {
		p         Production
		weight    int
		recursive bool
	}
	cs := []candidate{{Literal, 5, false}}
	if bindings {
		cs = append(cs, candidate{Binding, 8, false})
	}
	if calls {
		cs = append(cs, candidate{Call, 3, false})
	}
	if depth < cap && bindings {
		damp := cap - depth
		cs = append(cs, candidate{Binary, 2 * damp, true}, candidate{Convert, damp, true})
	}
	total := 0
	for _, c := range cs {
		total += c.weight
	}
	pick := r.Intn(total)
	for _, c := range cs {
		if pick < c.weight {
			return c.p
		}
		pick -= c.weight
	}
	return Literal
}
