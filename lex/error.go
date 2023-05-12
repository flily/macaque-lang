package lex

import (
	"github.com/flily/macaque-lang/errors"
	"github.com/flily/macaque-lang/token"
)

type LexicalError struct {
	token.ErrorWithContext
}

func NewLexicalError(ctx *token.TokenContext, format string, args ...interface{}) *LexicalError {
	context := ctx.ToContext()
	base := token.NewErrorWithContext(context, errors.ErrCodeLexicalError, format, args...)
	e := &LexicalError{
		ErrorWithContext: *base,
	}

	return e
}
