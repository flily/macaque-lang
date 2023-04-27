package float

import (
	"math"
)

func splitFloat32(f float32) (uint32, int64, uint64) {
	u := math.Float32bits(f)
	sign := uint32(u >> 31)
	exp := int64(u>>23&0xff) - 127
	frac := uint64(u & 0x7fffff)
	return sign, exp, frac
}

func buildFloat32(sign uint32, exp int64, frac uint64) float32 {
	u := (uint32(sign&1) << 31) |
		(uint32(exp+127))<<23 |
		(uint32(frac & 0x7fffff))

	return math.Float32frombits(u)
}

func splitFloat64(f float64) (uint32, int64, uint64) {
	u := math.Float64bits(f)
	sign := uint32(u >> 63)
	exp := int64(u>>52&0x7ff) - 1023
	frac := uint64(u & 0xfffffffffffff)
	return sign, exp, frac
}

func buildFloat64(sign uint32, exp int64, frac uint64) float64 {
	u := (uint64(sign)&1)<<63 |
		(uint64(exp+1023))<<52 |
		(uint64(frac) & 0xfffffffffffff)
	return math.Float64frombits(u)
}
