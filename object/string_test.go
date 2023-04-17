package object

import (
	"testing"

	"github.com/flily/macaque-lang/token"
)

func TestStringObject(t *testing.T) {
	s := NewString("foobar")

	if s.Type() != ObjectTypeString {
		t.Errorf("string.Type() is not ObjectTypeString")
	}

	if s.Inspect() != "foobar" {
		t.Errorf("string.Inspect() wrong, expected %q, got %q",
			"foobar", s.Inspect())
	}

	if !s.Hashable() {
		t.Errorf("string.Hashable() shoud be true")
	}

	if s.HashKey() != "foobar" {
		t.Errorf("string.HashKey() wrong, expected %q, got %q",
			"foobar", s.HashKey())
	}

	l1 := NewString("foobar")
	l2 := NewString("lorem ipsum")
	if !l1.EqualTo(l1) {
		t.Errorf("string(%s) should be equal to string(%s)",
			l1.Inspect(), l1.Inspect())
	}

	if l1.EqualTo(l2) {
		t.Errorf("string(%s) should not be equal to string(%s)",
			l1.Inspect(), l2.Inspect())
	}
}

func TestStringObjectEvalutation(t *testing.T) {
	s := NewString("foobar")

	tests := []testObjectEvaluationCase{
		evalTest("STRING + STRING").
			call(s.OnInfix(token.Plus, s)).
			expect(NewString("foobarfoobar"), true),
		evalTest("STRING == STRING").
			call(s.OnInfix(token.EQ, s)).
			expect(NewBoolean(true), true),
		evalTest("STRING != STRING").
			call(s.OnInfix(token.NE, s)).
			expect(NewBoolean(false), true),
		evalTest("STRING == INTEGER").
			call(s.OnInfix(token.EQ, NewInteger(42))).
			expect(NewBoolean(false), true),
		evalTest("!STRING").
			call(s.OnPrefix(token.Bang)).
			expect(NewBoolean(false), true),
	}

	testObjectEvaluation(t, tests)
}
