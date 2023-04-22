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
}

func TestStringObjectEvalutation(t *testing.T) {
	s := NewString("foobar")

	tests := []testObjectEvaluationCase{
		evalTest("!STRING(foobar)").
			call(s.OnPrefix(token.Bang)).
			expect(NewBoolean(false), true),
		evalTest("-STRING(foobar)").
			call(s.OnPrefix(token.Minus)).
			expect(nil, false),
		evalTest("STRING(foobar) + STRING(foobar)").
			call(s.OnInfix(token.Plus, s)).
			expect(NewString("foobarfoobar"), true),
		evalTest("STRING(foobar) == STRING(foobar)").
			call(s.OnInfix(token.EQ, s)).
			expect(NewBoolean(true), true),
		evalTest("STRING(foobar) != STRING(foobar)").
			call(s.OnInfix(token.NE, s)).
			expect(NewBoolean(false), true),
		evalTest("STRING(foobar) == STRING(foobar)").
			call(s.OnInfix(token.EQ, NewString("foobar"))).
			expect(NewBoolean(true), true),
		evalTest("STRING(foobar) == STRING(lorem ipsum)").
			call(s.OnInfix(token.EQ, NewString("lorem ipsum"))).
			expect(NewBoolean(false), true),
		evalTest("STRING(foobar) == INTEGER(42)").
			call(s.OnInfix(token.EQ, NewInteger(42))).
			expect(NewBoolean(false), true),
	}

	testObjectEvaluation(t, tests)
}

func TestStringObjectIndexEvaluation(t *testing.T) {
	s := NewString("foobar")

	tests := []testObjectEvaluationCase{
		evalTest("STRING(foobar)[INTEGER(0)]").
			call(s.OnIndex(NewInteger(0))).
			expect(NewString("f"), true),
		evalTest("STRING(foobar)[INTEGER(1)]").
			call(s.OnIndex(NewInteger(1))).
			expect(NewString("o"), true),
		evalTest("STRING(foobar)[INTEGER(2)]").
			call(s.OnIndex(NewInteger(2))).
			expect(NewString("o"), true),
		evalTest("STRING(foobar)[INTEGER(3)]").
			call(s.OnIndex(NewInteger(3))).
			expect(NewString("b"), true),
		evalTest("STRING(foobar)[INTEGER(10)]").
			call(s.OnIndex(NewInteger(10))).
			expect(NewNull(), true),
	}

	testObjectEvaluation(t, tests)
}
