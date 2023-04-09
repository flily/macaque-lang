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

func (c *CodeContext) MakeHighlight() string {
	leadingSpaces := strings.Repeat(" ", c.NumColumn-1)
	highlight := strings.Repeat("^", c.Length)

	return leadingSpaces + highlight
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
	highlight := e.Context.MakeHighlight()
	lines := []string{
		e.Context.Line,
		highlight,
		leadingSpaces + e.Message,
		fmt.Sprintf("  at %s:%d:%d", e.Context.Filename, e.Context.NumLine, e.Context.NumColumn),
	}

	return strings.Join(lines, "\n")
}
