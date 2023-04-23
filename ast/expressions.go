package ast

import (
	"fmt"
	"strings"

	"github.com/flily/macaque-lang/token"
)

type ExpressionList struct {
	Expressions []Expression
}

func (l *ExpressionList) expressionNode() {}

func (l *ExpressionList) CanonicalCode() string {
	elems := make([]string, len(l.Expressions))
	for i, expr := range l.Expressions {
		elems[i] = expr.CanonicalCode()
	}

	return strings.Join(elems, ", ")
}

func (l *ExpressionList) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *ExpressionList:
		if len(n.Expressions) == len(l.Expressions) {
			result = true
			for i, expr := range l.Expressions {
				if !expr.EqualTo(n.Expressions[i]) {
					result = false
					break
				}
			}
		}
	}

	return result
}

func (l *ExpressionList) Length() int {
	if l != nil {
		return len(l.Expressions)
	}

	return 0
}

func (l *ExpressionList) AddExpression(expr Expression) {
	l.Expressions = append(l.Expressions, expr)
}

func (l *ExpressionList) IsIdentifierList() bool {
	result := true
	for _, expr := range l.Expressions {
		if _, ok := expr.(*Identifier); !ok {
			result = false
			break
		}
	}

	return result
}

func (l *ExpressionList) ToIdentifierList() *IdentifierList {
	ids := make([]*Identifier, len(l.Expressions))
	for i, expr := range l.Expressions {
		id, ok := expr.(*Identifier)
		if !ok {
			return nil
		}

		ids[i] = id
	}

	list := &IdentifierList{
		Identifiers: ids,
	}

	return list
}

type IdentifierList struct {
	Identifiers []*Identifier
}

func (l *IdentifierList) expressionNode() {}

func (l *IdentifierList) CanonicalCode() string {
	elems := make([]string, len(l.Identifiers))

	for i, id := range l.Identifiers {
		elems[i] = id.CanonicalCode()
	}

	return strings.Join(elems, ", ")
}

func (l *IdentifierList) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *IdentifierList:
		if len(n.Identifiers) == len(l.Identifiers) {
			result = true
			for i, id := range l.Identifiers {
				if !id.EqualTo(n.Identifiers[i]) {
					result = false
					break
				}
			}
		}
	}

	return result
}

func (l *IdentifierList) Length() int {
	if l != nil {
		return len(l.Identifiers)
	}

	return 0
}

func (l *IdentifierList) AddIdentifier(id *Identifier) {
	l.Identifiers = append(l.Identifiers, id)
}

type PrefixExpression struct {
	PrefixOperator token.Token
	Operand        Expression
}

func (e *PrefixExpression) expressionNode() {}

func (e *PrefixExpression) CanonicalCode() string {
	s := fmt.Sprintf("(%s %s)",
		e.PrefixOperator.CodeName(),
		e.Operand.CanonicalCode())

	return s
}

func (e *PrefixExpression) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *PrefixExpression:
		if e.PrefixOperator == n.PrefixOperator {
			result = e.Operand.EqualTo(n.Operand)
		}
	}

	return result
}

type InfixExpression struct {
	LeftOperand  Expression
	Operator     token.Token
	RightOperand Expression
}

func (e *InfixExpression) expressionNode() {}

func (e *InfixExpression) CanonicalCode() string {
	s := fmt.Sprintf("(%s %s %s)",
		e.LeftOperand.CanonicalCode(),
		e.Operator.CodeName(),
		e.RightOperand.CanonicalCode())

	return s
}

func (e *InfixExpression) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *InfixExpression:
		if e.Operator == n.Operator {
			result = e.LeftOperand.EqualTo(n.LeftOperand) &&
				e.RightOperand.EqualTo(n.RightOperand)
		}
	}

	return result
}

type IndexExpression struct {
	Base     Expression
	Operator token.Token
	Index    Expression
	End      token.Token
}

func (e *IndexExpression) expressionNode() {}

func (e *IndexExpression) CanonicalCode() string {
	s := fmt.Sprintf("(%s[%s])",
		e.Base.CanonicalCode(),
		e.Index.CanonicalCode())

	return s
}

func (e *IndexExpression) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *IndexExpression:
		result = e.Base.EqualTo(n.Base) &&
			e.Index.EqualTo(n.Index)
	}

	return result
}

type CallExpression struct {
	Callable  Expression
	Token     token.Token
	Member    *Identifier
	Args      *ExpressionList
	Recursion bool
}

func (e *CallExpression) expressionNode() {}

func (e *CallExpression) CanonicalCode() string {
	var result string

	switch e.Token {
	case token.LParen:
		result = fmt.Sprintf("%s(%s)",
			e.Callable.CanonicalCode(),
			e.Args.CanonicalCode(),
		)

	case token.DualColon:
		result = fmt.Sprintf("%s::%s(%s)",
			e.Callable.CanonicalCode(),
			e.Member.CanonicalCode(),
			e.Args.CanonicalCode(),
		)

	case token.Fn:
		result = fmt.Sprintf("fn(%s)",
			e.Args.CanonicalCode(),
		)
	}

	return result
}

func (e *CallExpression) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *CallExpression:
		if e.Callable == nil {
			result = e.Callable == nil
		} else {
			result = e.Callable.EqualTo(n.Callable)
		}

		result = result && e.Args.EqualTo(n.Args)
		result = result && e.Token == n.Token && e.Recursion == n.Recursion

		if e.Member != nil {
			result = result && e.Member.EqualTo(n.Member)
		}
	}

	return result
}

type IfExpression struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative BlockStatementNode
}

func (e *IfExpression) expressionNode() {}

func (e *IfExpression) CanonicalCode() string {
	result := fmt.Sprintf("if ( %s ) %s",
		e.Condition.CanonicalCode(),
		e.Consequence.CanonicalCode(),
	)

	if e.Alternative != nil {
		result += fmt.Sprintf(" else %s", e.Alternative.CanonicalCode())
	}

	return result
}

func (e *IfExpression) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *IfExpression:
		result = e.Condition.EqualTo(n.Condition) &&
			e.Consequence.EqualTo(n.Consequence)

		if result && e.Alternative != nil {
			result = result && e.Alternative.EqualTo(n.Alternative)
		}
	}

	return result
}
