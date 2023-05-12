package lex

import (
	"github.com/flily/macaque-lang/errors"
	"github.com/flily/macaque-lang/token"
)

var (
	ErrScannerHasContentAlready = errors.NewError(errors.ErrScannerError, "scanner has content already")
)

type Scanner interface {
	Scan() (*token.TokenContext, error)
}
