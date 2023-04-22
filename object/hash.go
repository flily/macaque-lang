package object

import (
	"strings"

	"github.com/flily/macaque-lang/token"
)

type HashPair struct {
	Key   Object
	Value Object
}

type HashObject struct {
	Elements []HashPair
	Map      map[interface{}]HashPair
}

func NewHash(elements []HashPair) Object {
	o := &HashObject{
		Elements: elements,
		Map:      make(map[interface{}]HashPair),
	}

	for _, e := range elements {
		o.Map[e.Key.HashKey()] = e
	}

	return o
}

func (h *HashObject) Type() ObjectType {
	return ObjectTypeHash
}

func (h *HashObject) Inspect() string {
	parts := make([]string, len(h.Elements))

	for i, e := range h.Elements {
		parts[i] = e.Key.Inspect() + ": " + e.Value.Inspect()
	}

	return "{" + strings.Join(parts, ", ") + "}"
}

func (h *HashObject) Hashable() bool {
	return false
}

func (h *HashObject) HashKey() interface{} {
	return nil
}

func (h *HashObject) EqualTo(o Object) bool {
	switch v := o.(type) {
	case *HashObject:
		if len(h.Elements) != len(v.Elements) {
			return false
		}

		for _, e := range h.Elements {
			if !e.Value.EqualTo(v.Map[e.Key.HashKey()].Value) {
				return false
			}
		}

		return true
	}

	return false
}

func (h *HashObject) OnPrefix(t token.Token) (Object, bool) {
	var r Object
	ok := false
	switch t {
	case token.Bang:
		r, ok = NewBoolean(false), true
	}

	return r, ok
}

func (h *HashObject) OnInfix(t token.Token, o Object) (Object, bool) {
	if t == token.EQ || t == token.NE {
		return doEqualCompare(t, h.EqualTo(o))
	}

	var r Object
	ok := false
	switch t {
	}

	return r, ok
}

func (h *HashObject) OnIndex(o Object) (Object, bool) {
	if !o.Hashable() {
		return nil, false
	}

	if e, ok := h.Map[o.HashKey()]; ok {
		return e.Value, true
	}

	return NewNull(), true
}
