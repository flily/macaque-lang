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
