package errors

import (
	"errors"
	"fmt"
	"strings"
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

type BaseError struct {
	Kind    int
	Message string
	Base    *BaseError
	Context *CodeContext
}

type HasKind interface {
	ErrorKind() int
}

type CanDerive interface {
	Derive(format string, args ...interface{}) *BaseError
}

func NewRawError(kind int, format string, args ...interface{}) *BaseError {
	e := &BaseError{
		Kind:    kind,
		Message: fmt.Sprintf(format, args...),
	}

	return e
}

func NewError(kind int, format string, args ...interface{}) error {
	return NewRawError(kind, format, args...)
}

func (e *BaseError) Error() string {
	return e.Message
}

func KindOf(err error) int {
	if err, ok := err.(HasKind); ok {
		return err.ErrorKind()
	}

	return UnknownError
}

func (e *BaseError) ErrorKind() int {
	return e.Kind
}

// Is wraps errors.Is
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func (e *BaseError) Is(target error) bool {
	if err, ok := target.(HasKind); ok {
		return err.ErrorKind() == e.ErrorKind()
	}

	return false
}

func Derive(base error, format string, args ...interface{}) error {
	if err, ok := base.(CanDerive); ok {
		return err.Derive(format, args...)
	}

	return NewError(UnknownError, format, args...)
}

func (e *BaseError) Derive(format string, args ...interface{}) *BaseError {
	err := NewRawError(e.Kind, format, args...)
	err.Base = e
	return err
}

func (e *BaseError) WithContext(c *CodeContext) *BaseError {
	e.Context = c
	return e
}
