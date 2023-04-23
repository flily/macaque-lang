package object

import (
	"github.com/flily/macaque-lang/token"
)

type StringObject struct {
	Value string
}

func NewString(value string) Object {
	o := &StringObject{
		Value: value,
	}

	return o
}

func (s *StringObject) Type() ObjectType {
	return ObjectTypeString
}

func (s *StringObject) Inspect() string {
	return s.Value
}

func (s *StringObject) Hashable() bool {
	return true
}

func (s *StringObject) HashKey() interface{} {
	return s.Value
}

func (s *StringObject) EqualTo(o Object) bool {
	switch v := o.(type) {
	case *StringObject:
		return s.Value == v.Value
	}

	return false
}

func (s *StringObject) OnPrefix(t token.Token) (Object, bool) {
	var r Object
	ok := false
	switch t {
	case token.Bang:
		r, ok = NewBoolean(false), true
	}

	return r, ok
}

func (s *StringObject) OnInfix(t token.Token, o Object) (Object, bool) {
	if t == token.EQ || t == token.NE {
		return doEqualCompare(t, s.EqualTo(o))
	}

	var r Object
	ok := false

	switch v := o.(type) {
	case *StringObject:
		switch t {
		case token.Plus:
			r, ok = NewString(s.Value+v.Value), true
		}
	}

	return r, ok
}

func (s *StringObject) OnIndex(o Object) (Object, bool) {
	var r Object
	ok := false
	switch v := o.(type) {
	case *IntegerObject:
		if v.Value < 0 || v.Value >= int64(len(s.Value)) {
			r = objectNull
		} else {
			r = NewString(string(s.Value[v.Value]))
		}
		ok = true
	}

	return r, ok
}
