package object

import (
	"testing"

	"github.com/flily/macaque-lang/token"
)

func TestFunctionObject(t *testing.T) {
	f := NewFunction(0, 0, 42, nil)

	if f.Type() != ObjectTypeFunction {
		t.Errorf("f.Type() is not ObjectTypeFunction. got=%T (%+v)", f.Type(), f.Type())
	}

	if f.Inspect() != "function[0]" {
		t.Errorf("f.Inspect() wrong. got=%q", f.Inspect())
	}

	if f.Hashable() {
		t.Errorf("f.Hashable() is true")
	}

	if f.HashKey() != nil {
		t.Errorf("f.HashKey() is not nil")
	}
}

func TestFunctionObjectEvaluation(t *testing.T) {
	f1 := NewFunction(10, 2, 42, nil)
	f2 := NewFunction(10, 2, 42, nil)
	f3 := NewFunction(10, 2, 43, nil)
	f4 := NewFunction(10, 3, 42, []Object{NewInteger(1)})
	f5 := NewFunction(10, 3, 42, []Object{NewInteger(2)})

	tests := []testObjectEvaluationCase{
		evalTest("!FUNCTION(10, 2, 42)").
			call(f1.OnPrefix(token.Bang)).
			expect(NewBoolean(false), true),
		evalTest("-FUNCTION(10, 2, 42)").
			call(f1.OnPrefix(token.Minus)).
			expect(nil, false),
		evalTest("FUNCTION(10, 2, 42) == FUNCTION(10, 2, 42)").
			call(f1.OnInfix(token.EQ, f2)).
			expect(NewBoolean(true), true),
		evalTest("FUNCTION(10, 2, 42) != FUNCTION(10, 2, 42)").
			call(f1.OnInfix(token.NE, f2)).
			expect(NewBoolean(false), true),
		evalTest("FUNCTION(10, 2, 42) == FUNCTION(10, 2, 43)").
			call(f1.OnInfix(token.EQ, f3)).
			expect(NewBoolean(false), true),
		evalTest("FUNCTION(10, 2, 42) != FUNCTION(10, 2, 43)").
			call(f1.OnInfix(token.NE, f3)).
			expect(NewBoolean(true), true),
		evalTest("FUNCTION(10, 2, 42, [1]) == FUNCTION(10, 3, 42, [2])").
			call(f4.OnInfix(token.EQ, f5)).
			expect(NewBoolean(false), true),
		evalTest("FUNCTION(10, 2, 42, [1]) != FUNCTION(10, 3, 42, [2])").
			call(f4.OnInfix(token.NE, f5)).
			expect(NewBoolean(true), true),
		evalTest("FUNCTION(10, 2, 42) == INTEGER(42)").
			call(f1.OnInfix(token.EQ, NewInteger(42))).
			expect(NewBoolean(false), true),
		evalTest("FUNCTION(10, 2, 42) != INTEGER(42)").
			call(f1.OnInfix(token.NE, NewInteger(42))).
			expect(NewBoolean(true), true),
		evalTest("FUNCTION(10, 2, 42) + INTEGER(42)").
			call(f1.OnInfix(token.Plus, NewInteger(42))).
			expect(nil, false),
	}

	testObjectEvaluation(t, tests)
}

func TestFunctionObjectOnIndexEvaluation(t *testing.T) {
	f := NewFunction(10, 2, 42, nil)

	tests := []testObjectEvaluationCase{
		evalTest("FUNCTION(10, 2, 42)[0]").
			call(f.OnIndex(NewInteger(0))).
			expect(nil, false),
	}

	testObjectEvaluation(t, tests)
}
