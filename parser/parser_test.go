package parser

import (
	"testing"

	"github.com/flily/macaque-lang/token"
)

func TestGetPrecedence(t *testing.T) {
	tests := []struct {
		token    token.Token
		expected int
	}{
		{token.Illegal, 0},
		{token.Bang, 0},
		{token.Plus, PrecedenceSum},
		{token.Minus, PrecedenceSum},
		{token.Asterisk, PrecedenceProduct},
		{token.Slash, PrecedenceProduct},
		{token.EQ, PrecedenceComparisonEqual},
		{token.NE, PrecedenceComparisonEqual},
		{token.LT, PrecedenceComparisonLessGreater},
		{token.LE, PrecedenceComparisonLessGreater},
		{token.GT, PrecedenceComparisonLessGreater},
		{token.GE, PrecedenceComparisonLessGreater},
		{token.AND, PrecedenceLogicalAND},
		{token.OR, PrecedenceLogicalOR},
		{token.Assign, 0},
		{token.Period, PrecedenceIndex},
		{token.Comma, 0},
		{token.Colon, 0},
		{token.Semicolon, 0},
		{token.DualColon, PrecedenceCall},
		{token.LParen, PrecedenceCall},
		{token.RParen, 0},
		{token.LBrace, 0},
		{token.RBrace, 0},
		{token.LBracket, PrecedenceIndex},
		{token.RBracket, 0},
		{token.Identifier, 0},
		{token.Integer, 0},
		{token.Float, 0},
		{token.String, 0},
		{token.Null, 0},
		{token.Modulo, PrecedenceProduct},
		{token.BITAND, PrecedenceBitwiseAND},
		{token.BITOR, PrecedenceBitwiseOR},
		{token.BITXOR, PrecedenceBitwiseXOR},
		{token.BITNOT, PrecedencePrefix},
	}

	for _, test := range tests {
		if GetPrecedence(test.token) != test.expected {
			t.Errorf("GetPrecedence(%s) expected %d, got %d", test.token, test.expected, GetPrecedence(test.token))
		}
	}
}

func TestIsInfixOperator(t *testing.T) {
	tests := []struct {
		token    token.Token
		expected bool
	}{
		{token.Illegal, false},
		{token.Bang, false},
		{token.Plus, true},
		{token.Minus, true},
		{token.Asterisk, true},
		{token.Slash, true},
		{token.EQ, true},
		{token.NE, true},
		{token.LT, true},
		{token.LE, true},
		{token.GT, true},
		{token.GE, true},
		{token.AND, true},
		{token.OR, true},
		{token.Assign, false},
		{token.Period, false},
		{token.Comma, false},
		{token.Colon, false},
		{token.Semicolon, false},
		{token.LParen, false},
		{token.RParen, false},
		{token.LBrace, false},
		{token.RBrace, false},
		{token.LBracket, false},
		{token.RBracket, false},
		{token.Illegal, false},
	}

	for _, c := range tests {
		if got := IsInfixOperator(c.token); got != c.expected {
			t.Errorf("IsInfixOperator(%s) is %t, expected %t", c.token, got, c.expected)
		}
	}
}
