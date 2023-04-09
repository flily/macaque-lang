package errors

import (
	"testing"

	"strings"
)

func TestSyntaxErrorMessage(t *testing.T) {
	ctx := &CodeContext{
		Filename:  "test",
		Line:      "the quick brown fox jumps over the lazy dog",
		NumLine:   1,
		NumColumn: 21,
		Length:    5,
	}

	err := ctx.NewSyntaxError("invalid verb '%s'", "jumps")
	expected := []string{
		`the quick brown fox jumps over the lazy dog`,
		`                    ^^^^^`,
		`                    invalid verb 'jumps'`,
		`  at test:1:21`,
	}

	if err.Error() != strings.Join(expected, "\n") {
		t.Errorf("expected error message got:\n%s", err.Error())
	}
}
