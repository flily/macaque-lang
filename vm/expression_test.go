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
			assertRegister(sp(2)),
		},
		{
			`"hello" + " " + "world"`,
			stack(object.NewString("hello world")),
			assertRegister(sp(2)),
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
			assertRegister(sp(5)),
		},
	}

	runVMTest(t, tests)
}
