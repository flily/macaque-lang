package vm

import (
	"testing"
)

func TestLetStatement(t *testing.T) {
	tests := []vmTest{
		{
			`let a = 1`,
			stack(),
			assertRegister(sp(1), bp(0)),
		},
		{
			`let a = 1; let b = 2`,
			stack(),
			assertRegister(sp(1), bp(0)),
		},
	}

	runVMTest(t, tests)
}
