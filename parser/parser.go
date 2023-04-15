package parser

import (
	"github.com/flily/macaque-lang/token"
)

const (
	_ int = iota
	PrecedenceLowest
	PrecedenceLogicalOR             // ||
	PrecedenceLogicalAND            // &&
	PrecedenceComparisonEqual       // == !=
	PrecedenceComparisonLessGreater // > < >= <=
	PrecedenceSum                   // + -
	PrecedenceProduct               // * / %
	PrecedenceBitwiseOR             // |
	PrecedenceBitwiseXOR            // ^
	PrecedenceBitwiseAND            // &
	PrecedencePrefix                // -X or !X or ~X
	PrecedenceCall                  // myFunction(X)
	PrecedenceIndex                 // array[index]
)

var precedenceMap = map[token.Token]int{
	token.AND:      PrecedenceLogicalAND,
	token.OR:       PrecedenceLogicalOR,
	token.EQ:       PrecedenceComparisonEqual,
	token.NE:       PrecedenceComparisonEqual,
	token.LT:       PrecedenceComparisonLessGreater,
	token.GT:       PrecedenceComparisonLessGreater,
	token.LE:       PrecedenceComparisonLessGreater,
	token.GE:       PrecedenceComparisonLessGreater,
	token.Plus:     PrecedenceSum,
	token.Minus:    PrecedenceSum,
	token.Asterisk: PrecedenceProduct,
	token.Slash:    PrecedenceProduct,
	token.Modulo:   PrecedenceProduct,
	token.BITAND:   PrecedenceBitwiseAND,
	token.BITOR:    PrecedenceBitwiseOR,
	token.BITXOR:   PrecedenceBitwiseXOR,
	token.BITNOT:   PrecedencePrefix,
	token.LParen:   PrecedenceCall,
	token.Colon:    PrecedenceCall,
	token.LBracket: PrecedenceIndex,
	token.Period:   PrecedenceIndex,
}

func GetPrecedence(token token.Token) int {
	if precedence, ok := precedenceMap[token]; ok {
		return precedence
	}

	return 0
}

func IsInfixOperator(t token.Token) bool {
	return token.Plus <= t && t <= token.BITOR
}
