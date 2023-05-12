package token

import (
	"strings"
	"testing"
)

func TestSyntaxErrorMessage(t *testing.T) {
	text := `the quick brown fox jumps over the lazy dog`
	tokenizer := NewSimpleTokenizer("testcase")
	tokenizer.SetContent([]byte(text))
	tokens := tokenizer.TokenList()

	ctx := &Context{
		tokens[4:6],
	}

	err := ctx.NewSyntaxError("invalid verb '%s' for '%s'", "jumps over", "fox")
	expected := strings.Join([]string{
		"the quick brown fox jumps over the lazy dog",
		"                    ^^^^^ ^^^^",
		"                    invalid verb 'jumps over' for 'fox'",
		"  at testcase:1:21",
	}, "\n")

	if err.Error() != expected {
		t.Errorf("got wrong error message, got:\n%v", err)
	}
}

func TestCompilationErrorMessage(t *testing.T) {
	text := `the quick brown fox jumps over the lazy dog`
	tokenizer := NewSimpleTokenizer("testcase")
	tokenizer.SetContent([]byte(text))
	tokens := tokenizer.TokenList()

	ctx1 := &Context{
		tokens[4:6],
	}

	ctx2 := &Context{
		tokens[2:4],
	}

	err := ctx1.NewCompilationError(
		"invalid verb '%s' for '%s'", "jumps over", "brown fox",
	).WithInfo(ctx2, "subject '%s' is here", "brown fox")

	expected := strings.Join([]string{
		"the quick brown fox jumps over the lazy dog",
		"                    ^^^^^ ^^^^",
		"                    invalid verb 'jumps over' for 'brown fox'",
		"  at testcase:1:21",
		"the quick brown fox jumps over the lazy dog",
		"          ^^^^^ ^^^",
		"          subject 'brown fox' is here",
		"  at testcase:1:11",
	}, "\n")

	if err.Error() != expected {
		t.Errorf("got wrong error message, got:\n%v", err)
	}
}
