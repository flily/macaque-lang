package vm

import (
	"testing"

	"github.com/flily/macaque-lang/object"
)

func TestExpressions(t *testing.T) {
	tests := []vmTest{
		{
			`1 + 2`,
			stack(object.NewInteger(3)),
			assertRegister(sp(1)),
		},
		{
			`"hello" + " " + "world"`,
			stack(object.NewString("hello world")),
			assertRegister(sp(1)),
		},
	}

	runVMTest(t, tests)
}

func TestExpressionList(t *testing.T) {
	tests := []vmTest{
		{
			`1, 2, 3, 4`,
			stack(
				object.NewInteger(4),
				object.NewInteger(3),
				object.NewInteger(2),
				object.NewInteger(1),
			),
			assertRegister(sp(4)),
		},
	}

	runVMTest(t, tests)
}

func TestIfExpression(t *testing.T) {
	tests := []vmTest{
		{
			`if (true) { 10 }`,
			stack(object.NewInteger(10)),
			assertRegister(sp(1)),
		},
		{
			`if (true) { 10 } else { 20 }`,
			stack(object.NewInteger(10)),
			assertRegister(sp(1)),
		},
		{
			`if (false) { 10 } else { 20 }`,
			stack(object.NewInteger(20)),
			assertRegister(sp(1)),
		},
	}

	runVMTest(t, tests)
}
