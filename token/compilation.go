package token

import (
	"strings"

	"github.com/flily/macaque-lang/errors"
)

type CompilationNError struct {
	errors.BaseError

	Context *Context
	Info    []string
}

func MakeCompilationError(context *Context, format string, args ...interface{}) *CompilationNError {
	base := errors.NewRawError(ErrCodeCompilationError, format, args...)
	e := &CompilationNError{
		BaseError: *base,
		Context:   context,
	}

	return e
}

func (e *CompilationNError) WithInfo(ctx *Context, format string, args ...interface{}) *CompilationNError {
	s := ctx.Message(format, args...)
	e.Info = append(e.Info, s)
	return e
}

func (e *CompilationNError) Error() string {
	parts := make([]string, 1+len(e.Info))
	parts[0] = e.Context.Message(e.Message)
	for i, info := range e.Info {
		parts[i+1] = info
	}

	return strings.Join(parts, "\n")
}
