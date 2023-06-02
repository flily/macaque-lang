package compiler

import (
	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/opcode"
)

type CodePage struct {
	Codes     []opcode.Opcode
	Data      []object.Object
	Functions []*opcode.Function
}

func (p *CodePage) Main() *opcode.Function {
	return p.Functions[0]
}

func (p *CodePage) Info(i int) *opcode.Function {
	return p.Functions[i]
}
