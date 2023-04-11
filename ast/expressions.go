package ast

import (
	"fmt"
	"strings"

	"github.com/flily/macaque-lang/token"
)

type ExpressionList struct {
	Expressions []Expression
}

func (e *ExpressionList) expressionNode() {}

func (e *ExpressionList) Children() []Node {
	nodes := make([]Node, len(e.Expressions))
	for i, exp := range e.Expressions {
		nodes[i] = exp
	}

	return nodes
}

func (e *ExpressionList) CanonicalCode() string {
	elems := make([]string, len(e.Expressions))
	for i, exp := range e.Expressions {
		elems[i] = exp.CanonicalCode()
	}

	return strings.Join(elems, ", ")
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

type CallExpression struct {
	Callable Expression
	Colon    token.Token
	Member   *Identifier
	Args     []Expression
}

func (e *CallExpression) expressionNode() {}

func (e *CallExpression) Children() []Node {
	hasMember := 1
	if e.Member != nil {
		hasMember = 2
	}

	nodes := make([]Node, len(e.Args)+hasMember)
	nodes[0] = e.Callable
	if e.Member != nil {
		nodes[1] = e.Member
	}

	for i, arg := range e.Args {
		nodes[i+hasMember] = arg
	}

	return nodes
}

func (e *CallExpression) CanonicalCode() string {
	args := make([]string, len(e.Args))
	for i, arg := range e.Args {
		args[i] = arg.CanonicalCode()
	}

	var result string
	if e.Member != nil {
		result = fmt.Sprintf("%s:%s(%s)",
			e.Callable.CanonicalCode(),
			e.Member.CanonicalCode(),
			strings.Join(args, ", "),
		)
	} else {
		result = fmt.Sprintf("%s(%s)",
			e.Callable.CanonicalCode(),
			strings.Join(args, ", "),
		)
	}

	return result
}
