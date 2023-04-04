package token

import (
	"testing"
)

func TestInfo(t *testing.T) {
	f := NewFileInfo("sample.txt")

	if f.Filename != "sample.txt" {
		t.Errorf("filename is %s, expected %s", f.Filename, "sample.txt")
	}

	if len(f.Lines) != 0 {
		t.Errorf("lines is %d, expected %d", len(f.Lines), 0)
	}

	l1 := f.NewLine("the quick brown fox", 1)
	_ = f.NewLine("jumps over", 2)
	_ = f.NewLine("the lazy dog", 3)

	if len(f.Lines) != 3 {
		t.Errorf("len(f.Lines) is %d, expected %d", len(f.Lines), 3)
	}

	if l1.File != f {
		t.Errorf("l1.File is %v, expected %v", l1.File, f)
	}

	t1 := l1.NewToken(0, 3, "the")
	if t1.Column != 0 {
		t.Errorf("t1.Column is %d, expected %d", t1.Column, 0)
	}

	if t1.Length != 3 {
		t.Errorf("t1.Length is %d, expected %d", t1.Length, 3)
	}

	if t1.Content != "the" {
		t.Errorf("t1.Content is %s, expected %s", t1.Content, "the")
	}

	if t1.Line != l1 {
		t.Errorf("t1.Line is %v, expected %v", t1.Line, l1)
	}

	if len(l1.Tokens) != 1 {
		t.Errorf("len(l1.Tokens) is %d, expected %d", len(l1.Tokens), 1)
	}
}
