package vm

import (
	"testing"

	"github.com/flily/macaque-lang/object"
)

func TestParallelLetStatementWithParameterNumberNotMatched(t *testing.T) {
	// FIXED
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
				"a + b, c;",
			),
			stack(null, object.NewInteger(3)),
			assertRegister(sp(2), bp(0)),
		},
		{
			text(
				"let a, b, c = 1, 2, 3, 4;",
				"a + c;",
			),
			stack(object.NewInteger(4)),
			assertRegister(sp(1), bp(0)),
		},
		{
			text(
				"let a = 1;",
				"let b, c, d = 2, 3, 4, 5;",
				"a + d;",
			),
			stack(object.NewInteger(5)),
			assertRegister(sp(1), bp(0)),
		},
		{
			text(
				"7;",
				"let a, b, c, d = 1, 2, 3;",
				"a, b, c, d;",
			),
			stack(null, object.NewInteger(3), object.NewInteger(2), object.NewInteger(1)),
			assertRegister(sp(4), bp(0)),
		},
		{
			text(
				"let f = fn() { 5, 7, 9 };",
				"let a, b, c = f(), 3;",
				"a, b, c;",
			),
			stack(object.NewInteger(9), object.NewInteger(7), object.NewInteger(5)),
			assertRegister(sp(3), bp(0)),
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
			stack(null, object.NewInteger(2)),
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
