package token

import (
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

type SyntaxError struct {
	errors.BaseError

	Context *Context
}

func MakeSyntaxError(context *Context, format string, args ...interface{}) *SyntaxError {
	base := errors.NewRawError(ErrCodeSyntaxError, format, args...)
	e := &SyntaxError{
		BaseError: *base,
		Context:   context,
	}

	return e
}

func (e *SyntaxError) Error() string {
	return e.Context.Message(e.Message)
}

type CompilationError struct {
	errors.BaseError

	Context *Context
	Info    []string
}

func MakeCompilationError(context *Context, format string, args ...interface{}) *CompilationError {
	base := errors.NewRawError(ErrCodeCompilationError, format, args...)
	e := &CompilationError{
		BaseError: *base,
		Context:   context,
	}

	return e
}

func (e *CompilationError) WithInfo(ctx *Context, format string, args ...interface{}) *CompilationError {
	s := ctx.Message(format, args...)
	e.Info = append(e.Info, s)
	return e
}

func (e *CompilationError) Error() string {
	parts := make([]string, 1+len(e.Info))
	parts[0] = e.Context.Message(e.Message)
	for i, info := range e.Info {
		parts[i+1] = info
	}

	return strings.Join(parts, "\n")
}
