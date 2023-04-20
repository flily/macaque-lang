package lex

import (
	"testing"

	"strings"

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

func TestReadEOF(t *testing.T) {
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

	elem, err := lex.Scan()
	if elem == nil || elem.Token != token.EOF {
		t.Errorf("Scan() after EOF should return non-nil EOF element, got: %+v",
			elem)
	}

	if err != nil {
		t.Errorf("Scan() after EOF should return nil, got: %v", err)
	}
}

func TestReadTokenAtTheEnd(t *testing.T) {
	code := `42
		3.1415926
		0xdeadbeef`

	lex := NewRecursiveScanner("testcase")
	_ = lex.SetContent([]byte(code))

	expected := []expectedTokenInfo{
		{token.Integer, "42", 1, 1},
		{token.Float, "3.1415926", 2, 3},
		{token.Integer, "0xdeadbeef", 3, 3},
		{token.EOF, "", 3, 13},
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
		{token.EOF, "", 3, 2},
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
		{token.EOF, "", 3, 11},
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
	if elem != nil {
		t.Fatalf("Scan() should return nil, got: %v", elem)
	}

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
}

func TestScanStringErrorOnInvalidEscape(t *testing.T) {
	code := `"foo\z"`

	lex := NewRecursiveScanner("testcase")
	_ = lex.SetContent([]byte(code))

	elem, err := lex.Scan()
	if elem != nil {
		t.Fatalf("Scan() should return nil, got: %v", elem)
	}

	if err == nil {
		t.Fatal("Scan() should fail")
	}

	expected := strings.Join([]string{
		`"foo\z"`,
		`     ^`,
		`     invalid escape sequence: \z`,
		`  at testcase:1:6`,
	}, "\n")

	if err.Error() != expected {
		t.Errorf("got wrong error message, got:\n%v", err)
	}
}

func TestStringEscapeAtTheEnd1(t *testing.T) {
	code := `"the quick brown fox\`

	lex := NewRecursiveScanner("testcase")
	_ = lex.SetContent([]byte(code))

	elem, err := lex.Scan()
	if elem != nil {
		t.Fatalf("Scan() should return nil, got: %v", elem)
	}

	if err == nil {
		t.Fatal("Scan() should fail")
	}

	expected := strings.Join([]string{
		`"the quick brown fox\`,
		`                     ^`,
		`                     unexpected EOF`,
		`  at testcase:1:22`,
	}, "\n")

	if err.Error() != expected {
		t.Errorf("got wrong error message, got:\n%v", err)
	}
}

func TestStringEscapeAtTheEnd2(t *testing.T) {
	code := `"the quick brown fox\x2`

	lex := NewRecursiveScanner("testcase")
	_ = lex.SetContent([]byte(code))

	elem, err := lex.Scan()
	if elem != nil {
		t.Fatalf("Scan() should return nil, got: %v", elem)
	}

	if err == nil {
		t.Fatal("Scan() should fail")
	}

	expected := strings.Join([]string{
		`"the quick brown fox\x2`,
		`                      ^^`,
		`                      insufficient characters for escape sequence`,
		`  at testcase:1:23`,
	}, "\n")

	if err.Error() != expected {
		t.Errorf("got wrong error message, got:\n%v", err)
	}
}

func TestScanPunctuations(t *testing.T) {
	code := `(){}[];,.
	=== !=== <= >=
	/-*+`

	lex := NewRecursiveScanner("testcase")
	_ = lex.SetContent([]byte(code))

	expected := []expectedTokenInfo{
		{token.LParen, "(", 1, 1},
		{token.RParen, ")", 1, 2},
		{token.LBrace, "{", 1, 3},
		{token.RBrace, "}", 1, 4},
		{token.LBracket, "[", 1, 5},
		{token.RBracket, "]", 1, 6},
		{token.Semicolon, ";", 1, 7},
		{token.Comma, ",", 1, 8},
		{token.Period, ".", 1, 9},
		{token.EQ, "==", 2, 2},
		{token.Assign, "=", 2, 4},
		{token.NE, "!=", 2, 6},
		{token.EQ, "==", 2, 8},
		{token.LE, "<=", 2, 11},
		{token.GE, ">=", 2, 14},
		{token.Slash, "/", 3, 2},
		{token.Minus, "-", 3, 3},
		{token.Asterisk, "*", 3, 4},
		{token.Plus, "+", 3, 5},
		{token.EOF, "", 3, 6},
	}

	checkTokenScan(t, lex, expected)
}

func TestScanPunctuationError(t *testing.T) {
	code := ` =#`

	lex := NewRecursiveScanner("testcase")
	_ = lex.SetContent([]byte(code))

	elem1, err := lex.Scan()
	if elem1 == nil || err != nil {
		t.Fatalf("Scan() failed: elem=%v err=%v", elem1, err)
	}

	elem2, err := lex.Scan()
	if elem2 != nil {
		t.Fatalf("unexpected elem2: %v", elem2)
	}

	if err == nil {
		t.Fatal("Scan() should fail")
	}

	expected := strings.Join([]string{
		` =#`,
		`  ^`,
		`  unknown operator '#'`,
		`  at testcase:1:3`,
	}, "\n")

	if err.Error() != expected {
		t.Errorf("got wrong error message, got:\n%v", err)
	}
}

func TestScanComment(t *testing.T) {
	code := `foorbar // a comment
	+`

	lex := NewRecursiveScanner("testcase")
	_ = lex.SetContent([]byte(code))

	expected := []expectedTokenInfo{
		{token.Identifier, "foorbar", 1, 1},
		{token.Comment, "// a comment", 1, 9},
		{token.Plus, "+", 2, 2},
	}

	checkTokenScan(t, lex, expected)
}

func TestEOFHighlight(t *testing.T) {
	code := strings.Join([]string{
		`foo`,
		`  bar`,
	}, "\n")

	lex := NewRecursiveScanner("testcase")
	_ = lex.SetContent([]byte(code))

	elem, err := lex.Scan()
	if elem == nil || err != nil {
		t.Fatalf("Scan() failed: elem=%v err=%v", elem, err)
	}

	elem, err = lex.Scan()
	if elem == nil || err != nil {
		t.Fatalf("Scan() failed: elem=%v err=%v", elem, err)
	}

	elem, err = lex.Scan()
	if elem == nil || err != nil {
		t.Fatalf("Scan() should return a EOF element, but gotelem=%v err=%v", elem, err)
	}

	highlight := elem.Position.MakeLineHighlight()
	expected := strings.Join([]string{
		`  bar`,
		`     ^`,
	}, "\n")
	if highlight != expected {
		t.Errorf("got wrong highlight, got:\n%s\nexpected:\n%s",
			highlight, expected)
	}
}
