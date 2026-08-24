// Package env tracks declared bindings indexed by type.
package env

import "github.com/brainfucknow/slopgen/internal/scalar"

// Table records the names declared for each type. M0 declares every binding
// at function scope, so a flat table is sufficient; a scope stack returns
// with block-scoped declarations in M1.
type Table struct {
	byType map[scalar.ID][]string
}

func New() *Table { return &Table{byType: make(map[scalar.ID][]string)} }

func (e *Table) Declare(name string, typ scalar.ID) {
	e.byType[typ] = append(e.byType[typ], name)
}

// Lookup returns the names bound with type typ. The returned slice aliases
// internal state and must not be mutated.
func (e *Table) Lookup(typ scalar.ID) []string { return e.byType[typ] }
