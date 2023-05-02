package token

import (
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
