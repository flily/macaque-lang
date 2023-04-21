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
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}

func TestCompileFunctions(t *testing.T) {
	tests := []testCompilerCase{
		{
			text(
				"let answer = fn() { 42; };",
			),
			code(
				inst(opcode.IMakeFunc, 0, 0),
				inst(opcode.ISStore, 1),
				inst(opcode.IHalt),
				inst(opcode.ILoadInt, 42),
				inst(opcode.IReturn, 1),
				inst(opcode.IHalt),
			),
			data(),
		},
		{
			text(
				"let answer = 42;",
				"let final_answer = fn() { answer; };",
				"final_answer",
			),
			code(
				inst(opcode.ILoadInt, 42),
				inst(opcode.ISStore, 1),
				inst(opcode.ISLoad, 1),
				inst(opcode.IMakeFunc, 0, 1),
				inst(opcode.ISStore, 2),
				inst(opcode.ISLoad, 2),
				inst(opcode.IHalt),
				inst(opcode.ILoadBind, 0),
				inst(opcode.IReturn, 1),
				inst(opcode.IHalt),
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}
