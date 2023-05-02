package token

import (
	"strings"
)

var (
	ErrCompilationError = NewPositionError(ErrCodeCompilationError, "compilation error")
)

type CompilationError struct {
	PositionError
	Info []string
}

func NewCompilationError(context *CodeContext, format string, args ...interface{}) *CompilationError {
	base := ErrCompilationError.Derive(format, args...).WithContext(context)
	e := &CompilationError{
		PositionError: *base,
	}

	return e
}

func (e *CompilationError) Error() string {
	parts := make([]string, 1+len(e.Info))
	parts[0] = e.Context.MakeMessage(e.Message)
	for i, info := range e.Info {
		parts[i+1] = info
	}

	return strings.Join(parts, "\n")
}

func (e *CompilationError) WithInfo(ctx *CodeContext, format string, args ...interface{}) *CompilationError {
	s := ctx.MakeMessage(format, args...)
	e.Info = append(e.Info, s)
	return e
}