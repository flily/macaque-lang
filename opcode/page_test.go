package opcode

import (
	"strings"
	"testing"
)

func TestCodeBlockBasic(t *testing.T) {
	b := NewCodeBlock()
	b.IL(nil, ILoad, 6)
	b.IL(nil, ILoadInt, 42)

	if b.Length() != 2 {
		t.Errorf("Expected length 2, got %d", b.Length())
	}

	got := b.String()
	exp := strings.Join([]string{
		"CodeBlock[2]{",
		" 0    LOAD 6",
		" 1    LOADINT 42",
		"}",
	}, "\n")

	if got != exp {
		t.Errorf("wrong output string, expect %s, got %s", exp, got)
	}
}

func TestCodeBlockAppendAndPrepend(t *testing.T) {
	b := NewCodeBlock()
	b.IL(nil, ILoad, 6)
	b.IL(nil, ILoadInt, 42)

	b2 := NewCodeBlock()
	b2.IL(nil, ILoad, 7)
	b2.IL(nil, ILoadInt, 43)

	_ = b.Append(b2, nil)
	if b.Length() != 4 {
		t.Errorf("Expected length 4, got %d", b.Length())
	}

	got := b.String()
	exp := strings.Join([]string{
		"CodeBlock[4]{",
		" 0    LOAD 6",
		" 1    LOADINT 42",
		" 2    LOAD 7",
		" 3    LOADINT 43",
		"}",
	}, "\n")

	if got != exp {
		t.Errorf("wrong output string, expect %s, got %s", exp, got)
	}

	b3 := NewCodeBlock()
	b3.IL(nil, ILoadInt, 44)
	b3.PrependIL(nil, ILoad, 8)

	b.Prepend(b3)
	if b.Length() != 6 {
		t.Errorf("Expected length 6, got %d", b.Length())
	}

	got = b.String()
	exp = strings.Join([]string{
		"CodeBlock[6]{",
		" 0    LOAD 8",
		" 1    LOADINT 44",
		" 2    LOAD 6",
		" 3    LOADINT 42",
		" 4    LOAD 7",
		" 5    LOADINT 43",
		"}",
	}, "\n")

	if got != exp {
		t.Errorf("wrong output string, expect %s, got %s", exp, got)
	}
}
