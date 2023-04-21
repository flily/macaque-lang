package object

import (
	"testing"
)

func TestFunctionObject(t *testing.T) {
	f := NewFunction(0, 42, nil)

	if f.Type() != ObjectTypeFunction {
		t.Errorf("f.Type() is not ObjectTypeFunction. got=%T (%+v)", f.Type(), f.Type())
	}

	if f.Inspect() != "function" {
		t.Errorf("f.Inspect() wrong. got=%q", f.Inspect())
	}

	if f.Hashable() {
		t.Errorf("f.Hashable() is true")
	}

	if f.HashKey() != nil {
		t.Errorf("f.HashKey() is not nil")
	}
}
