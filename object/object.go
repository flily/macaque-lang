package object

import (
	"github.com/flily/macaque-lang/token"
)

type ObjectType int

const (
	ObjectTypeInvalid    ObjectType = iota
	ObjectTypeNull       ObjectType = 1
	ObjectTypeBoolean    ObjectType = 2
	ObjectTypeInteger    ObjectType = 3
	ObjectTypeFloat      ObjectType = 4
	ObjectTypeString     ObjectType = 5
	ObjectTypeArray      ObjectType = 6
	ObjectTypeHash       ObjectType = 7
	ObjectTypeFunction   ObjectType = 8
	ObjectTypeSystemFlag ObjectType = 64
)

type Object interface {
	Type() ObjectType
	Inspect() string
	Hashable() bool
	HashKey() interface{}
	EqualTo(Object) bool
	OnPrefix(token.Token) (Object, bool)
	OnInfix(token.Token, Object) (Object, bool)
}

type (
	UnaryFunction  func(op token.Token, self Object) (result Object, ok bool)
	BinaryFunction func(op token.Token, self Object, target Object) (result Object, ok bool)
)

var (
	objectNull  = &NullObject{}
	objectTrue  = &BooleanObject{Value: true}
	objectFalse = &BooleanObject{Value: false}
)

type NullObject struct {
}

func NewNull() Object {
	return objectNull
}

func (n *NullObject) Type() ObjectType {
	return ObjectTypeNull
}

func (n *NullObject) Inspect() string {
	return token.SNull
}

func (n *NullObject) Hashable() bool {
	return false
}

func (n *NullObject) HashKey() interface{} {
	return nil
}

func (n *NullObject) EqualTo(o Object) bool {
	switch o.(type) {
	case *NullObject:
		return true
	}

	return false
}

func (n *NullObject) OnPrefix(t token.Token) (Object, bool) {
	var r Object
	ok := false

	switch t {
	case token.Bang:
		r, ok = objectTrue, true
	}

	return r, ok
}

func (n *NullObject) OnInfix(t token.Token, o Object) (Object, bool) {
	if t == token.EQ || t == token.NE {
		return doEqualCompare(t, n.EqualTo(o))
	}

	return nil, false
}

type BooleanObject struct {
	Value bool
}

func NewBoolean(value bool) Object {
	if value {
		return objectTrue
	}

	return objectFalse
}

func (b *BooleanObject) Type() ObjectType {
	return ObjectTypeBoolean
}

func (b *BooleanObject) Inspect() string {
	if b.Value {
		return token.STrue
	}

	return token.SFalse
}

func (b *BooleanObject) Hashable() bool {
	return true
}

func (b *BooleanObject) HashKey() interface{} {
	return b.Value
}

func (b *BooleanObject) EqualTo(o Object) bool {
	switch v := o.(type) {
	case *BooleanObject:
		return b.Value == v.Value
	}

	return false
}

func (b *BooleanObject) OnPrefix(t token.Token) (Object, bool) {
	var r Object
	ok := false
	switch t {
	case token.Bang:
		r, ok = NewBoolean(!b.Value), true
	}

	return r, ok
}

func (b *BooleanObject) OnInfix(t token.Token, o Object) (Object, bool) {
	if t == token.EQ || t == token.NE {
		return doEqualCompare(t, b.EqualTo(o))
	}

	return nil, false
}

func doEqualCompare(t token.Token, equal bool) (Object, bool) {
	if t == token.EQ {
		return NewBoolean(equal), true
	}

	return NewBoolean(!equal), true
}
