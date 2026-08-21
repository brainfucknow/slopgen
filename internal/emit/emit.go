// Package emit is a minimal gofmt-compatible source writer.
package emit

import "bytes"

type Writer struct {
	b      bytes.Buffer
	indent int
}

func (w *Writer) Line(s string) {
	for range w.indent {
		w.b.WriteByte('\t')
	}
	w.b.WriteString(s)
	w.b.WriteByte('\n')
}
func (w *Writer) Blank()         { w.b.WriteByte('\n') }
func (w *Writer) Open(s string)  { w.Line(s + " {"); w.indent++ }
func (w *Writer) Close()         { w.indent--; w.Line("}") }
func (w *Writer) Bytes() []byte  { return w.b.Bytes() }
func (w *Writer) String() string { return w.b.String() }
