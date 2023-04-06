package lex

import (
	"testing"

	"github.com/flily/macaque-lang/token"
)

func TestRelexerSetContent(t *testing.T) {
	code := `42
	3.1415926
	0xdeadbeef`

	lex := NewRecursiveScanner("testcase")
	err := lex.SetContent([]byte(code))
	if err != nil {
		t.Fatalf("SetContent() failed: %v", err)
	}

	if lex.source == nil {
		t.Fatal("lex.source is nil")
	}

	if len(lex.FileInfo.Lines) != 3 {
		t.Fatalf("len(lex.FileInfo.Lines) is %d, expected %d", len(lex.FileInfo.Lines), 3)
	}
}

func TestScanNumber(t *testing.T) {
	code := `42
		3.1415926
		0xdeadbeef  `

	lex := NewRecursiveScanner("testcase")
	_ = lex.SetContent([]byte(code))

	expected := []struct {
		token   token.Token
		content string
		line    int
		column  int
	}{
		{token.Integer, "42", 1, 1},
		{token.Float, "3.1415926", 2, 3},
		{token.Integer, "0xdeadbeef", 3, 3},
	}

	for _, e := range expected {
		token, err := lex.Scan()
		if err != nil {
			t.Fatalf("Scan() failed: %v", err)
		}

		if token.Token != e.token || token.Content != e.content {
			t.Errorf("unexpected token: %s", token)
		}

		if token.Position.Line.Line != e.line || token.Position.Column != e.column {
			t.Errorf("unexpected position: %+v", token.Position)
		}
	}
}
