package token

import (
	"testing"
)

func TestTokenString(t *testing.T) {
	tests := []struct {
		token    Token
		expected string
	}{
		{Illegal, "ILLEGAL"},
		{EOF, "EOF"},
		{Identifier, "IDENTIFIER"},
		{Integer, "INTEGER"},
		{Let, "LET"},
		{Return, "RETURN"},
		{Plus, "PLUS(+)"},
		{LT, "LT(<)"},
		{AND, "AND(&&)"},
		{Assign, "ASSIGN(=)"},
		{Semicolon, "SEMICOLON(;)"},
		{LBracket, "LBRACKET('[')"},
		{RBracket, "RBRACKET(']')"},
		{LBrace, "LBRACE('{')"},
		{RBrace, "RBRACE('}')"},
		{LParen, "LPAREN('(')"},
		{RParen, "RPAREN(')')"},
		{-1, "ILLEGAL"},
	}

	for _, c := range tests {
		s := c.token.String()
		if s != c.expected {
			t.Errorf("wrong token.String() of %d, expected %s, got %s", c.token, c.expected, s)
		}
	}
}

func TestTokenCodeName(t *testing.T) {
	tests := []struct {
		token    Token
		expected string
	}{
		{Illegal, "ILLEGAL"},
		{EOF, "EOF"},
		{Identifier, "IDENTIFIER"},
		{Integer, "INTEGER"},
		{Let, "LET"},
		{Return, "RETURN"},
		{Plus, "+"},
		{LT, "<"},
		{AND, "&&"},
		{Assign, "="},
		{Semicolon, ";"},
		{-1, "ILLEGAL"},
	}

	for _, c := range tests {
		s := c.token.CodeName()
		if s != c.expected {
			t.Errorf("wrong token.String() of %d, expected %s, got %s", c.token, c.expected, s)
		}
	}
}

func TestTokenType(t *testing.T) {
	tests := []struct {
		token      Token
		isLiteral  bool
		isKeyword  bool
		isOperator bool
	}{
		{Illegal, false, false, false},
		{EOF, false, false, false},
		{Identifier, true, false, false},
		{Integer, true, false, false},
		{Let, false, true, false},
		{Return, false, true, false},
		{Plus, false, false, true},
		{AND, false, false, true},
		{Assign, false, false, true},
		{Semicolon, false, false, true},
	}

	for _, c := range tests {
		if got := c.token.IsLiteral(); got != c.isLiteral {
			t.Errorf("%s.IsLiteral() is %t, expected %t", c.token, got, c.isLiteral)
		}

		if got := c.token.IsKeyword(); got != c.isKeyword {
			t.Errorf("%s.IsKeyword() is %t, expected %t", c.token, got, c.isKeyword)
		}

		if got := c.token.IsOperator(); got != c.isOperator {
			t.Errorf("%s.IsOperator() is %t, expected %t", c.token, got, c.isOperator)
		}
	}
}

func TestCheckKeywordToken(t *testing.T) {
	tests := []struct {
		token    string
		expected Token
	}{
		{SLet, Let},
		{SReturn, Return},
		{SFn, Fn},
		{SIf, If},
		{SElse, Else},
		{SImport, Import},
		{SNull, Null},
		{SFalse, False},
		{STrue, True},
		{"foobar", Identifier},
	}

	for _, c := range tests {
		token := CheckKeywordToken(c.token)
		if token != c.expected {
			t.Errorf("wrong token of %s, expected %s, got %s", c.token, c.expected, token)
		}
	}
}

func TestCheckOperatorToken(t *testing.T) {
	tests := []struct {
		token    string
		expected Token
	}{
		{"+", Plus},
		{"-", Minus},
		{"*", Asterisk},
		{"/", Slash},
		{"!", Bang},
		{"==", EQ},
		{"!=", NE},
		{"<", LT},
		{"<=", LE},
		{">", GT},
		{">=", GE},
		{"&&", AND},
		{"||", OR},
		{"=", Assign},
		{".", Period},
		{",", Comma},
		{":", Colon},
		{";", Semicolon},
		{"(", LParen},
		{")", RParen},
		{"{", LBrace},
		{"}", RBrace},
		{"[", LBracket},
		{"]", RBracket},
		{"foobar", Illegal},
	}

	for _, c := range tests {
		token := CheckOperatorToken(c.token)
		if token != c.expected {
			t.Errorf("wrong token of %s, expected %s, got %s", c.token, c.expected, token)
		}
	}
}
