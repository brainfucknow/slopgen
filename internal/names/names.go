// Package names generates identifier-like names from a tiny character Markov model.
package names

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/brainfucknow/slopgen/internal/rng"
)

var corpus = []string{"value", "index", "count", "result", "current", "source", "target", "node", "item", "total", "limit", "state", "buffer", "offset"}

// Generator holds first-order transitions learned from the built-in corpus.
type Generator struct {
	next   map[byte][]byte
	starts []byte
}

func New() *Generator {
	g := &Generator{next: make(map[byte][]byte)}
	for _, word := range corpus {
		g.starts = append(g.starts, word[0])
		prev := byte('^')
		for i := 0; i < len(word); i++ {
			g.next[prev] = append(g.next[prev], word[i])
			prev = word[i]
		}
		g.next[prev] = append(g.next[prev], 0)
	}
	return g
}

func (g *Generator) Name(r *rng.Stream, exported bool, suffix int) string {
	var b strings.Builder
	c := g.starts[r.Intn(len(g.starts))]
	for b.Len() < 10 && c != 0 {
		b.WriteByte(c)
		choices := g.next[c]
		c = choices[r.Intn(len(choices))]
	}
	name := b.String()
	if len(name) < 2 {
		name = "value"
	}
	if exported {
		rs := []rune(name)
		rs[0] = unicode.ToUpper(rs[0])
		name = string(rs)
	}
	if suffix >= 0 {
		name += strconv.Itoa(suffix)
	}
	return name
}
