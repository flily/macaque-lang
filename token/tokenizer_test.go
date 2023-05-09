package token

import (
	"strings"
	"testing"
)

type testTokenContextInfo struct {
	Token   Token
	Line    int
	Column  int
	Content string
}

func (i testTokenContextInfo) Check(c *TokenContext) bool {
	if i.Token != c.Token {
		return false
	}

	if i.Line != c.LineNo() {
		return false
	}

	if i.Column != c.ColumnStart() {
		return false
	}

	if i.Content != c.Content {
		return false
	}

	return true
}

func TestSimpleTokenizer1(t *testing.T) {
	text := "lorem ipsum"
	expected := []testTokenContextInfo{
		{
			Token:   String,
			Line:    1,
			Column:  1,
			Content: "lorem",
		},
		{
			Token:   String,
			Line:    1,
			Column:  7,
			Content: "ipsum",
		},
	}

	tokenizer := NewSimpleTokenizer("test")
	tokenizer.SetContent([]byte(text))

	got := tokenizer.TokenList()
	if len(got) != 2 {
		t.Errorf("got %d tokens, want 2", len(got))
	}

	for i, token := range got {
		if !expected[i].Check(token) {
			t.Errorf("got %v, want %v", token, expected[i])
		}
	}
}

func TestRejectToken(t *testing.T) {
	text := strings.Join([]string{
		"lorem",
		"  ipsum~",
	}, "\n")

	expected := []testTokenContextInfo{
		{
			Token:   String,
			Line:    1,
			Column:  1,
			Content: "lorem",
		},
		{
			Token:   Illegal,
			Line:    2,
			Column:  8,
			Content: "~",
		},
	}

	tokenizer := NewSimpleTokenizer("test")
	tokenizer.SetContent([]byte(text))

	got := tokenizer.TokenList()
	if len(got) != 2 {
		t.Errorf("got %d tokens, want 2", len(got))
	}

	for i, token := range got {
		if !expected[i].Check(token) {
			t.Errorf("got %v, want %v", token, expected[i])
		}
	}
}

func TestPeekCharAndString(t *testing.T) {
	text := "lorem ipsum"

	tokenizer := NewSimpleTokenizer("test")
	tokenizer.SetContent([]byte(text))

	if c := tokenizer.PeekChar(0); c != 'l' {
		t.Errorf("PeekChar(0) != 'l', got '%c'", c)
	}

	if c := tokenizer.PeekChar(1); c != 'o' {
		t.Errorf("PeekChar(1) != 'o', got '%c'", c)
	}

	if !tokenizer.Peek("lorem") {
		t.Errorf("Peek('lorem') failed")
	}

	if tokenizer.Peek("loREm") {
		t.Errorf("Peek('loREm') failed")
	}

	if tokenizer.Peek("the quick brown fox jumps over the lazy dog") {
		t.Errorf("Peek() long string failed")
	}
}

func TestCharsLeft(t *testing.T) {
	text := "lorem ipsum"

	tokenizer := NewSimpleTokenizer("test")
	tokenizer.SetContent([]byte(text))

	l := len(text)
	for i := 0; i < l; i++ {
		if n := tokenizer.CharsLeft(); n != l {
			t.Errorf("%d chars left after %d shifts, expect %d", n, i, l)
		}
	}
}
