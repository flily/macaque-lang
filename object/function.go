package object

import (
	"github.com/flily/macaque-lang/token"
)

type FunctionObject struct {
	FrameSize int
	Arguments int
	IP        uint64
	Bounds    []Object
	Context   *token.Context
}

func NewFunction(frameSize int, args int, ip uint64, bounds []Object) *FunctionObject {
	f := &FunctionObject{
		FrameSize: frameSize,
		Arguments: args,
		IP:        ip,
		Bounds:    bounds,
	}

	return f
}

func (f *FunctionObject) Type() ObjectType {
	return ObjectTypeFunction
}

func (f *FunctionObject) Inspect() string {
	return "function"
}

func (f *FunctionObject) Hashable() bool {
	return false
}

func (f *FunctionObject) HashKey() interface{} {
	return nil
}

func (f *FunctionObject) equals(g *FunctionObject) bool {
	if f.IP != g.IP {
		return false
	}

	// Number of bounds must be equal
	for i := 0; i < len(f.Bounds); i++ {
		if !f.Bounds[i].EqualTo(g.Bounds[i]) {
			return false
		}
	}

	return true
}

func (f *FunctionObject) EqualTo(o Object) bool {
	switch o := o.(type) {
	case *FunctionObject:
		return f.equals(o)

	}

	return false
}

func (f *FunctionObject) OnPrefix(t token.Token) (Object, bool) {
	var r Object
	ok := false
	switch t {
	case token.Bang:
		r, ok = objectFalse, true
	}

	return r, ok
}

func (f *FunctionObject) OnInfix(t token.Token, o Object) (Object, bool) {
	if t == token.EQ || t == token.NE {
		return doEqualCompare(t, f.EqualTo(o))
	}

	return nil, false
}

func (f *FunctionObject) OnIndex(o Object) (Object, bool) {
	return nil, false
}
