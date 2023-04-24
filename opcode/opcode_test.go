package opcode

import (
	"testing"
)

func TestCode(t *testing.T) {
	i := Code(IInvalid, 1, 2)
	if i.Name != IInvalid {
		t.Errorf("i.Name is not IInvalid, got %v", i.Name)
	}

	if i.Operand0 != 1 {
		t.Errorf("i.Operand0 is not 1, got %v", i.Operand0)
	}

	if i.Operand1 != 2 {
		t.Errorf("i.Operand1 is not 2, got %v", i.Operand1)
	}
}

func TestCodeName(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{-1, "INVALID"},
		{IInvalid, "INVALID"},
		{ILoadInt, "LOADINT"},
		{10000, "INVALID"},
	}

	for _, tt := range tests {
		if CodeName(tt.code) != tt.expected {
			t.Errorf("CodeName(%d) is not %s, got %s", tt.code, tt.expected, CodeName(tt.code))
		}
	}
}

func TestCodeString(t *testing.T) {
	code := Code(IInvalid, 1, 2)
	expected := "INVALID 1 2"
	if code.String() != expected {
		t.Errorf("code.String() is not %s, got %s", expected, code.String())
	}
}
