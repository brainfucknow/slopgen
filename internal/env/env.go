// Package env tracks lexical bindings using type-indexed scope tables.
package env

import "github.com/openai/slopgen/internal/scalar"

type Binding struct {
	Name string
	Type scalar.ID
	Used bool
}
type Scope struct {
	byType  map[scalar.ID][]*Binding
	names   map[string]bool
	ordered []*Binding
}
type Stack struct{ scopes []*Scope }

func New() *Stack { e := &Stack{}; e.Push(); return e }
func (e *Stack) Push() {
	e.scopes = append(e.scopes, &Scope{byType: make(map[scalar.ID][]*Binding), names: make(map[string]bool)})
}
func (e *Stack) Pop() []*Binding {
	s := e.scopes[len(e.scopes)-1]
	e.scopes = e.scopes[:len(e.scopes)-1]
	return s.ordered
}
func (e *Stack) Declare(name string, typ scalar.ID) *Binding {
	s := e.scopes[len(e.scopes)-1]
	b := &Binding{Name: name, Type: typ}
	s.names[name] = true
	s.byType[typ] = append(s.byType[typ], b)
	s.ordered = append(s.ordered, b)
	return b
}
func (e *Stack) ContainsCurrent(name string) bool { return e.scopes[len(e.scopes)-1].names[name] }
func (e *Stack) Lookup(typ scalar.ID) []*Binding {
	var out []*Binding
	for i := len(e.scopes) - 1; i >= 0; i-- {
		out = append(out, e.scopes[i].byType[typ]...)
	}
	return out
}
