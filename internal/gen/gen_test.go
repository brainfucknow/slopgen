package gen

import (
	"bytes"
	"go/format"
	"testing"

	"github.com/openai/slopgen/internal/check"
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

func BenchmarkGenerate(b *testing.B) {
	cfg := Config{Functions: 20, Statements: 100}
	b.ReportAllocs()
	b.SetBytes(int64(len(Generate(0, cfg))))
	for i := 0; i < b.N; i++ {
		Generate(uint64(i), cfg)
	}
}
