package token

import "github.com/flily/macaque-lang/errors"

type SyntaxNError struct {
	errors.BaseError

	Context *Context
}

func MakeSyntaxError(context *Context, format string, args ...interface{}) *SyntaxNError {
	base := errors.NewRawError(ErrCodeSyntaxError, format, args...)
	e := &SyntaxNError{
		BaseError: *base,
		Context:   context,
	}

	return e
}

func (e *SyntaxNError) Error() string {
	return e.Context.Message(e.Message)
}
