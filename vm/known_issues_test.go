package vm

import (
	"testing"

	"github.com/flily/macaque-lang/object"
)

func TestParallelLetStatementWithParameterNumberNotMatched(t *testing.T) {
	tests := []vmTest{
		{
			text(
				"let a, b, c = 1, 2, 3;",
				"a + c;",
			),
			stack(object.NewInteger(4)),
			assertRegister(sp(5), bp(0)),
		},
		{
			text(

				"let a, b, c = 1, 2;",
				"b + c;",
			),
			stack(object.NewInteger(3)),
			assertRegister(sp(5), bp(0)),
		},
		{
			text(
				"let a, b, c = 1, 2, 3, 4;",
				"a + c;",
			),
			stack(object.NewInteger(6)),
			assertRegister(sp(6), bp(0)),
		},
		{
			text(
				"let a = 1;",
				"let b, c, d = 2, 3, 4, 5;",
				"a + d;",
			),
			stack(object.NewInteger(6)),
			assertRegister(sp(7), bp(0)),
		},
		{
			text(
				"7;",
				"let a, b, c, d = 1, 2, 3;",
				"a + d;",
			),
			stack(object.NewInteger(10)),
			assertRegister(sp(6), bp(0)),
		},
	}

	runVMTest(t, tests)
}

func TestReturnStatementInTheMiddle(t *testing.T) {
	tests := []vmTest{
		{
			text(
				"let f = fn() {",
				"	5;",
				"	return 2;",
				"	3;",
				"};",
				"let a, b = f();",
				"b + 7;",
			),
			stack(object.NewInteger(9)),
			assertRegister(sp(5), bp(0)),
		},
	}

	runVMTest(t, tests)
}
