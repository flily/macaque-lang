package vm

import (
	"testing"

	"github.com/flily/macaque-lang/object"
)

func TestLetStatement(t *testing.T) {
	tests := []vmTest{
		{
			`let a = 1;`,
			stack(),
			assertRegister(sp(0), bp(0)),
		},
		{
			`let a = 1; let b = 2;`,
			stack(),
			assertRegister(sp(0), bp(0)),
		},
		{
			text(
				`let answer = 30 + 6;`,
				`let final_answer = answer + 6;`,
				"final_answer",
			),
			stack(object.NewInteger(42)),
			assertRegister(sp(1), bp(0)),
		},
	}

	runVMTest(t, tests)
}

// func TestReturnInIfExpression(t *testing.T) {
// 	tests := []vmTest{
// 		{
// 			text(
// 				`let r = if (10 > 1) {`,
// 				`	if (10 > 1) {`,
// 				`		return 10;`,
// 				`	}`,
// 				`	return 1;`,
// 				`};`,
// 				`r + 1;`,
// 			),
// 			stack(object.NewInteger(11)),
// 			assertRegister(sp(2), bp(0)),
// 		},
// 	}

// 	runVMTest(t, tests)
// }
