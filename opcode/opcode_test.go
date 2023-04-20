package opcode

import (
	"testing"
)

func TestCode(t *testing.T) {
	i := Code(IInvalid, 1, 2, 3)
	if i.Name != IInvalid {
		t.Errorf("i.Name is not IInvalid, got %v", i.Name)
	}

	if i.Operand0 != 1 {
		t.Errorf("i.Operand0 is not 1, got %v", i.Operand0)
	}

	if i.Operand1 != 2 {
		t.Errorf("i.Operand1 is not 2, got %v", i.Operand1)
	}

	if i.Operand2 != 3 {
		t.Errorf("i.Operand2 is not 3, got %v", i.Operand2)
	}
}
