package object

import (
	"testing"

	"github.com/flily/macaque-lang/token"
)

func TestIntegerObject(t *testing.T) {
	i := NewInteger(42)

	if i.Type() != ObjectTypeInteger {
		t.Errorf("integer.Type() is not ObjectTypeInteger")
	}

	if i.Inspect() != "42" {
		t.Errorf("integer.Inspect() wrong, expected %q, got %q",
			"42", i.Inspect())
	}

	if !i.Hashable() {
		t.Errorf("integer.Hashable() is not true")
	}

	if i.HashKey() != int64(42) {
		t.Errorf("integer.HashKey() is not 42, got %v", i.HashKey())
	}
}

func TestIntegerObjectPrefixEvaluation(t *testing.T) {
	i := NewInteger(42)

	tests := []testObjectEvaluationCase{
		evalTest("-INTEGER").
			call(i.OnPrefix(token.Minus)).
			expect(NewInteger(-42), true),
		evalTest("!INTEGER").
			call(i.OnPrefix(token.Bang)).
			expect(NewBoolean(false), true),
		evalTest("~INTEGER").
			call(i.OnPrefix(token.BITNOT)).
			expect(NewInteger(^42), true),
	}

	testObjectEvaluation(t, tests)
}

func TestIntegerObjectInfixOnIntegerEvaluation(t *testing.T) {
	i := NewInteger(42)
	j := NewInteger(2)

	tests := []testObjectEvaluationCase{
		evalTest("INTEGER(42) == INTEGER(2)").
			call(i.OnInfix(token.EQ, j)).
			expect(NewBoolean(false), true),
		evalTest("INTEGER(42) == INTEGER(42)").
			call(i.OnInfix(token.EQ, NewInteger(42))).
			expect(NewBoolean(true), true),
		evalTest("INTEGER(42) != INTEGER(2)").
			call(i.OnInfix(token.NE, j)).
			expect(NewBoolean(true), true),
		evalTest("INTEGER(42) == STRING(42)").
			call(i.OnInfix(token.EQ, NewString("42"))).
			expect(NewBoolean(false), true),
		evalTest("INTEGER(42) != STRING(42)").
			call(i.OnInfix(token.NE, NewString("42"))).
			expect(NewBoolean(true), true),
		evalTest("INTEGER(42) + INTEGER(2)").
			call(i.OnInfix(token.Plus, j)).
			expect(NewInteger(44), true),
		evalTest("INTEGER(42) - INTEGER(2)").
			call(i.OnInfix(token.Minus, j)).
			expect(NewInteger(40), true),
		evalTest("INTEGER(42) * INTEGER(2)").
			call(i.OnInfix(token.Asterisk, j)).
			expect(NewInteger(84), true),
		evalTest("INTEGER(42) / INTEGER(2)").
			call(i.OnInfix(token.Slash, j)).
			expect(NewInteger(21), true),
		evalTest("INTEGER(42) % INTEGER(2)").
			call(i.OnInfix(token.Modulo, j)).
			expect(NewInteger(0), true),
		evalTest("INTEGER(42) < INTEGER(2)").
			call(i.OnInfix(token.LT, j)).
			expect(NewBoolean(false), true),
		evalTest("INTEGER(42) > INTEGER(2)").
			call(i.OnInfix(token.GT, j)).
			expect(NewBoolean(true), true),
		evalTest("INTEGER(42) <= INTEGER(2)").
			call(i.OnInfix(token.LE, j)).
			expect(NewBoolean(false), true),
		evalTest("INTEGER(42) >= INTEGER(2)").
			call(i.OnInfix(token.GE, j)).
			expect(NewBoolean(true), true),
		evalTest("INTEGER(42) && INTEGER(2)").
			call(i.OnInfix(token.AND, j)).
			expect(NewBoolean(true), true),
		evalTest("INTEGER(42) || INTEGER(2)").
			call(i.OnInfix(token.OR, j)).
			expect(NewBoolean(true), true),
		evalTest("INTEGER(42) & INTEGER(2)").
			call(i.OnInfix(token.BITAND, j)).
			expect(NewInteger(42&2), true),
		evalTest("INTEGER(42) | INTEGER(2)").
			call(i.OnInfix(token.BITOR, j)).
			expect(NewInteger(42|2), true),
		evalTest("INTEGER(42) ^ INTEGER(2)").
			call(i.OnInfix(token.BITXOR, j)).
			expect(NewInteger(42^2), true),
	}

	testObjectEvaluation(t, tests)
}

func TestIntegerObjectOnIndexEvaluation(t *testing.T) {
	i := NewInteger(42)

	tests := []testObjectEvaluationCase{
		evalTest("INTEGER(42)[0]").
			call(i.OnIndex(NewInteger(0))).
			expect(nil, false),
	}

	testObjectEvaluation(t, tests)
}
