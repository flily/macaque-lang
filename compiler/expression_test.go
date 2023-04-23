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
				`let answer = 30 + 6;`,
				`let final_anser = answer + 6;`,
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
				inst(opcode.IJumpIf, 2),
				inst(opcode.ILoadInt, 10),
				inst(opcode.IJumpFWD, 1),
				inst(opcode.ILoadNull),
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
				inst(opcode.IJumpIf, 7),
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 3),
				inst(opcode.ILoadInt, 4),
				inst(opcode.IClean),
				inst(opcode.ILoadInt, 5),
				inst(opcode.IJumpFWD, 1),
				inst(opcode.ILoadNull),
			),
			data(),
		},
		{
			text(
				"if (true) { 42 } else { 24 }",
			),
			code(
				inst(opcode.ILoadBool, 1),
				inst(opcode.IJumpIf, 2),
				inst(opcode.ILoadInt, 42),
				inst(opcode.IJumpFWD, 1),
				inst(opcode.ILoadInt, 24),
			),
			data(),
		},
		{
			text(
				"if (false) { 42 } else { 24 }",
			),
			code(
				inst(opcode.ILoadBool, 0),
				inst(opcode.IJumpIf, 2),
				inst(opcode.ILoadInt, 42),
				inst(opcode.IJumpFWD, 1),
				inst(opcode.ILoadInt, 24),
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}

func TestCompileIfExpressionError(t *testing.T) {
	tests := []testCompilerErrorCase{
		{
			text(
				"if (10 > 4) {",
				"    let a = 5;",
				"}",
				"4 + a",
			),
			text(
				"4 + a",
				"    ^",
				"    variable a undefined",
				"  at testcase:4:5",
			),
		},
	}

	runCompilerErrorTestCases(t, tests)
}

func TestCallExpression(t *testing.T) {
	tests := []testCompilerCase{
		{
			text(
				"let add = fn(a, b) { a + b };",
				"add(1, 2)",
			),
			code(
				inst(opcode.IMakeFunc, 1),
				inst(opcode.ISStore, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 1),
				inst(opcode.ISLoad, 1),
				inst(opcode.ICall, 2),
				inst(opcode.IHalt),
				inst(opcode.ISLoad, -1),
				inst(opcode.ISLoad, -2),
				inst(opcode.IBinOp, int(token.Plus)),
				inst(opcode.IReturn,),
				inst(opcode.IHalt),
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}

func TestIndexExpression(t *testing.T) {
	tests := []testCompilerCase{
		{
			text(
				"let arr = [1, 2, 3];",
				"arr[3]",
			),
			code(
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 3),
				inst(opcode.IMakeList, 3),
				inst(opcode.ISStore, 1),
				inst(opcode.ISLoad, 1),
				inst(opcode.ILoadInt, 3),
				inst(opcode.IIndex),
			),
			data(),
		},
		{
			text(
				`let hash = {"one": 1, "two": 2};`,
				`hash["one"]`,
			),
			code(
				inst(opcode.ILoad, 0),
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoad, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.IMakeHash, 2),
				inst(opcode.ISStore, 1),
				inst(opcode.ISLoad, 1),
				inst(opcode.ILoad, 0),
				inst(opcode.IIndex),
			),
			data(
				object.NewString("one"),
				object.NewString("two"),
			),
		},
	}

	runCompilerTestCases(t, tests)
}
