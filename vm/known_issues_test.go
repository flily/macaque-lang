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
			assertRegister(sp(1), bp(0)),
		},
		{
			text(
				"let a, b, c = 1, 2;",
				"b + c;",
			),
			stack(object.NewInteger(3)),
			assertRegister(sp(1), bp(0)),
		},
		{
			text(
				"let a, b, c = 1, 2, 3, 4;",
				"a + c;",
			),
			stack(object.NewInteger(6)),
			assertRegister(sp(1), bp(0)),
		},
		{
			text(
				"let a = 1;",
				"let b, c, d = 2, 3, 4, 5;",
				"a + d;",
			),
			stack(object.NewInteger(6)),
			assertRegister(sp(1), bp(0)),
		},
		{
			text(
				"7;",
				"let a, b, c, d = 1, 2, 3;",
				"a + d;",
			),
			stack(object.NewInteger(10)),
			assertRegister(sp(1), bp(0)),
		},
	}

	runVMTest(t, tests)
}

func TestReturnStatementInTheMiddle(t *testing.T) {
	// FIXED
	tests := []vmTest{
		{
			text(
				"let f = fn() {",
				"	5;",
				"	return 2;",
				"	3;",
				"};",
				"let a, b = f();",
				"a, b;",
			),
			stack(object.NewInteger(2), null),
			assertRegister(sp(2), bp(0)),
		},
	}

	runVMTest(t, tests)
}

func TestDirtyStackAssignment(t *testing.T) {
	// FIXED
	tests := []vmTest{
		{
			text(
				"let n = 42;",
				"let a, b, c  = 3, 5, if (n > 5) { 9; };",
				"a, b, c;",
			),
			// stack(object.NewInteger(9), object.NewNull(), object.NewNull()),
			stack(object.NewInteger(9), object.NewInteger(5), object.NewInteger(3)),
			assertRegister(sp(3), bp(0)),
		},
	}

	runVMTest(t, tests)
}
