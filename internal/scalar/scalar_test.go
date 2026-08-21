package scalar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"
)

func TestLiteralsRepresentable(t *testing.T) {
	for _, typ := range All() {
		for _, bits := range []uint64{0, 1, ^uint64(0), 1 << 63} {
			lit := typ.Literal(bits, false)
			if typ == String || typ == Bool || typ == Float64 {
				continue
			}
			src := "package p\nvar _ " + typ.String() + " = " + lit
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "literal.go", src, 0)
			if err == nil {
				_, err = new(types.Config).Check("p", fset, []*ast.File{f}, nil)
			}
			if err != nil {
				t.Fatalf("%s is not a %s: %v", lit, typ, err)
			}
		}
	}
}

func TestNonzero(t *testing.T) {
	for _, typ := range All() {
		if typ.Numeric() && typ.Literal(0, true) == "0" {
			t.Fatalf("zero %s", typ)
		}
	}
}
