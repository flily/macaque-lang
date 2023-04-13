package parser

import (
	"github.com/flily/macaque-lang/token"
)

const (
	_ int = iota
	PrecedenceLowest
	PrecedenceLogical               // && ||
	PrecedenceComparisonEqual       // == !=
	PrecedenceComparisonLessGreater // > or <
	PrecedenceSum                   // + -
	PrecedenceProduct               // * / %
	PrecedenceBitwise               // & | ^ ~
	PrecedencePrefix                // -X or !X
	PrecedenceCall                  // myFunction(X)
	PrecedenceIndex                 // array[index]
)

var precedenceMap = map[token.Token]int{
	token.AND:      PrecedenceLogical,
	token.OR:       PrecedenceLogical,
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
	token.BITAND:   PrecedenceBitwise,
	token.BITOR:    PrecedenceBitwise,
	token.BITXOR:   PrecedenceBitwise,
	token.BITNOT:   PrecedenceBitwise,
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
