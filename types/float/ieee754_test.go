package float

import (
	"testing"
)

func TestFloat64Convert(t *testing.T) {
	{
		pi := float64(3.1415926)
		sign, exp, frac := splitFloat64(pi)
		got := buildFloat32(sign, exp, frac>>29)
		expected := float32(pi)
		if got != expected {
			t.Errorf("buildFloat32(%v) = %v, want %v", pi, got, expected)
		}
	}

	{
		pi := float32(3.1415926)
		sign, exp, frac := splitFloat32(pi)
		got := buildFloat64(sign, exp, frac<<29)
		expected := float64(pi)
		if got != expected {
			t.Errorf("buildFloat64(%v) = %v, want %v", pi, got, expected)
		}
	}
}
