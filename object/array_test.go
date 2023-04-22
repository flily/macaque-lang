package object

import (
	"testing"

	"github.com/flily/macaque-lang/token"
)

func TestArrayObject(t *testing.T) {
	a := NewArray([]Object{
		NewInteger(1),
		NewInteger(2),
		NewInteger(3),
	})

	if a.Type() != ObjectTypeArray {
		t.Errorf("array.Type() is not ObjectTypeArray")
	}

	if a.Inspect() != "[1, 2, 3]" {
		t.Errorf("array.Inspect() wrong, expected %q, got %q",
			"[1, 2, 3]", a.Inspect())
	}

	if a.Hashable() {
		t.Errorf("array.Hashable() is true")
	}

	if a.HashKey() != nil {
		t.Errorf("array.HashKey() is not nil")
	}
}

func TestArrayObjectEvaluation(t *testing.T) {
	a := NewArray([]Object{
		NewInteger(1),
		NewInteger(2),
		NewInteger(3),
	})

	tests := []testObjectEvaluationCase{
		evalTest("!ARRAY([1, 2, 3])").
			call(a.OnPrefix(token.Bang)).
			expect(NewBoolean(false), true),
		evalTest("-ARRAY([1, 2, 3])").
			call(a.OnPrefix(token.Minus)).
			expect(nil, false),
		evalTest("ARRAY[1, 2, 3] == ARRAY[1, 2, 3]").
			call(a.OnInfix(token.EQ, a)).
			expect(NewBoolean(true), true),
		evalTest("ARRAY[1, 2, 3] != ARRAY[1, 2, 3]").
			call(a.OnInfix(token.NE, a)).
			expect(NewBoolean(false), true),
		evalTest("ARRAY[1, 2, 3] == ARRAY[1, 2, 3]").
			call(a.OnInfix(token.EQ, NewArray([]Object{
				NewInteger(1),
				NewInteger(2),
				NewInteger(3),
			}))).
			expect(NewBoolean(true), true),
		evalTest("ARRAY[1, 2, 3] == ARRAY[1, 2]").
			call(a.OnInfix(token.EQ, NewArray([]Object{
				NewInteger(1),
				NewInteger(2),
			}))).
			expect(NewBoolean(false), true),
		evalTest("ARRAY[1, 2, 3] == ARRAY[1, 2, 4]").
			call(a.OnInfix(token.EQ, NewArray([]Object{
				NewInteger(1),
				NewInteger(2),
				NewInteger(4),
			}))).
			expect(NewBoolean(false), true),
		evalTest("ARRAY[1, 2, 3] == INTEGER(42)").
			call(a.OnInfix(token.EQ, NewInteger(42))).
			expect(NewBoolean(false), true),
		evalTest("ARRAY[1, 2, 3] + INTEGER(42)").
			call(a.OnInfix(token.Plus, NewInteger(42))).
			expect(nil, false),
	}

	testObjectEvaluation(t, tests)
}

func TestArrayOnIndexEvalutation(t *testing.T) {
	a := NewArray([]Object{
		NewInteger(1),
		NewInteger(2),
		NewInteger(3),
	})

	tests := []testObjectEvaluationCase{
		evalTest("ARRAY[1, 2, 3][0]").
			call(a.OnIndex(NewInteger(0))).
			expect(NewInteger(1), true),
		evalTest("ARRAY[1, 2, 3][1]").
			call(a.OnIndex(NewInteger(1))).
			expect(NewInteger(2), true),
		evalTest("ARRAY[1, 2, 3][2]").
			call(a.OnIndex(NewInteger(2))).
			expect(NewInteger(3), true),
		evalTest("ARRAY[1, 2, 3][3]").
			call(a.OnIndex(NewInteger(3))).
			expect(NewNull(), true),
		evalTest("ARRAY[1, 2, 3][-1]").
			call(a.OnIndex(NewInteger(-1))).
			expect(NewInteger(3), true),
		evalTest("ARRAY[1, 2, 3][-3]").
			call(a.OnIndex(NewInteger(-3))).
			expect(NewInteger(1), true),
		evalTest("ARRAY[1, 2, 3][-5]").
			call(a.OnIndex(NewInteger(-5))).
			expect(NewNull(), true),
	}

	testObjectEvaluation(t, tests)
}
