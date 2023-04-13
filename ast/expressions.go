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

func (l *ExpressionList) Children() []Node {
	nodes := make([]Node, len(l.Expressions))
	for i, expr := range l.Expressions {
		nodes[i] = expr
	}

	return nodes
}

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

type IdentifierList struct {
	Identifiers []*Identifier
}

func (l *IdentifierList) expressionNode() {}

func (l *IdentifierList) Children() []Node {
	nodes := make([]Node, len(l.Identifiers))

	for i, id := range l.Identifiers {
		nodes[i] = id
	}

	return nodes
}

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

func (e *PrefixExpression) Children() []Node {
	return []Node{e.Operand}
}

func (e *PrefixExpression) CanonicalCode() string {
	s := fmt.Sprintf("%s %s",
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

func (e *InfixExpression) Children() []Node {
	return []Node{e.LeftOperand, e.RightOperand}
}

func (e *InfixExpression) CanonicalCode() string {
	s := fmt.Sprintf("%s %s %s",
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

func (e *IndexExpression) Children() []Node {
	return []Node{e.Base, e.Index}
}

func (e *IndexExpression) CanonicalCode() string {
	s := fmt.Sprintf("%s[%s]",
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
	Callable Expression
	Colon    token.Token
	Member   *Identifier
	Args     *ExpressionList
}

func (e *CallExpression) expressionNode() {}

func (e *CallExpression) Children() []Node {
	elementCount := 2
	if e.Member != nil {
		elementCount = 3
	}

	nodes := make([]Node, elementCount)
	nodes[0] = e.Callable
	if e.Member != nil {
		nodes[1] = e.Member
	}

	nodes[elementCount-1] = e.Args
	return nodes
}

func (e *CallExpression) CanonicalCode() string {
	var result string

	if e.Member != nil {
		result = fmt.Sprintf("%s:%s(%s)",
			e.Callable.CanonicalCode(),
			e.Member.CanonicalCode(),
			e.Args.CanonicalCode(),
		)
	} else {
		result = fmt.Sprintf("%s(%s)",
			e.Callable.CanonicalCode(),
			e.Args.CanonicalCode(),
		)
	}

	return result
}

func (e *CallExpression) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *CallExpression:
		result = e.Callable.EqualTo(n.Callable) &&
			e.Args.EqualTo(n.Args)
		if e.Member != nil {
			result = result && e.Member.EqualTo(n.Member)
		}
	}

	return result
}

type GroupExpression struct {
	Expression Expression
}

func (e *GroupExpression) expressionNode() {}

func (e *GroupExpression) Children() []Node {
	return []Node{e.Expression}
}

func (e *GroupExpression) CanonicalCode() string {
	return fmt.Sprintf("(%s)", e.Expression.CanonicalCode())
}

func (e *GroupExpression) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *GroupExpression:
		result = e.Expression.EqualTo(n.Expression)
	}

	return result
}
