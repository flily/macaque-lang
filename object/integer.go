package object

import (
	"fmt"

	"github.com/flily/macaque-lang/token"
)

type IntegerObject struct {
	Value int64
}

func NewInteger(value int64) Object {
	o := &IntegerObject{
		Value: value,
	}

	return o
}

func (i *IntegerObject) Type() ObjectType {
	return ObjectTypeInteger
}

func (i *IntegerObject) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *IntegerObject) Hashable() bool {
	return true
}

func (i *IntegerObject) HashKey() interface{} {
	return i.Value
}

func (i *IntegerObject) EqualTo(o Object) bool {
	switch v := o.(type) {
	case *IntegerObject:
		return i.Value == v.Value
	}

	return false
}

func (i *IntegerObject) OnPrefix(t token.Token) (Object, bool) {
	var r Object
	ok := false

	switch t {
	case token.Bang:
		r, ok = NewBoolean(false), true

	case token.Minus:
		r, ok = NewInteger(-i.Value), true

	case token.BITNOT:
		r, ok = NewInteger(^i.Value), true
	}

	return r, ok
}

func (i *IntegerObject) OnInfix(t token.Token, o Object) (Object, bool) {
	var r Object
	ok := false

	if t == token.EQ || t == token.NE {
		return doEqualCompare(t, i.EqualTo(o))
	}

	switch v := o.(type) {
	case *IntegerObject:
		r, ok = i.onIntegerInfix(t, v)

	}

	return r, ok
}

func (i *IntegerObject) OnIndex(o Object) (Object, bool) {
	return nil, false
}

func (i *IntegerObject) onIntegerInfix(t token.Token, o *IntegerObject) (Object, bool) {
	var r Object
	ok := false
	switch t {
	case token.Plus:
		r, ok = NewInteger(i.Value+o.Value), true

	case token.Minus:
		r, ok = NewInteger(i.Value-o.Value), true

	case token.Asterisk:
		r, ok = NewInteger(i.Value*o.Value), true

	case token.Slash:
		// NOTES: will crash if o.Value == 0
		r, ok = NewInteger(i.Value/o.Value), true

	case token.Modulo:
		r, ok = NewInteger(i.Value%o.Value), true

	case token.LT:
		r, ok = NewBoolean(i.Value < o.Value), true

	case token.GT:
		r, ok = NewBoolean(i.Value > o.Value), true

	case token.LE:
		r, ok = NewBoolean(i.Value <= o.Value), true

	case token.GE:
		r, ok = NewBoolean(i.Value >= o.Value), true

	case token.AND:
		r, ok = NewBoolean(i.Value != 0 && o.Value != 0), true

	case token.OR:
		r, ok = NewBoolean(i.Value != 0 || o.Value != 0), true

	case token.BITAND:
		r, ok = NewInteger(i.Value&o.Value), true

	case token.BITOR:
		r, ok = NewInteger(i.Value|o.Value), true

	case token.BITXOR:
		r, ok = NewInteger(i.Value^o.Value), true
	}

	return r, ok
}
