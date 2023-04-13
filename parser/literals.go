package parser

import (
	"github.com/flily/macaque-lang/ast"
	"github.com/flily/macaque-lang/lex"
	"github.com/flily/macaque-lang/token"
)

const (
	BooleanTrue  = token.True
	BooleanFalse = token.False
)

func makeInteger(content string, position *token.TokenInfo) *ast.IntegerLiteral {
	literal := &ast.IntegerLiteral{
		Value:    ConvertInteger(content),
		Content:  content,
		Position: position,
	}

	return literal
}

func NewInteger(token *lex.LexicalElement) *ast.IntegerLiteral {
	return makeInteger(token.Content, token.Position)
}

func makeFloat(content string, position *token.TokenInfo) *ast.FloatLiteral {
	literal := &ast.FloatLiteral{
		Value:    ConvertFloat(content),
		Content:  content,
		Position: position,
	}

	return literal
}

func NewFloat(token *lex.LexicalElement) *ast.FloatLiteral {
	return makeFloat(token.Content, token.Position)
}

func makeString(content string, position *token.TokenInfo) *ast.StringLiteral {
	literal := &ast.StringLiteral{
		Value:    ConvertString(content),
		Content:  content,
		Position: position,
	}

	return literal
}

func NewString(token *lex.LexicalElement) *ast.StringLiteral {
	return makeString(token.Content, token.Position)
}

func makeNull(position *token.TokenInfo) *ast.NullLiteral {
	literal := &ast.NullLiteral{
		Position: position,
	}

	return literal
}

func NewNull(token *lex.LexicalElement) *ast.NullLiteral {
	return makeNull(token.Position)
}

func makeBoolean(value bool, position *token.TokenInfo) *ast.BooleanLiteral {
	literal := &ast.BooleanLiteral{
		Value:    value,
		Position: position,
	}

	return literal
}

func NewBoolean(token *lex.LexicalElement) *ast.BooleanLiteral {
	return makeBoolean(token.Token == BooleanTrue, token.Position)
}
