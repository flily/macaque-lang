package token

import "github.com/flily/macaque-lang/errors"

var (
	ErrSyntaxError = NewPositionError(ErrCodeSyntaxError, "syntax error")
)

type SyntaxError struct {
	PositionError
}

type SyntaxNError struct {
	errors.BaseError

	Context *TokenContext
}

func NewSyntaxError(context *CodeContext, format string, args ...interface{}) *SyntaxError {
	base := ErrSyntaxError.Derive(format, args...).WithContext(context)
	e := &SyntaxError{
		PositionError: *base,
	}

	return e
}

func (e *SyntaxError) Error() string {
	return e.Context.MakeMessage(e.Message)
}

func MakeSyntaxError(context *TokenContext, format string, args ...interface{}) *SyntaxNError {
	base := errors.NewRawError(ErrCodeSyntaxError, format, args...)
	e := &SyntaxNError{
		BaseError: *base,
		Context:   context,
	}

	return e
}

func (e *SyntaxNError) Error() string {
	return e.Context.Position.MakeMessage(e.Message)
}
