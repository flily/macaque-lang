package token

import (
	"testing"

	"strings"
)

func TestInfo(t *testing.T) {
	f := NewFileInfo("sample.txt")

	if f.Filename != "sample.txt" {
		t.Errorf("filename is %s, expected %s", f.Filename, "sample.txt")
	}

	if len(f.Lines) != 0 {
		t.Errorf("lines is %d, expected %d", len(f.Lines), 0)
	}

	l1 := f.NewLine("the quick brown fox")
	_ = f.NewLine("jumps over")
	_ = f.NewLine("the lazy dog")

	if len(f.Lines) != 3 {
		t.Errorf("len(f.Lines) is %d, expected %d", len(f.Lines), 3)
	}

	if l1.File != f {
		t.Errorf("l1.File is %v, expected %v", l1.File, f)
	}

	t1 := l1.NewToken(1, 3, "the")
	if t1.Column != 1 {
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

	if t1.String() != "Token{the, sample.txt:1:1}" {
		t.Errorf("t1.String() is %s, expected %s", t1.String(), "Token{the, sample.txt:1:1}")
	}

	ctx1 := t1.MakeContext()
	if ctx1.Filename != "sample.txt" {
		t.Errorf("ctx1.Filename is %s, expected %s", ctx1.Filename, "sample.txt")
	}

	if ctx1.Line != "the quick brown fox" {
		t.Errorf("ctx1.Line is %s, expected %s", ctx1.Line, "the quick brown fox")
	}

	if ctx1.NumLine != 1 {
		t.Errorf("ctx1.NumLine is %d, expected %d", ctx1.NumLine, 1)
	}

	if ctx1.NumColumn != 1 {
		t.Errorf("ctx1.NumColumn is %d, expected %d", ctx1.NumColumn, 1)
	}

	if ctx1.Length != 3 {
		t.Errorf("ctx1.Length is %d, expected %d", ctx1.Length, 3)
	}
}

func TestTokenInfoMessage(t *testing.T) {
	f := NewFileInfo("sample.txt")
	l1 := f.NewLine("the quick brown fox")
	t1 := l1.NewToken(5, 5, "quick")

	expected := strings.Join([]string{
		"the quick brown fox",
		"    ^^^^^",
		"    lorem ipsum",
		"  at sample.txt:1:5",
	}, "\n")

	got := t1.MakeMessage("lorem ipsum")
	if got != expected {
		t.Errorf("got\n%v, expected\n%v", got, expected)
	}
}

func TestTokenInfoMessageWithSpaces(t *testing.T) {
	f := NewFileInfo("sample.txt")
	l1 := f.NewLine("   the 	 quick brown fox")
	t1 := l1.NewToken(10, 5, "quick")

	expected := strings.Join([]string{
		"   the 	 quick brown fox",
		"       	 ^^^^^",
		"       	 lorem ipsum",
		"  at sample.txt:1:10",
	}, "\n")

	got := t1.MakeMessage("lorem ipsum")
	if got != expected {
		t.Errorf("got\n%v, expected\n%v", got, expected)
	}
}

func TestTokenInfoHighlight(t *testing.T) {
	f := NewFileInfo("sample.txt")
	l1 := f.NewLine("the quick brown fox")
	t1 := l1.NewToken(5, 5, "quick")

	expected := strings.Join([]string{
		"the quick brown fox",
		"    ^^^^^",
	}, "\n")

	got := t1.MakeLineHighlight()
	if got != expected {
		t.Errorf("got %v, expected %v", got, expected)
	}
}
