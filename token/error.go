package token

import (
	"github.com/flily/macaque-lang/errors"
)

type ErrorWithContext struct {
	errors.BaseError

	Context *Context
}

func NewErrorWithContext(context *Context, code int, format string, args ...interface{}) *ErrorWithContext {
	base := errors.NewRawError(code, format, args...)
	e := &ErrorWithContext{
		BaseError: *base,
		Context:   context,
	}

	return e
}

func (e *ErrorWithContext) Error() string {
	return e.Context.Message(e.Message)
}
