package token

import (
	"strings"
	"testing"
)

func TestContextHighligh(t *testing.T) {
	code := strings.Join([]string{
		"   the quick	brown fox",
		"		jumps   over	 the lazy dog",
	}, "\n")

	tokenizer := NewSimpleTokenizer("testcase")
	tokenizer.SetContent([]byte(code))
	tokens := tokenizer.TokenList()
	if len(tokens) != 9 {
		t.Fatalf("got %d tokens, want 9", len(tokens))
	}

	ctx := &Context{
		Tokens: tokens[3:7],
	}

	expected := strings.Join([]string{
		`   the quick	brown fox`,
		`            	      ^^^`,
		`		jumps   over	 the lazy dog`,
		`		^^^^^   ^^^^	 ^^^`,
	}, "\n")

	got := ctx.HighLight()
	if got != expected {
		t.Errorf("expected error highlight got:\n%s", got)
	}
}

func TestContextMessage(t *testing.T) {
	code := strings.Join([]string{
		"   the quick	brown fox",
		"		jumps   over	 the lazy dog",
	}, "\n")

	tokenizer := NewSimpleTokenizer("testcase")
	tokenizer.SetContent([]byte(code))
	tokens := tokenizer.TokenList()
	if len(tokens) != 9 {
		t.Fatalf("got %d tokens, want 9", len(tokens))
	}

	ctx := &Context{
		Tokens: tokens[3:6],
	}

	expected := strings.Join([]string{
		`   the quick	brown fox`,
		`            	      ^^^`,
		`		jumps   over	 the lazy dog`,
		`		^^^^^   ^^^^`,
		`		invalid verb 'jumps' for 'fox'`,
		`  at testcase:1:20`,
	}, "\n")

	got := ctx.Message("invalid verb 'jumps' for 'fox'")
	if got != expected {
		t.Errorf("expected error message got:\n%s", got)
	}
}
