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
			),
			data(),
		},
		{
			`42, "answer"`,
			code(
				inst(opcode.ILoadInt, 42),
				inst(opcode.ILoad, 0),
			),
			data(
				object.NewString("answer"),
			),
		},
	}

	runCompilerTestCases(t, tests)
}

func TestEvaluationExpressionWithVariables(t *testing.T) {
	tests := []testCompilerCase{
		{
			text(
				`let answer = 30 + 6`,
				`let final_anser = answer + 6`,
			),
			code(
				inst(opcode.ILoadInt, 30),
				inst(opcode.ILoadInt, 6),
				inst(opcode.IBinOp, int(token.Plus)),
				inst(opcode.ISStore, 1),
				inst(opcode.ISLoad, 1),
				inst(opcode.ILoadInt, 6),
				inst(opcode.IBinOp, int(token.Plus)),
				inst(opcode.ISStore, 2),
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}

func TestCompileIfExpression(t *testing.T) {
	tests := []testCompilerCase{
		{
			text(
				"if (true) { 10 }",
			),
			code(
				inst(opcode.ILoadBool, 1),
				inst(opcode.IJumpIf, 1),
				inst(opcode.ILoadInt, 10),
			),
			data(),
		},
		{
			text(
				"if (10 > 4) { 1, 2, 3, 4; 5 }",
			),
			code(
				inst(opcode.ILoadInt, 10),
				inst(opcode.ILoadInt, 4),
				inst(opcode.IBinOp, int(token.GT)),
				inst(opcode.IJumpIf, 6),
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 3),
				inst(opcode.ILoadInt, 4),
				inst(opcode.IPop, 4),
				inst(opcode.ILoadInt, 5),
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}
