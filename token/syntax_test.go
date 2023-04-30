package token

import (
	"testing"

	"strings"
)

func TestCodeContextLineHighlight(t *testing.T) {
	ctx := &CodeContext{
		Line:      "the quick brown fox jumps over the lazy dog",
		NumColumn: 21,
		Length:    5,
	}

	expected := []string{
		`the quick brown fox jumps over the lazy dog`,
		`                    ^^^^^`,
	}

	if ctx.MakeLineHighlight() != strings.Join(expected, "\n") {
		t.Errorf("expected line highlight got:\n%s", ctx.MakeLineHighlight())
	}

}

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

func TestSyntaxErrorMessageWithSpace(t *testing.T) {
	ctx := &CodeContext{
		Filename:  "test",
		Line:      "   the quick	brown fox jumps over the lazy dog",
		NumLine:   1,
		NumColumn: 24,
		Length:    5,
	}

	err := ctx.NewSyntaxError("invalid verb '%s'", "jumps")
	expected := strings.Join([]string{
		`   the quick	brown fox jumps over the lazy dog`,
		`            	          ^^^^^`,
		`            	          invalid verb 'jumps'`,
		`  at test:1:24`,
	}, "\n")

	if err.Error() != expected {
		t.Errorf("expected error message got:\n%s\nexpected:\n%s",
			err.Error(), expected)
	}
}
