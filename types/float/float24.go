// Special encoding for 24-bit floating point numbers, used in 32-bit instruction
package float

import (
	"fmt"
	"math"
)

// IEEE-754 float number
//  Size   Sign   Exponent   Fraction
//   16      1         5          10
//   24      1         6          17
//   32      1         8          23
//   64      1        11          52
//  128      1        15         112

type Float24 uint32

var (
	PositiveZero64     = float64(0.0)
	NegativeZero64     = buildFloat64(1, -1023, 0)
	PositiveInfinity64 = float64(math.Inf(1))
	NegativeInfinity64 = float64(math.Inf(-1))
	NaN24              = buildFloat24(0, 32, 1)

	PositiveZero32     = float32(0.0)
	NegativeZero32     = buildFloat32(1, -127, 0)
	PositiveInfinity32 = float32(math.Inf(1))
	NegativeInfinity32 = float32(math.Inf(-1))

	PositiveZero24     = buildFloat24(0, -31, 0)
	NegativeZero24     = buildFloat24(1, -31, 0)
	PositiveInfinity24 = buildFloat24(0, 32, 0)
	NegativeInfinity24 = buildFloat24(1, 32, 0)
)

func buildFloat24(sign uint32, exp int64, frac uint64) Float24 {
	u := (sign&1)<<23 |
		(uint32(exp+31)&0x3f)<<17 |
		(uint32(frac) & 0x1ffff)

	return Float24(u)
}

func splitFloat24(f Float24) (uint32, int64, uint64) {
	u := uint32(f)
	sign := uint32(u >> 23)
	exp := int64(u>>17&0x3f) - 31
	frac := uint64(u & 0x1ffff)
	return sign, exp, frac
}

func NewFloat24(f float64) Float24 {
	sign, exp, frac := splitFloat64(f)
	switch {
	case exp >= 1024 && frac != 0: // NaN
		if frac&0x1ffff == 0 {
			frac = 1
		}
		exp = 32

	case exp >= 32: // Infinity
		exp, frac = 32, 0

	case exp <= -31: // Zero
		exp, frac = -31, 0

	default:
		frac = frac >> 35
	}

	return buildFloat24(sign, exp, frac)
}

func (f Float24) String() string {
	sign, exp, frac := splitFloat24(f)
	var s string
	switch {
	case sign == 0 && exp == -31 && frac == 0:
		s = "+0.000"
	case sign == 1 && exp == -31 && frac == 0:
		s = "-0.000"
	case sign == 0 && exp == 32 && frac == 0:
		s = "+Inf"
	case sign == 1 && exp == 32 && frac == 0:
		s = "-Inf"
	case exp == 32 && frac != 0:
		s = "NaN"
	default:
		s = fmt.Sprintf("%f", buildFloat32(sign, exp, frac<<6))
	}

	return s
}

func (f Float24) Float32() float32 {
	sign, exp, frac := splitFloat24(f)
	if exp <= -31 {
		exp = -127
	} else if exp >= 32 {
		exp = 128
	}

	return buildFloat32(sign, exp, frac<<6)
}

func (f Float24) Float64() float64 {
	sign, exp, frac := splitFloat24(f)
	if exp <= -31 {
		exp = -1023
	} else if exp >= 32 {
		exp = 1024
	}

	return buildFloat64(sign, exp, frac<<35)
}

func (f Float24) IsNaN() bool {
	_, exp, frac := splitFloat24(f)
	return exp == 32 && frac != 0
}

func (f Float24) IsInf() bool {
	_, exp, frac := splitFloat24(f)
	return exp == 32 && frac == 0
}
