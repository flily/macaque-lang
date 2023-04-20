package compiler

import (
	"testing"

	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/opcode"
	"github.com/flily/macaque-lang/token"
)

func TestCompileExpression(t *testing.T) {
	tests := []testCompilerCase{
		{
			`1 + 2`,
			code(
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.IBinOp, int(token.Plus)),
				inst(opcode.ISetAX, 1),
			),
			data(),
		},
		{
			`1 + 2 + 3`,
			code(
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.IBinOp, int(token.Plus)),
				inst(opcode.ILoadInt, 3),
				inst(opcode.IBinOp, int(token.Plus)),
				inst(opcode.ISetAX, 1),
			),
			data(),
		},
		{
			`1 + 2 * 3 + 4`,
			code(
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 3),
				inst(opcode.IBinOp, int(token.Asterisk)),
				inst(opcode.IBinOp, int(token.Plus)),
				inst(opcode.ILoadInt, 4),
				inst(opcode.IBinOp, int(token.Plus)),
				inst(opcode.ISetAX, 1),
			),
			data(),
		},
		{
			`1, 2 + 3, 4`,
			code(
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 3),
				inst(opcode.IBinOp, int(token.Plus)),
				inst(opcode.ILoadInt, 4),
				inst(opcode.ISetAX, 3),
			),
			data(),
		},
		{
			`42, "answer"`,
			code(
				inst(opcode.ILoadInt, 42),
				inst(opcode.ILoad, 0),
				inst(opcode.ISetAX, 2),
			),
			data(
				object.NewString("answer"),
			),
		},
	}

	runCompilerTestCases(t, tests)
}
