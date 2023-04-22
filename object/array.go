package object

import (
	"strings"

	"github.com/flily/macaque-lang/token"
)

type ArrayObject struct {
	Elements []Object
}

func NewArray(elements []Object) Object {
	o := &ArrayObject{
		Elements: elements,
	}

	return o
}

func (a *ArrayObject) Type() ObjectType {
	return ObjectTypeArray
}

func (a *ArrayObject) Inspect() string {
	parts := make([]string, len(a.Elements))

	for i, e := range a.Elements {
		parts[i] = e.Inspect()
	}

	return "[" + strings.Join(parts, ", ") + "]"
}

func (a *ArrayObject) Hashable() bool {
	return false
}

func (a *ArrayObject) HashKey() interface{} {
	return nil
}

func (a *ArrayObject) EqualTo(o Object) bool {
	switch v := o.(type) {
	case *ArrayObject:
		if len(a.Elements) != len(v.Elements) {
			return false
		}

		for i, e := range a.Elements {
			if !e.EqualTo(v.Elements[i]) {
				return false
			}
		}

		return true
	}

	return false
}

func (a *ArrayObject) OnPrefix(t token.Token) (Object, bool) {
	var r Object
	ok := false
	switch t {
	case token.Bang:
		r, ok = NewBoolean(false), true
	}

	return r, ok
}

func (a *ArrayObject) OnInfix(t token.Token, o Object) (Object, bool) {
	if t == token.EQ || t == token.NE {
		return doEqualCompare(t, a.EqualTo(o))
	}

	var r Object
	ok := false
	switch t {
	}

	return r, ok
}

func (a *ArrayObject) OnIndex(o Object) (Object, bool) {
	var r Object
	ok := false
	switch v := o.(type) {
	case *IntegerObject:
		ok = true
		l := int64(len(a.Elements))
		switch {
		case -l <= v.Value && v.Value < 0:
			r = a.Elements[l+v.Value]
		case 0 <= v.Value && v.Value < l:
			r = a.Elements[v.Value]
		default:
			r = NewNull()
		}
	}

	return r, ok
}
