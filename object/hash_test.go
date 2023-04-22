package object

import (
	"testing"

	"github.com/flily/macaque-lang/token"
)

func TestHashObject(t *testing.T) {
	h := NewHash([]HashPair{
		{NewString("one"), NewInteger(1)},
		{NewString("two"), NewInteger(2)},
		{NewString("three"), NewInteger(3)},
	})

	if h.Type() != ObjectTypeHash {
		t.Errorf("hash.Type() is not ObjectTypeHash")
	}

	if h.Inspect() != "{one: 1, two: 2, three: 3}" {
		t.Errorf("hash.Inspect() wrong, expected %q, got %q",
			"{one: 1, two: 2, three: 3}", h.Inspect())
	}

	if h.Hashable() {
		t.Errorf("hash.Hashable() is true")
	}

	if h.HashKey() != nil {
		t.Errorf("hash.HashKey() is not nil")
	}
}

func TestHashObjectEvaluation(t *testing.T) {
	h := NewHash([]HashPair{
		{NewString("one"), NewInteger(1)},
		{NewString("two"), NewInteger(2)},
		{NewString("three"), NewInteger(3)},
	})

	tests := []testObjectEvaluationCase{
		evalTest("!HASH({one: 1, two: 2, three: 3})").
			call(h.OnPrefix(token.Bang)).
			expect(NewBoolean(false), true),
		evalTest("-HASH({one: 1, two: 2, three: 3})").
			call(h.OnPrefix(token.Minus)).
			expect(nil, false),
		evalTest("HASH{one: 1, two: 2, three: 3} == HASH{one: 1, two: 2, three: 3}").
			call(h.OnInfix(token.EQ, h)).
			expect(NewBoolean(true), true),
		evalTest("HASH{one: 1, two: 2, three: 3} != HASH{one: 1, two: 2, three: 3}").
			call(h.OnInfix(token.NE, h)).
			expect(NewBoolean(false), true),
		evalTest("HASH{one: 1, two: 2, three: 3} == HASH{one: 1, two: 2, three: 3}").
			call(h.OnInfix(token.EQ, NewHash([]HashPair{
				{NewString("one"), NewInteger(1)},
				{NewString("two"), NewInteger(2)},
				{NewString("three"), NewInteger(3)},
			}))).
			expect(NewBoolean(true), true),
		evalTest("HASH{one: 1, two: 2, three: 3} == HASH{one: 1, two: 2}").
			call(h.OnInfix(token.EQ, NewHash([]HashPair{
				{NewString("one"), NewInteger(1)},
				{NewString("two"), NewInteger(2)},
			}))).
			expect(NewBoolean(false), true),
		evalTest("HASH{one: 1, two: 2, three: 3} == HASH{one: 1, two: 2}").
			call(h.OnInfix(token.EQ, NewHash([]HashPair{
				{NewString("one"), NewInteger(1)},
				{NewString("two"), NewInteger(2)},
				{NewString("three"), NewInteger(33)},
			}))).
			expect(NewBoolean(false), true),
		evalTest("HASH{one: 1, two: 2, three: 3} == 42").
			call(h.OnInfix(token.EQ, NewInteger(42))).
			expect(NewBoolean(false), true),
		evalTest("HASH{one: 1, two: 2, three: 3} + 42").
			call(h.OnInfix(token.Plus, NewInteger(42))).
			expect(nil, false),
	}

	testObjectEvaluation(t, tests)
}

func TestHashOnIndexEvaluation(t *testing.T) {
	h := NewHash([]HashPair{
		{NewString("one"), NewInteger(1)},
		{NewString("two"), NewInteger(2)},
		{NewString("three"), NewNull()},
	})

	tests := []testObjectEvaluationCase{
		evalTest(`HASH{"one": 1, "two": 2, "three": null}["one"]`).
			call(h.OnIndex(NewString("one"))).
			expect(NewInteger(1), true),
		evalTest(`HASH{"one": 1, "two": 2, "three": null}["two"]`).
			call(h.OnIndex(NewString("two"))).
			expect(NewInteger(2), true),
		evalTest(`HASH{"one": 1, "two": 2, "three": null}["three"]`).
			call(h.OnIndex(NewString("three"))).
			expect(NewNull(), true),
		evalTest(`HASH{"one": 1, "two": 2, "three": null}["four"]`).
			call(h.OnIndex(NewString("four"))).
			expect(NewNull(), true),
		evalTest(`HASH{"one": 1, "two": 2, "three": null}[true]`).
			call(h.OnIndex(NewBoolean(true))).
			expect(NewNull(), true),
		evalTest(`HASH{"one": 1, "two": 2, "three": null}[false]`).
			call(h.OnIndex(NewBoolean(false))).
			expect(NewNull(), true),
		evalTest(`HASH{"one": 1, "two": 2, "three": null}[null]`).
			call(h.OnIndex(NewNull())).
			expect(nil, false),
	}

	testObjectEvaluation(t, tests)
}
