package compiler

import (
	"testing"

	"github.com/flily/macaque-lang/opcode"
)

func TestCompileLetStatement(t *testing.T) {
	tests := []testCompilerCase{
		{
			`let a = 42`,
			code(
				inst(opcode.ILoadInt, 42),
				inst(opcode.ISetAX, 1),
				inst(opcode.ISStore, 1),
			),
			data(),
		},
		{
			`let a, b, c, d = 42, 43, 44, 45`,
			code(
				inst(opcode.ILoadInt, 42),
				inst(opcode.ILoadInt, 43),
				inst(opcode.ILoadInt, 44),
				inst(opcode.ILoadInt, 45),
				inst(opcode.ISetAX, 4),
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
			[]string{
				`let a = 42; let a = 43`,
			},
			[]string{
				`let a = 42; let a = 43`,
				"                ^",
				"                variable a redefined",
				"  at testcase:1:17",
			},
		},
	}

	runCompilerErrorTestCases(t, tests)
}
