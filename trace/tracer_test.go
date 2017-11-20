package trace

import (
	"testing"
	"bytes"
)

func TestNew(t *testing.T)  {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Return from New should not be nil")
	} else {
		traceStr := "Hello trace package."
		tracer.Trace(traceStr)
		if buf.String() != traceStr + "\n" {
			t.Errorf("Trace should not write '%s'.", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("something")
}
