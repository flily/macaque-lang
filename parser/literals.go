package parser

import (
	"github.com/flily/macaque-lang/ast"
	"github.com/flily/macaque-lang/token"
)

const (
	BooleanTrue  = token.True
	BooleanFalse = token.False
)

func newInteger(token *token.TokenContext) *ast.IntegerLiteral {
	content := token.Content

	literal := &ast.IntegerLiteral{
		Value:   ConvertInteger(content),
		Content: content,
		Context: token.ToContext(),
	}

	return literal
}

func newFloat(token *token.TokenContext) *ast.FloatLiteral {
	content := token.Content
	literal := &ast.FloatLiteral{
		Value:   ConvertFloat(content),
		Content: content,
		Context: token.ToContext(),
	}

	return literal
}

func newString(token *token.TokenContext) *ast.StringLiteral {
	content := token.Content
	literal := &ast.StringLiteral{
		Value:   ConvertString(content),
		Content: content,
		Context: token.ToContext(),
	}

	return literal
}

func newNull(token *token.TokenContext) *ast.NullLiteral {
	literal := &ast.NullLiteral{
		Context: token.ToContext(),
	}

	return literal
}

func newBoolean(token *token.TokenContext) *ast.BooleanLiteral {
	literal := &ast.BooleanLiteral{
		Value:   token.Token == BooleanTrue,
		Context: token.ToContext(),
	}

	return literal
}
