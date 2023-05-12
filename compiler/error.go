package compiler

import (
	"strings"

	"github.com/flily/macaque-lang/errors"
	"github.com/flily/macaque-lang/token"
)

type CompilationError struct {
	token.ErrorWithContext

	Info []string
}

func NewCompilationError(context *token.Context, format string, args ...interface{}) *CompilationError {
	base := token.NewErrorWithContext(context, errors.ErrCodeCompilationError, format, args...)
	e := &CompilationError{
		ErrorWithContext: *base,
	}

	return e
}

func (e *CompilationError) WithInfo(ctx *token.Context, format string, args ...interface{}) *CompilationError {
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
