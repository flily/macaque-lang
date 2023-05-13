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

type UnexpectedTokenError struct {
	SyntaxError

	Expecteds []token.Token
	Actual    *token.TokenContext
}

func NewUnexpectedTokenError(ctx *token.Context, actual *token.TokenContext, expecteds ...token.Token) *UnexpectedTokenError {
	base := token.NewErrorWithContext(ctx, errors.ErrCodeUnexpectedToken,
		"unexpected token %s, expected %s", actual.Token, expecteds)

	e := &UnexpectedTokenError{
		SyntaxError: SyntaxError{
			ErrorWithContext: *base,
		},

		Expecteds: expecteds,
		Actual:    actual,
	}

	return e
}

func (e *UnexpectedTokenError) WithMessage(format string, args ...interface{}) *UnexpectedTokenError {
	_ = e.SyntaxError.WithMessage(format, args...)
	return e
}
