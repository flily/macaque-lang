package object

import (
	"testing"

	"github.com/flily/macaque-lang/token"
)

type testObjectEvaluationCase struct {
	desciption string
	gotObject  Object
	gotOk      bool
	expected   Object
	ok         bool
}

func (c testObjectEvaluationCase) call(got Object, ok bool) testObjectEvaluationCase {
	c.gotObject = got
	c.gotOk = ok
	return c
}

func (c testObjectEvaluationCase) expect(expected Object, ok bool) testObjectEvaluationCase {
	c.expected = expected
	c.ok = ok
	return c
}

func evalTest(s string) testObjectEvaluationCase {
	c := testObjectEvaluationCase{
		desciption: s,
	}

	return c
}

func testObjectEvaluation(t *testing.T, cases []testObjectEvaluationCase) {
	for _, c := range cases {
		if c.gotObject != c.expected {
			t.Errorf("%s: got %v, expected %v", c.desciption, c.gotObject, c.expected)
		}

		if c.gotOk != c.ok {
			t.Errorf("%s: got %v, expected %v", c.desciption, c.gotOk, c.ok)
		}
	}
}

func TestNullObject(t *testing.T) {
	null := NewNull()

	if null.Type() != ObjectTypeNull {
		t.Errorf("null.Type() is not ObjectTypeNull")
	}

	if null.Inspect() != token.SNull {
		t.Errorf("null.Inspect() is not '%s'", token.SNull)
	}

	if null.Hashable() {
		t.Errorf("null.Hashable() is not false")
	}

	if null.HashKey() != nil {
		t.Errorf("null.HashKey() is not nil")
	}

	if !null.EqualTo(NewNull()) {
		t.Errorf("null.EqualTo(NewNull()) is not true")
	}

	if null.EqualTo(NewBoolean(true)) {
		t.Errorf("null.EqualTo(NewBoolean(true)) is not false")
	}

	if r, ok := null.OnPrefix(token.Minus); r != nil || ok {
		t.Errorf("null.OnPrefix(token.Minus) is not (nil, false)")
	}

	if r, ok := null.OnPrefix(token.Bang); !r.EqualTo(NewBoolean(true)) || !ok {
		t.Errorf("null.OnPrefix(token.Bang) is not (true, true)")
	}

	if r, ok := null.OnInfix(token.Plus, NewNull()); r != nil || ok {
		t.Errorf("null.OnInfix(token.Plus, NewNull()) is not (nil, false)")
	}
}

func TestBooleanObject(t *testing.T) {
	vTrue, vFalse := NewBoolean(true), NewBoolean(false)

	if vTrue.Type() != ObjectTypeBoolean {
		t.Errorf("vTrue.Type() is not ObjectTypeBoolean")
	}

	if vTrue.Inspect() != token.STrue {
		t.Errorf("vTrue.Inspect() is not '%s'", token.STrue)
	}

	if vFalse.Inspect() != token.SFalse {
		t.Errorf("vFalse.Inspect() is not '%s'", token.SFalse)
	}

	if !vTrue.Hashable() || !vFalse.Hashable() {
		t.Errorf("vTrue.Hashable() = %v, vFalse.Hashable() = %v",
			vTrue.Hashable(), vFalse.Hashable())
	}

	if vTrue.HashKey() != true || vFalse.HashKey() != false {
		t.Errorf("vTrue.HashKey() = %v, vFalse.HashKey() = %v",
			vTrue.HashKey(), vFalse.HashKey())
	}

	if vTrue.EqualTo(vFalse) || !vTrue.EqualTo(vTrue) {
		t.Errorf("vTrue.EqualTo(vFalse) = %v, vTrue.EqualTo(vTrue) = %v",
			vTrue.EqualTo(vFalse), vTrue.EqualTo(vTrue))
	}

	if vFalse.EqualTo(vTrue) || !vFalse.EqualTo(vFalse) {
		t.Errorf("vFalse.EqualTo(vTrue) = %v, vFalse.EqualTo(vFalse) = %v",
			vFalse.EqualTo(vTrue), vFalse.EqualTo(vFalse))
	}

	if vTrue.EqualTo(NewNull()) || vFalse.EqualTo(NewNull()) {
		t.Errorf("vTrue.EqualTo(NewNull()) = %v, vFalse.EqualTo(NewNull()) = %v",
			vTrue.EqualTo(NewNull()), vFalse.EqualTo(NewNull()))
	}
}

func TestBooleanObjectEvaluation(t *testing.T) {
	tests := []testObjectEvaluationCase{
		evalTest("!true").
			call(objectTrue.OnPrefix(token.Bang)).
			expect(objectFalse, true),
		evalTest("!false").
			call(objectFalse.OnPrefix(token.Bang)).
			expect(objectTrue, true),
		evalTest("true + true").
			call(objectTrue.OnInfix(token.Plus, objectTrue)).
			expect(nil, false),
	}

	testObjectEvaluation(t, tests)
}
