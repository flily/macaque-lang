package lex

import (
	"github.com/flily/macaque-lang/errors"
	"github.com/flily/macaque-lang/token"
)

var (
	ErrScannerHasContentAlready = errors.NewError(token.ErrScannerError, "scanner has content already")
)

type LexicalElement struct {
	Token    token.Token
	Content  string
	Position *token.TokenInfo
}

type Scanner interface {
	Scan() (*LexicalElement, error)
}
