package errors

import (
	"fmt"
	"strings"
)

var ErrSyntaxError = NewRawError(ErrCodeSyntaxError, "syntax error")

type CodeContext struct {
	Filename  string
	NumLine   int
	NumColumn int
	Line      string
	Length    int
}

func (c *CodeContext) NewSyntaxError(format string, args ...interface{}) *SyntaxError {
	return NewSyntaxError(c, format, args...)
}

type SyntaxError struct {
	BaseError
	Context *CodeContext
}

func NewSyntaxError(context *CodeContext, format string, args ...interface{}) *SyntaxError {
	e := &SyntaxError{
		BaseError: BaseError{
			Kind:    ErrCodeSyntaxError,
			Message: fmt.Sprintf(format, args...),
			Base:    ErrSyntaxError,
		},
		Context: context,
	}

	return e
}

func (e *SyntaxError) Error() string {
	leadingSpaces := strings.Repeat(" ", e.Context.NumColumn-1)
	highlight := strings.Repeat("^", e.Context.Length)
	lines := []string{
		e.Context.Line,
		leadingSpaces + highlight,
		leadingSpaces + e.Message,
		fmt.Sprintf("  at %s:%d:%d", e.Context.Filename, e.Context.NumLine, e.Context.NumColumn),
	}

	return strings.Join(lines, "\n")
}
