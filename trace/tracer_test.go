package trace

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Return from New should not be a nil")
	} else {
		msg := "Hello trace package"
		tracer.Trace(msg)
		if buf.String() != fmt.Sprint(msg, "\n") {
			t.Errorf("Trace should not write '%s'.", buf.String())
		}
	}
}
