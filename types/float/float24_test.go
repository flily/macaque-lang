package float

import (
	"fmt"
	"math"
	"testing"
)

func TestFloat24Zero(t *testing.T) {
	pz := PositiveZero24
	pz32 := PositiveZero32
	if pz.Float32() != pz32 {
		t.Errorf("+0.0.Float32() = %v, want %f", pz, pz32)
	}

	pz64 := float64(0.0)
	if pz.Float64() != pz64 {
		t.Errorf("+0.0.Float32() = %v, want %f", pz, pz64)
	}

	if pz.String() != "+0.000" {
		t.Errorf("+0.0.String() = %v, want %v", pz, "+0.000")
	}

	nz := NegativeZero24
	nz32 := NegativeZero32
	if nz.Float32() != nz32 {
		t.Errorf("-0.0.Float32() = %v, want %f", nz, nz32)
	}

	nz64 := NegativeZero64
	if nz.Float64() != nz64 {
		t.Errorf("-0.0.Float32() = %v, want %f", nz, nz64)
	}

	if nz.String() != "-0.000" {
		t.Errorf("-0.0.String() = %v, want %v", nz, "-0.000")
	}
}

func TestFloat24Infinity(t *testing.T) {
	pi := PositiveInfinity24

	pi32 := PositiveInfinity32
	if pi.Float32() != pi32 {
		t.Errorf("+Inf.Float32() = %v, want %f", pi, pi32)
	}

	pi64 := PositiveInfinity64
	if pi.Float64() != pi64 {
		t.Errorf("+Inf.Float64() = %v, want %f", pi, pi64)
	}

	if pi.String() != "+Inf" {
		t.Errorf("+Inf.String() = %v, want %v", pi, "+Inf")
	}

	if !pi.IsInf() {
		t.Errorf("+Inf.IsInf() = %v, want %v", pi, true)
	}

	ni := NegativeInfinity24
	ni32 := NegativeInfinity32
	if ni.Float32() != ni32 {
		t.Errorf("-Inf.Float32() = %v, want %f", ni, ni32)
	}

	ni64 := NegativeInfinity64
	if ni.Float64() != ni64 {
		t.Errorf("-Inf.Float64() = %v, want %f", ni, ni64)
	}

	if ni.String() != "-Inf" {
		t.Errorf("-Inf.String() = %v, want %v", ni, "-Inf")
	}

	if !ni.IsInf() {
		t.Errorf("-Inf.IsInf() = %v, want %v", ni, true)
	}
}

func TestFloat24NaN(t *testing.T) {
	nan64s := []float64{
		math.NaN(),
		buildFloat64(0, 1024, 0xffffffff),
	}

	for _, f := range nan64s {
		if !math.IsNaN(f) {
			t.Errorf("math.IsNaN(%v) = %v, want %v", f, math.IsNaN(f), true)
		}
	}

	nan24s := []Float24{
		NaN24,
		NewFloat24(math.NaN()),
		NewFloat24(buildFloat64(0, 1024, 0xffffffff)),
		NewFloat24(buildFloat64(0, 1024, 0xfff00000)),
	}

	for _, nan := range nan24s {
		if !nan.IsNaN() {
			t.Errorf("%s.IsNaN() = %v, want %v", nan, nan.IsInf(), true)
		}

		if nan.String() != "NaN" {
			t.Errorf("NaN.String() = %v, want %v", nan, "NaN")
		}
	}
}

func TestNewInfinity(t *testing.T) {
	bigNumber := float64(1e100)
	if math.IsInf(bigNumber, 1) {
		t.Errorf("math.IsInf(%v, 1) = %v, want %v", bigNumber, math.IsInf(bigNumber, 1), true)
	}

	if math.IsNaN(bigNumber) {
		t.Errorf("math.IsNaN(%v) = %v, want %v", bigNumber, math.IsNaN(bigNumber), false)
	}

	f24 := NewFloat24(bigNumber)
	if !f24.IsInf() {
		t.Errorf("NewFloat24(%v).IsInf() = %v, want %v", bigNumber, f24.IsInf(), true)
	}
}

func TestNewZero(t *testing.T) {
	smallNumber := float64(1e-100)

	if smallNumber < 0 {
		t.Errorf("%v < 0", smallNumber)
	}

	if !(smallNumber > 0) {
		t.Errorf("%v > 0", smallNumber)
	}

	f24 := NewFloat24(smallNumber)
	f64 := f24.Float64()
	if f64 < 0 {
		t.Errorf("%v < 0", f64)
	}

	if f64 > 0 {
		t.Errorf("%v > 0", f64)
	}
}

func TestFloat24String(t *testing.T) {
	tests := []struct {
		number   float64
		expected string
	}{
		{math.NaN(), "NaN"},
		{math.Inf(1), "+Inf"},
		{math.Inf(-1), "-Inf"},
		{0.0, "+0.000"},
		{NegativeZero64, "-0.000"},
		{1.0, "1.000000"},
		{99.5, "99.500000"},

		// Cases which lose precision.
		{3.3, "3.299988"},
		{99.2, "99.199707"},
		{9.80665, "9.806641"},
		{3.1415926, "3.141586"},
		{299792458, "299792384.000000"},
		{9192631770, "+Inf"},
	}

	for _, test := range tests {
		f24 := NewFloat24(test.number)
		if f24.String() != test.expected {
			t.Errorf("%v.String() = %v, want %v", test.number, f24.String(), test.expected)
		}
	}
}

func TestFloatString(t *testing.T) {
	tests := []struct {
		number   float64
		expected string
	}{
		{math.NaN(), "NaN"},
		{math.Inf(1), "+Inf"},
		{math.Inf(-1), "-Inf"},
		{0.0, "0.000000"},
		{NegativeZero64, "-0.000000"},
		{1.0, "1.000000"},
		{99.5, "99.500000"},

		{3.3, "3.300000"},
		{99.2, "99.200000"},
		{9.80665, "9.806650"},
		{3.1415926, "3.141593"},
		{299792458, "299792458.000000"},
		{9192631770, "9192631770.000000"},
	}

	for _, test := range tests {
		got := fmt.Sprintf("%f", test.number)
		if got != test.expected {
			t.Errorf("%v.String() = %v, want %v", test.number, got, test.expected)
		}
	}
}
