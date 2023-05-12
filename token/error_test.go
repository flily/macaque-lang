package token

import (
	"strings"
	"testing"

	"github.com/flily/macaque-lang/errors"
)

func TestErrorMessage(t *testing.T) {
	text := `the quick brown fox jumps over the lazy dog`
	tokenizer := NewSimpleTokenizer("testcase")
	tokenizer.SetContent([]byte(text))
	tokens := tokenizer.TokenList()

	ctx := &Context{
		tokens[4:6],
	}

	err := ctx.NewError(errors.UnknownError, "invalid verb '%s' for '%s'", "jumps over", "fox")
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
