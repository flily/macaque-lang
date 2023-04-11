package parser

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestConvertDecimalInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"42", 42},
		{"299792458", 299792458},
		{"299_792_458", 299792458},
	}

	for _, test := range tests {
		got := ConvertDecimalInteger(test.input)
		if got != test.expected {
			t.Errorf("ConvertDecimalInteger(%s) = %d, expected=%d",
				test.input, got, test.expected)
		}
	}
}

func BenchmarkConvertDecimalInteger(b *testing.B) {
	makeNums := func(n int) []string {
		nums := make([]string, n)
		for i := 0; i < n; i++ {
			nums[i] = fmt.Sprintf("%d", rand.Int63())
		}

		return nums
	}

	b.Run("benchmark ConvertDecimalInteger", func(bb *testing.B) {
		nums := makeNums(bb.N)

		for i := 0; i < bb.N; i++ {
			ConvertDecimalInteger(nums[i])
		}
	})

	b.Run("benchmark strconv.Parse", func(bb *testing.B) {
		nums := makeNums(bb.N)

		for i := 0; i < bb.N; i++ {
			_, _ = strconv.ParseInt(nums[i], 10, 64)
		}
	})
}

func TestConvertHexdecimalInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"0x42", 66},
		{"0xdeadbeef", 3735928559},
		{"0xdead_beef", 3735928559},
	}

	for _, test := range tests {
		got := ConvertHexdecimalInteger(test.input)
		if got != test.expected {
			t.Errorf("ConvertHexdecimalInteger(%s) = %d, expected=%d",
				test.input, got, test.expected)
		}
	}
}

func BenchmarkConvertHexdecimalInteger(b *testing.B) {
	makeNums := func(n int) []string {
		nums := make([]string, n)
		for i := 0; i < n; i++ {
			nums[i] = fmt.Sprintf("0x%x", rand.Int63())
		}

		return nums
	}

	b.Run("benchmark ConvertHexdecimalInteger", func(bb *testing.B) {
		nums := makeNums(bb.N)

		for i := 0; i < bb.N; i++ {
			ConvertHexdecimalInteger(nums[i])
		}
	})

	b.Run("benchmark strconv.Parse", func(bb *testing.B) {
		nums := makeNums(bb.N)

		for i := 0; i < bb.N; i++ {
			_, _ = strconv.ParseInt(nums[i], 16, 64)
		}
	})
}

func TestConvertFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"299792458", 299792458},
		{"3.1415926", 3.1415926},
		{"9.80_665", 9.80665},
	}

	for _, test := range tests {
		got := ConvertFloat(test.input)
		if got != test.expected {
			t.Errorf("ConvertFloat(%s) = %f, expected=%f",
				test.input, got, test.expected)
		}
	}
}
