package lex

import (
	"github.com/flily/macaque-lang/errors"
	"github.com/flily/macaque-lang/token"
)

var (
	ErrScannerHasContentAlready = errors.NewError(errors.LexicalError, "scanner has content already")
	InvalidLexicalElement       = &LexicalElement{
		Token: token.Illegal,
	}
)

type LexicalElement struct {
	Token    token.Token
	Content  string
	Position *token.TokenInfo
}

type Scanner interface {
	Scan() (*LexicalElement, error)
}