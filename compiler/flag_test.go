package compiler

import (
	"testing"
)

func TestCompilerFlagString(t *testing.T) {
	f := NewFlag(FlagPackValue, FlagCleanStack)
	if f.String() != "CompilerFlag(0000-0000-0000-0101)" {
		t.Errorf("wrong string: %s", f)
	}
}

func TestCompiilerFlagSet(t *testing.T) {
	f1 := NewFlag(FlagPackValue)

	f2 := f1.With(FlagCleanStack)
	f3 := f1.Set(FlagWithReturn)

	if !f2.Has(FlagPackValue) || !f2.Has(FlagCleanStack) {
		t.Errorf("FlagPackValue and FlagCleanStack should be set: %s", f2)
	}

	if f1.Has(FlagCleanStack) {
		t.Errorf("FlagCleanStack on f1 should not be set: %s", f1)
	}

	if f3 != f1 {
		t.Errorf("f3 should be f1: %s", f3)
	}
}

func TestCompilerClear(t *testing.T) {
	f1 := NewFlag(FlagPackValue, FlagCleanStack, FlagWithReturn)
	f2 := f1.Without(FlagCleanStack)
	f3 := f1.Clear(FlagWithReturn)

	if f2.Has(FlagCleanStack) {
		t.Errorf("FlagCleanStack on f2 should not be set: %s", f2)
	}

	if !f1.Has(FlagCleanStack) {
		t.Errorf("FlagCleanStack on f1 should be set: %s", f1)
	}

	if f3.Has(FlagWithReturn) {
		t.Errorf("FlagWithReturn on f3 should not be set: %s", f3)
	}

	if f1.Has(FlagWithReturn) {
		t.Errorf("FlagWithReturn on f1 should not be set: %s", f1)
	}
}
