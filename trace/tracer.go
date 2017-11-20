package trace

import (
	"io"
	"fmt"
)

// Tracer is the interface that describes an object capable of
// tracing events throughout code.
type Tracer struct {
	out io.Writer
}
func (t *Tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

func New(w io.Writer) *Tracer {
	return &Tracer{out: w}
}
