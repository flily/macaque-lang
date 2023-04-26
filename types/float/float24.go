// Special encoding for 24-bit floating point numbers, used in 32-bit instruction
package float

import "fmt"

// IEEE-754 float number
//  Size   Sign   Exponent   Fraction
//   16      1         5          10
//   24      1         6          17
//   32      1         8          23
//   64      1        11          52
//  128      1        15         112

type Float24 uint32

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

func NewPositiveZero() Float24 {
	return buildFloat24(0, -31, 0)
}

func NewNegativeZero() Float24 {
	return buildFloat24(1, -31, 0)
}

func NewPositiveInfinite() Float24 {
	return buildFloat24(0, 32, 0)
}

func NewNegativeInfinite() Float24 {
	return buildFloat24(1, 32, 0)
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
	default:
		s = fmt.Sprintf("%f", buildFloat32(sign, exp, frac<<6))
	}

	return s
}

func NewFloat24FromFloat32(f float32) Float24 {
	sign, exp, frac := splitFloat32(f)
	return buildFloat24(sign, exp, frac>>6)
}
