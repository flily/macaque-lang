package token

import (
	"testing"
)

func TestSimpleTokenizer1(t *testing.T) {
	text := "lorem ipsum"

	tokenizer := NewSimpleTokenizer("test")
	tokenizer.SetContent([]byte(text))

	got := tokenizer.TokenList()
	if len(got) != 2 {
		t.Errorf("got %d tokens, want 2", len(got))
	}
}
