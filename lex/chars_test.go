package lex

import "testing"

func inBytesSet(c byte, set []byte) bool {
	for _, v := range set {
		if c == v {
			return true
		}
	}

	return false
}

func TestIsSpace(t *testing.T) {
	expectedSpace := []byte{
		' ', '\t', '\n', '\r',
	}

	for i := 0; i < 256; i++ {
		c := byte(i)
		in := inBytesSet(c, expectedSpace)
		is := IsSpace(c)
		if in != is {
			t.Errorf("IsSpace(%d) = %v, want %v", i, is, in)
		}
	}
}

func TestIsDigit(t *testing.T) {
	expectedDigits := []byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}

	for i := 0; i < 256; i++ {
		c := byte(i)
		in := inBytesSet(c, expectedDigits)
		is := IsDigit(c)
		if in != is {
			t.Errorf("IsDigit(%c %d) = %v, want %v", i, i, is, in)
		}
	}
}

func TestIsHexDigit(t *testing.T) {
	expectedDigits := []byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f',
		'A', 'B', 'C', 'D', 'E', 'F',
	}

	for i := 0; i < 256; i++ {
		c := byte(i)
		in := inBytesSet(c, expectedDigits)
		is := IsHexDigit(c)
		if in != is {
			t.Errorf("IsHexDigit('%c' %d) = %v, want %v", i, i, is, in)
		}
	}
}

func TestIsLower(t *testing.T) {
	expectedLower := []byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}

	for i := 0; i < 256; i++ {
		c := byte(i)
		in := inBytesSet(c, expectedLower)
		is := IsLower(c)
		if in != is {
			t.Errorf("IsLower('%c' %d) = %v, want %v", i, i, is, in)
		}
	}
}

func TestIsUpper(t *testing.T) {
	expectedUpper := []byte{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	}

	for i := 0; i < 256; i++ {
		c := byte(i)
		in := inBytesSet(c, expectedUpper)
		is := IsUpper(c)
		if in != is {
			t.Errorf("IsUpper('%c' %d) = %v, want %v", i, i, is, in)
		}
	}
}

func TestIsPunct(t *testing.T) {
	expectedPunct := []byte{
		'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-',
		'.', '/', ':', ';', '<', '=', '>', '?', '@', '[', '\\', ']', '^',
		'_', '`', '{', '|', '}', '~',
	}

	for i := 0; i < 256; i++ {
		c := byte(i)
		in := inBytesSet(c, expectedPunct)
		is := IsPunct(c)
		if in != is {
			t.Errorf("IsPunct('%c' %d) = %v, want %v", i, i, is, in)
		}
	}
}
