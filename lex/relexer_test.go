package lex

import (
	"strings"
	"testing"

	"github.com/flily/macaque-lang/token"
)

type expectedTokenInfo struct {
	token   token.Token
	content string
	line    int
	column  int
}

func checkTokenScan(t *testing.T, s Scanner, expected []expectedTokenInfo) {
	for _, c := range expected {
		elem, err := s.Scan()
		if err != nil {
			t.Fatalf("Scan(%s) failed: %v", c.content, err)
		}

		if elem.Content != c.content {
			t.Errorf("Scan(%s) got wrong content: %s", c.content, elem.Content)
		}

		if elem.Token != c.token {
			t.Errorf("Scan(%s) got wrong token: %s, expected: %s",
				c.content, elem.Token, c.token)
		}

		line, column := elem.Position.GetPosition()

		if line != c.line || column != c.column {
			t.Errorf("Scan(%s) got wrong position: %d:%d, expected %d:%d",
				c.content, line, column, c.line, c.column)
		}
	}
}

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

	err = lex.SetContent([]byte(code))
	if err == nil {
		t.Fatal("SetContent() twice should fail")
	}

	if err != ErrScannerHasContentAlready {
		t.Errorf("SetContent() twice should return ErrScannerHasContentAlready, got: %v", err)
	}
}

func TestScanNumber(t *testing.T) {
	code := `42
		3.1415926
		0xdeadbeef  `

	lex := NewRecursiveScanner("testcase")
	_ = lex.SetContent([]byte(code))

	expected := []expectedTokenInfo{
		{token.Integer, "42", 1, 1},
		{token.Float, "3.1415926", 2, 3},
		{token.Integer, "0xdeadbeef", 3, 3},
	}

	checkTokenScan(t, lex, expected)
}

func TestScanIdentifierAndKeyworkd(t *testing.T) {
	code := `foobar
		if foo else bar return
	`

	lex := NewRecursiveScanner("testcase")
	_ = lex.SetContent([]byte(code))

	expected := []expectedTokenInfo{
		{token.Identifier, "foobar", 1, 1},
		{token.If, "if", 2, 3},
		{token.Identifier, "foo", 2, 6},
		{token.Else, "else", 2, 10},
		{token.Identifier, "bar", 2, 15},
		{token.Return, "return", 2, 19},
	}

	checkTokenScan(t, lex, expected)
}

func TestScannerAppend(t *testing.T) {
	lex := NewRecursiveScanner("testcase")
	lex.Append([]byte("foobar"))
	lex.Append([]byte("  42"))

	expected := []expectedTokenInfo{
		{token.Identifier, "foobar", 1, 1},
		{token.Integer, "42", 2, 3},
	}

	checkTokenScan(t, lex, expected)
}

func TestScanStrings(t *testing.T) {
	code := ` "foobar"
		"foo\nbar\""
		"c+\x2b"`

	expected := []expectedTokenInfo{
		{token.String, "\"foobar\"", 1, 2},
		{token.String, "\"foo\\nbar\\\"\"", 2, 3},
		{token.String, "\"c+\\x2b\"", 3, 3},
	}

	lex := NewRecursiveScanner("testcase")
	_ = lex.SetContent([]byte(code))

	checkTokenScan(t, lex, expected)
}

func TestScanStringErrorOnInvalidHexdecimalEscape(t *testing.T) {
	code := `"foo\x2"`

	lex := NewRecursiveScanner("testcase")
	_ = lex.SetContent([]byte(code))

	elem, err := lex.Scan()
	if err == nil {
		t.Fatal("Scan() should fail")
	}

	expected := strings.Join([]string{
		`"foo\x2"`,
		`      ^^`,
		`      invalid escape sequence: \x2"`,
		`  at testcase:1:7`,
	}, "\n")

	if err.Error() != expected {
		t.Errorf("got wrong error message, got:\n%v", err)
	}

	if elem != nil {
		t.Errorf("Scan() should return nil, got: %v", elem)
	}
}
