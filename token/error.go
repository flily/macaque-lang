package token

import (
	"fmt"
	"strings"

	"github.com/flily/macaque-lang/errors"
)

const (
	UnknownError            = iota
	ErrCodeSyntaxError      = 1
	ErrCodeRuntimeError     = 2
	ErrCodeCompilationError = 3
	ErrScannerError         = 100
)

type CodeContext struct {
	Filename  string
	NumLine   int
	NumColumn int
	Line      string
	Length    int
}

func (t *CodeContext) GetLeadingSpaces(index int) string {
	line := t.Line
	spaces := make([]byte, 0, index)

	for i := 0; i < index; i++ {
		switch line[i] {
		case '\t', '\v', '\f':
			spaces = append(spaces, line[i])
		default:
			spaces = append(spaces, ' ')
		}
	}

	return string(spaces)
}

func (c *CodeContext) MakeHighlight() string {
	leadingSpaces := c.GetLeadingSpaces(c.NumColumn - 1)
	highlight := strings.Repeat("^", c.Length)

	return leadingSpaces + highlight
}

func (c *CodeContext) MakeLineHighlight() string {
	highlight := c.MakeHighlight()
	return c.Line + "\n" + highlight
}

func (c *CodeContext) MakeMessage(format string, args ...interface{}) string {
	leadingSpaces := c.GetLeadingSpaces(c.NumColumn - 1)
	highlight := leadingSpaces + strings.Repeat("^", c.Length)
	lines := []string{
		c.Line,
		highlight,
		leadingSpaces + fmt.Sprintf(format, args...),
		fmt.Sprintf("  at %s:%d:%d", c.Filename, c.NumLine, c.NumColumn),
	}

	return strings.Join(lines, "\n")
}

func (c *CodeContext) NewSyntaxError(format string, args ...interface{}) *SyntaxError {
	return NewSyntaxError(c, format, args...)
}

func (c *CodeContext) NewCompilationError(format string, args ...interface{}) *CompilationError {
	return NewCompilationError(c, format, args...)
}

type PositionError struct {
	errors.BaseError

	Context *CodeContext
}

func NewPositionError(code int, format string, args ...interface{}) *PositionError {
	base := errors.NewRawError(code, format, args...)
	e := &PositionError{
		BaseError: *base,
	}

	return e
}

func (e *PositionError) WithContext(c *CodeContext) *PositionError {
	e.Context = c
	return e
}

func (e *PositionError) Derive(format string, args ...interface{}) *PositionError {
	base := e.BaseError.Derive(format, args...)
	err := &PositionError{
		BaseError: *base,
		Context:   e.Context,
	}

	return err
}
