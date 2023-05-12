package parser

import (
	"github.com/flily/macaque-lang/errors"
	"github.com/flily/macaque-lang/token"
)

type SyntaxError struct {
	token.ErrorWithContext
}

func NewSyntaxError(ctx *token.Context, format string, args ...interface{}) *SyntaxError {
	base := token.NewErrorWithContext(ctx, errors.ErrCodeSyntaxError, format, args...)
	e := &SyntaxError{
		ErrorWithContext: *base,
	}

	return e
}
