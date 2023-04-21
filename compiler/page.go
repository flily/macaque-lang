package compiler

import (
	"github.com/flily/macaque-lang/opcode"
)

type CodePage struct {
	Codes       []opcode.Opcode
	FunctionMap map[int]uint64
}
