package compiler

import (
	"testing"

	"github.com/flily/macaque-lang/opcode"
	"github.com/flily/macaque-lang/token"
)

func TestParseListLiteral(t *testing.T) {
	tests := []testCompilerCase{
		{
			text(
				"[]",
			),
			code(
				inst(opcode.IMakeList, 0),
				inst(opcode.IPop, 1),
			),
			data(),
		},
		{
			text(
				"[1, 2, 3, 4]",
			),
			code(
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 3),
				inst(opcode.ILoadInt, 4),
				inst(opcode.IMakeList, 4),
				inst(opcode.IPop, 1),
			),
			data(),
		},
		{
			text(
				"[1, 2 + 3, 4,]",
			),
			code(
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 3),
				inst(opcode.IBinOp, int(token.Plus)),
				inst(opcode.ILoadInt, 4),
				inst(opcode.IMakeList, 3),
				inst(opcode.IPop, 1),
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}
