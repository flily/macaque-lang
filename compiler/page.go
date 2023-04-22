package compiler

import (
	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/opcode"
)

type FunctionInfo struct {
	Index     int
	FrameSize int
	Arguments int
	IP        uint64
}

func (f FunctionInfo) Func(bounds []object.Object) *object.FunctionObject {
	return object.NewFunction(f.FrameSize, f.Arguments, f.IP, bounds)
}

type CodePage struct {
	Codes     []opcode.Opcode
	Data      []object.Object
	Functions []FunctionInfo
}

func (p *CodePage) Main() FunctionInfo {
	return p.Functions[0]
}

func (p *CodePage) Info(i int) FunctionInfo {
	return p.Functions[i]
}
