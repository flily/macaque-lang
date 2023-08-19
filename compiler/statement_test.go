package compiler

import (
	"testing"

	"github.com/flily/macaque-lang/opcode"
)

func TestCompileLetStatement(t *testing.T) {
	tests := []testCompilerCase{
		{
			`let a = 42;`,
			code(
				inst(opcode.ILoadInt, 42),
				inst(opcode.ISStore, 1),
			),
			data(),
		},
		{
			`let a, b, c, d = 42, 43, 44, 45;`,
			code(
				inst(opcode.ILoadInt, 42),
				inst(opcode.ILoadInt, 43),
				inst(opcode.ILoadInt, 44),
				inst(opcode.ILoadInt, 45),
				inst(opcode.ISStore, 4),
				inst(opcode.ISStore, 3),
				inst(opcode.ISStore, 2),
				inst(opcode.ISStore, 1),
			),
			data(),
		},
		{
			text(
				"let a = 1;",
				"let b, c, d = 2, 3, 4, 5;",
				"a + d;",
			),
			code(
				inst(opcode.ILoadInt, 1),
				inst(opcode.ISStore, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 3),
				inst(opcode.ILoadInt, 4),
				inst(opcode.ILoadInt, 5),
				inst(opcode.IPop, 1),
				inst(opcode.ISStore, 4),
				inst(opcode.ISStore, 3),
				inst(opcode.ISStore, 2),
				inst(opcode.IClean),
				inst(opcode.ISLoad, 1),
				inst(opcode.ISLoad, 4),
				inst(opcode.IBinOp, 23),
			),
			data(),
		},
		{
			text(
				"7;",
				"let a, b, c, d = 1, 2, 3;",
				"a + d;",
			),
			code(
				inst(opcode.ILoadInt, 7),
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 3),
				inst(opcode.ISStore, 3),
				inst(opcode.ISStore, 2),
				inst(opcode.ISStore, 1),
				inst(opcode.IClean),
				inst(opcode.ISLoad, 1),
				inst(opcode.ISLoad, 4),
				inst(opcode.IBinOp, 23),
			),
			data(),
		},
		{
			text(
				"let f = fn() { 42, 43, 44, 45 };",
				"let a, b, c = f();",
			),
			code(
				inst(opcode.IMakeFunc, 1, 0),
				inst(opcode.ISStore, 1),
				inst(opcode.IClean),
				inst(opcode.IScopeIn),
				inst(opcode.ISLoad, 1),
				inst(opcode.ICall, 0),
				inst(opcode.IStackRev),
				inst(opcode.ISStore, 2),
				inst(opcode.ISStore, 3),
				inst(opcode.ISStore, 4),
				inst(opcode.IScopeOut),
				inst(opcode.IHalt),
				inst(opcode.IScopeIn),
				inst(opcode.ILoadInt, 42),
				inst(opcode.ILoadInt, 43),
				inst(opcode.ILoadInt, 44),
				inst(opcode.ILoadInt, 45),
				inst(opcode.IReturn),
				inst(opcode.IHalt),
			),
			data(),
		},
		{
			text(
				"let f = fn() { 5, 7, 9 };",
				"let a, b, c = f(), 3;",
				"a, b, c;",
			),
			code(
				inst(opcode.IMakeFunc, 1, 0),
				inst(opcode.ISStore, 1),
				inst(opcode.IScopeIn),
				inst(opcode.ISLoad, 1),
				inst(opcode.ICall, 0),
				inst(opcode.ILoadInt, 3),
				inst(opcode.IStackRev),
				inst(opcode.ISStore, 2),
				inst(opcode.ISStore, 3),
				inst(opcode.ISStore, 4),
				inst(opcode.IScopeOut),
				inst(opcode.IClean),
				inst(opcode.ISLoad, 2),
				inst(opcode.ISLoad, 3),
				inst(opcode.ISLoad, 4),
				inst(opcode.IHalt),
				inst(opcode.IScopeIn),
				inst(opcode.ILoadInt, 5),
				inst(opcode.ILoadInt, 7),
				inst(opcode.ILoadInt, 9),
				inst(opcode.IReturn),
				inst(opcode.IHalt),
			),
			data(),
		},
		{
			text(
				"let n = 42;",
				"let a, b, c  = 3, 5, if (n > 5) { 9; };",
				"a, b, c;",
			),
			code(
				inst(opcode.ILoadInt, 42),
				inst(opcode.ISStore, 1),
				inst(opcode.ILoadInt, 3),
				inst(opcode.ILoadInt, 5),
				inst(opcode.IScopeIn),
				inst(opcode.ISLoad, 1),
				inst(opcode.ILoadInt, 5),
				inst(opcode.IBinOp, 30),
				inst(opcode.IScopeOut, 1),
				inst(opcode.IJumpIf, 4),
				inst(opcode.IScopeIn),
				inst(opcode.ILoadInt, 9),
				inst(opcode.IScopeOut, 0),
				inst(opcode.IJumpFWD, 1),
				inst(opcode.ILoadNull),
				inst(opcode.ISStore, 4),
				inst(opcode.ISStore, 3),
				inst(opcode.ISStore, 2),
				inst(opcode.IClean),
				inst(opcode.ISLoad, 2),
				inst(opcode.ISLoad, 3),
				inst(opcode.ISLoad, 4),
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}

func TestCompileLetStatementRedefined(t *testing.T) {
	tests := []testCompilerErrorCase{
		{
			text(
				`let a = 42; let a = 43;`,
			),
			text(
				"let a = 42; let a = 43;",
				"                ^",
				"                variable a redeclared",
				"  at testcase:1:17",
				"let a = 42; let a = 43;",
				"    ^",
				"    variable a is already declared here",
				"  at testcase:1:5",
			),
		},
	}

	runCompilerErrorTestCases(t, tests)
}

func TestCompileFunctionsWithoutReturnValue(t *testing.T) {
	tests := []testCompilerCase{
		{
			`
				let noReturn = fn() { };
				noReturn();
			`,
			code(
				inst(opcode.IMakeFunc, 1, 0),
				inst(opcode.ISStore, 1),
				inst(opcode.IClean),
				inst(opcode.ISLoad, 1),
				inst(opcode.ICall, 0),
				inst(opcode.IHalt),
				inst(opcode.ILoadNull),
				inst(opcode.IReturn),
				inst(opcode.IHalt),
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}
