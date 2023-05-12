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

func TestTokenContextGetToken(t *testing.T) {
	text := `the quick brown fox jumps over the lazy dog`
	tokenizer := NewSimpleTokenizer("testcase")
	tokenizer.SetContent([]byte(text))
	tokens := tokenizer.TokenList()

	tokenFox := tokens[3]
	if tokenFox.GetToken() != String {
		t.Errorf("expected token type String, got %s", tokenFox.GetToken())
	}

	var tokenNil *TokenContext
	if tokenNil.GetToken() != Nil {
		t.Errorf("expected token type Nil, got %s", tokenNil.GetToken())
	}
}

func TestTokenContextToContext(t *testing.T) {
	text := `the quick brown fox jumps over the lazy dog`
	tokenizer := NewSimpleTokenizer("testcase")
	tokenizer.SetContent([]byte(text))
	tokens := tokenizer.TokenList()

	tokenFox := tokens[3]
	{
		ctx := tokenFox.ToContext()
		expected := strings.Join([]string{
			`the quick brown fox jumps over the lazy dog`,
			`                ^^^`,
		}, "\n")

		if ctx.HighLight() != expected {
			t.Errorf("expected error highlight got:\n%s", ctx.HighLight())
		}
	}

	{
		var tokenNil *TokenContext
		ctx := tokenNil.ToContext()
		expected := ""

		if ctx.HighLight() != expected {
			t.Errorf("expected error highlight got:'%s'", ctx.HighLight())
		}
	}
}

func TestJoinContext(t *testing.T) {
	text := `the quick brown fox jumps over the lazy dog`
	tokenizer := NewSimpleTokenizer("testcase")
	tokenizer.SetContent([]byte(text))
	tokens := tokenizer.TokenList()

	tokenJumps := tokens[4]
	tokenOver := tokens[5]

	ctx1 := JoinContext(tokenJumps.ToContext(), tokenOver.ToContext())
	expected := strings.Join([]string{
		`the quick brown fox jumps over the lazy dog`,
		`                    ^^^^^ ^^^^`,
	}, "\n")

	if ctx1.HighLight() != expected {
		t.Errorf("expected error highlight got:\n%s", ctx1.HighLight())
	}

	var tokenNil *TokenContext
	ctx2 := JoinContext(
		tokenJumps.ToContext(),
		tokenNil.ToContext(),
		tokenOver.ToContext(),
	)

	if ctx2.HighLight() != expected {
		t.Errorf("expected error highlight got:\n%s", ctx2.HighLight())
	}
}
