package ast

import (
	"fmt"
	"strings"

	"github.com/flily/macaque-lang/token"
)

type ExpressionListItem struct {
	Expression Expression
	Comma      *token.TokenContext
}

func (i *ExpressionListItem) GetContext() *token.Context {
	c := token.JoinContext(
		i.Expression.GetContext(),
		i.Comma.ToContext(),
	)

	return c
}

func (i *ExpressionListItem) IsIdentifier() bool {
	_, ok := i.Expression.(*Identifier)
	return ok
}

func (i *ExpressionListItem) ToIdentifier() *IdentifierListItem {
	item := &IdentifierListItem{
		Identifier: i.Expression.(*Identifier),
		Comma:      i.Comma,
	}

	return item
}

type ExpressionList struct {
	Expressions []*ExpressionListItem
}

func (l *ExpressionList) expressionNode() {}

func (l *ExpressionList) CanonicalCode() string {
	elems := make([]string, len(l.Expressions))
	for i, item := range l.Expressions {
		elems[i] = item.Expression.CanonicalCode()
	}

	return strings.Join(elems, ", ")
}

func (l *ExpressionList) GetContext() *token.Context {
	ctxs := make([]*token.Context, len(l.Expressions))
	for i, n := range l.Expressions {
		ctxs[i] = n.GetContext()
	}

	return token.JoinContext(ctxs...)
}

