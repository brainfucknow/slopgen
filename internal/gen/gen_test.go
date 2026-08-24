package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/brainfucknow/slopgen/internal/check"
	"github.com/brainfucknow/slopgen/internal/rng"
)

func TestGeneratedFilesTypeCheckAndAreDeterministic(t *testing.T) {
	cfg := Config{Functions: 20, Statements: 30, MaxDepth: 4}
	for seed := uint64(0); seed < 100; seed++ {
		a, b := Generate(seed, cfg), Generate(seed, cfg)
		if !bytes.Equal(a, b) {
			t.Fatalf("seed %d is not deterministic", seed)
		}
		if err := check.Source("generated.go", a); err != nil {
			t.Fatalf("seed %d: %v\n%s", seed, err, a)
		}
		formatted, err := format.Source(a)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(a, formatted) {
			t.Fatalf("seed %d is not gofmt canonical", seed)
		}
	}
}

// Two runs with overlapping per-file index ranges model the CLI's seed
// derivation: with the old seed+index scheme these runs shared derived seeds
// and redeclared identical function names.
func TestGeneratedFilesTypeCheckAsOnePackage(t *testing.T) {
	cfg := Config{Functions: 20, Statements: 5, MaxDepth: 4}
	fset := token.NewFileSet()
	var files []*ast.File
	for _, root := range []uint64{42, 92} {
		for i := uint64(0); i < 60; i++ {
			name := fmt.Sprintf("run%d_generated_%06d.go", root, i)
			file, err := parser.ParseFile(fset, name, Generate(rng.Derive(root, i), cfg), 0)
			if err != nil {
				t.Fatalf("parse %s: %v", name, err)
			}
			files = append(files, file)
		}
	}
	if _, err := new(types.Config).Check("generated", fset, files, nil); err != nil {
		t.Fatalf("type check generated package: %v", err)
	}
}

func BenchmarkGenerate(b *testing.B) {
	cfg := Config{Functions: 20, Statements: 100}
	b.ReportAllocs()
	b.SetBytes(int64(len(Generate(0, cfg))))
	for i := 0; i < b.N; i++ {
		Generate(uint64(i), cfg)
	}
}
