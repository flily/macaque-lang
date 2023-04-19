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