func (l *ExpressionList) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *ExpressionList:
		if len(n.Expressions) == len(l.Expressions) {
			result = true
			for i, expr := range l.Expressions {
				if !expr.Expression.EqualTo(n.Expressions[i].Expression) {
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

func (l *ExpressionList) AddExpression(expr Expression, comma *token.TokenContext) {
	item := &ExpressionListItem{
		Expression: expr,
		Comma:      comma,
	}

	l.Expressions = append(l.Expressions, item)
}

func (l *ExpressionList) IsIdentifierList() bool {
	result := true
	for _, expr := range l.Expressions {
		if !expr.IsIdentifier() {
			result = false
			break
		}
	}

	return result
}

func (l *ExpressionList) ToIdentifierList() *IdentifierList {
	ids := make([]*IdentifierListItem, len(l.Expressions))
	for i, expr := range l.Expressions {
		ok := expr.IsIdentifier()
		if !ok {
			return nil
		}

		ids[i] = expr.ToIdentifier()
	}

	list := &IdentifierList{
		Identifiers: ids,
	}

	return list
}

type IdentifierListItem struct {
	Identifier *Identifier
	Comma      *token.TokenContext
}

func (i *IdentifierListItem) GetContext() *token.Context {
	c := token.JoinContext(
		i.Identifier.GetContext(),
		i.Comma.ToContext(),
	)

	return c
}

type IdentifierList struct {
	Identifiers []*IdentifierListItem
}

func (l *IdentifierList) expressionNode() {}

func (l *IdentifierList) CanonicalCode() string {
	elems := make([]string, len(l.Identifiers))

	for i, id := range l.Identifiers {
		elems[i] = id.Identifier.CanonicalCode()
	}

	return strings.Join(elems, ", ")
}

func (l *IdentifierList) GetContext() *token.Context {
	ctxs := make([]*token.Context, len(l.Identifiers))
	for i, n := range l.Identifiers {
		ctxs[i] = n.GetContext()
	}

	return token.JoinContext(ctxs...)
}

func (l *IdentifierList) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *IdentifierList:
		if len(n.Identifiers) == len(l.Identifiers) {
			result = true
			for i, id := range l.Identifiers {
				if !id.Identifier.EqualTo(n.Identifiers[i].Identifier) {
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

func (l *IdentifierList) AddIdentifier(id *Identifier, comma *token.TokenContext) {
	item := &IdentifierListItem{
		Identifier: id,
		Comma:      comma,
	}

	l.Identifiers = append(l.Identifiers, item)
}

type PrefixExpression struct {
	Prefix  *token.TokenContext
	Operand Expression
}

func (e *PrefixExpression) expressionNode() {}

func (e *PrefixExpression) CanonicalCode() string {
	s := fmt.Sprintf("(%s %s)",
		e.Prefix.Token.CodeName(),
		e.Operand.CanonicalCode())

	return s
}

func (e *PrefixExpression) GetContext() *token.Context {
	c := token.JoinContext(
		e.Prefix.ToContext(),
		e.Operand.GetContext(),
	)

	return c
}

func (e *PrefixExpression) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *PrefixExpression:
		if e.Prefix.Token == n.Prefix.Token {
			result = e.Operand.EqualTo(n.Operand)
		}
	}

	return result
}

type InfixExpression struct {
	LeftOperand  Expression
	Operator     *token.TokenContext
	RightOperand Expression
}

func (e *InfixExpression) expressionNode() {}

func (e *InfixExpression) CanonicalCode() string {
	s := fmt.Sprintf("(%s %s %s)",
		e.LeftOperand.CanonicalCode(),
		e.Operator.Token.CodeName(),
		e.RightOperand.CanonicalCode())

	return s
}

func (e *InfixExpression) GetContext() *token.Context {
	c := token.JoinContext(
		e.LeftOperand.GetContext(),
		e.Operator.ToContext(),
		e.RightOperand.GetContext(),
	)

	return c
}

func (e *InfixExpression) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *InfixExpression:
		if e.Operator.Token == n.Operator.Token {
			result = e.LeftOperand.EqualTo(n.LeftOperand) &&
				e.RightOperand.EqualTo(n.RightOperand)
		}
	}

	return result
}

// IndexExpression get the value of the index of an array or a hash.
// - Base[Expression]
// - Base.Identifier
type IndexExpression struct {
	Base     Expression
	Operator *token.TokenContext
	Index    Expression
	End      *token.TokenContext
}

func (e *IndexExpression) expressionNode() {}

func (e *IndexExpression) CanonicalCode() string {
	s := fmt.Sprintf("(%s[%s])",
		e.Base.CanonicalCode(),
		e.Index.CanonicalCode())

	return s
}

func (e *IndexExpression) GetContext() *token.Context {
	c := token.JoinContext(
		e.Base.GetContext(),
		e.Operator.ToContext(),
		e.Base.GetContext(),
		e.End.ToContext(),
	)

	return c
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

// CallExpression
// Callable()
// Callable::Member()
// fn()
type CallExpression struct {
	Base      Expression
	Token     *token.TokenContext
	Member    *Identifier
	LParen    *token.TokenContext
	Args      *ExpressionList
	RParen    *token.TokenContext
	Recursion bool
}

func (e *CallExpression) expressionNode() {}

func (e *CallExpression) CanonicalCode() string {
	var result string

	switch e.Token.GetToken() {
	case token.Nil:
		result = fmt.Sprintf("%s(%s)",
			e.Base.CanonicalCode(),
			e.Args.CanonicalCode(),
		)

	case token.DualColon:
		result = fmt.Sprintf("%s::%s(%s)",
			e.Base.CanonicalCode(),
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

func (e *CallExpression) GetContext() *token.Context {
	c := token.JoinContext(
		GetContext(e.Base),
		e.Token.ToContext(),
		e.Member.GetContext(),
		e.LParen.ToContext(),
		e.Args.GetContext(),
		e.RParen.ToContext(),
	)

	return c
}

func (e *CallExpression) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *CallExpression:
		if e.Base == nil {
			result = e.Base == nil
		} else {
			result = e.Base.EqualTo(n.Base)
		}

		result = result && e.Args.EqualTo(n.Args)
		result = result && e.Token.Token == n.Token.Token && e.Recursion == n.Recursion

		if e.Member != nil {
			result = result && e.Member.EqualTo(n.Member)
		}
	}

	return result
}

type IfExpression struct {
	If          *token.TokenContext
	LParen      *token.TokenContext
	Condition   Expression
	RParen      *token.TokenContext
	Else        *token.TokenContext
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

func (e *IfExpression) GetContext() *token.Context {
	c := token.JoinContext(
		e.If.ToContext(),
		e.LParen.ToContext(),
		e.Condition.GetContext(),
		e.RParen.ToContext(),
		e.Consequence.GetContext(),
		e.Else.ToContext(),
		GetContext(e.Alternative),
	)

	return c
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
