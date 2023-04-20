package vm

import (
	"testing"

	"github.com/flily/macaque-lang/object"
)

func TestLetStatement(t *testing.T) {
	tests := []vmTest{
		{
			`let a = 1`,
			stack(),
			assertRegister(sp(2), bp(0)),
		},
		{
			`let a = 1; let b = 2`,
			stack(),
			assertRegister(sp(3), bp(0)),
		},
		{
			text(
				`let answer = 30 + 6`,
				`let final_answer = answer + 6`,
				"final_answer",
			),
			stack(object.NewInteger(42)),
			assertRegister(sp(4), bp(0)),
		},
	}

	runVMTest(t, tests)
}
