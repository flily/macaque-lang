package ast

import (
	"errors"
)

const (
	WalkForAllNodes       = 0
	StopWalkingOnChildren = 1
	StopWalking           = 2
)

var (
	errStopWalking = errors.New("stop walking")
)

type Visitor interface {
	Visit(Node) int
}

// Walk traverses the AST in depth-first order.
func Walk(node Node, visitor Visitor) {
	defer func() {
		if r := recover(); r != nil {
			if r != errStopWalking {
				panic(r)
			}
		}
	}()

	doWalk(node, visitor)
}

func doWalk(node Node, v Visitor) {
	if node == nil {
		return
	}

	c := v.Visit(node)
	switch c {
	case StopWalkingOnChildren:
		return

	case StopWalking:
		panic(errStopWalking)
	}

	switch n := node.(type) {
	case *Program:
		walkOnStatements(n.Statements, v)

	case *Identifier, *IntegerLiteral, *FloatLiteral, *StringLiteral, *BooleanLiteral, *NullLiteral:
		// do nothing, node has already sent to visitor

	case *ArrayLiteral:
		walkOnExpressions(n.Expressions, v)

	case *HashLiteral:
		for _, pair := range n.Pairs {
			doWalk(pair.Key, v)
			doWalk(pair.Value, v)
		}

	case *ExpressionList:
		for _, expr := range n.Expressions {
			doWalk(expr.Expression, v)
		}

	case *IdentifierList:
		for _, id := range n.Identifiers {
			doWalk(id.Identifier, v)
		}

	case *PrefixExpression:
		doWalk(n.Operand, v)

	case *InfixExpression:
		doWalk(n.LeftOperand, v)
		doWalk(n.RightOperand, v)

	case *IndexExpression:
		doWalk(n.Base, v)
		doWalk(n.Index, v)

	case *CallExpression:
		doWalk(n.Base, v)
		doWalk(n.Member, v) // it is safe to call doWalk with nil
		doWalk(n.Args, v)

	case *IfExpression:
		doWalk(n.Condition, v)
		doWalk(n.Consequence, v)
		doWalk(n.Alternative, v)

	case *LetStatement:
		doWalk(n.Identifiers, v)
		doWalk(n.Expressions, v)

	case *ReturnStatement:
		doWalk(n.Expressions, v)

	case *IfStatement:
		doWalk(n.Expression, v)

	case *BlockStatement:
		walkOnStatements(n.Statements, v)

	case *ExpressionStatement:
		doWalk(n.Expressions, v)

	case *ImportStatement:

	}
}

func walkOnExpressions(list *ExpressionList, v Visitor) {
	for _, expr := range list.Expressions {
		doWalk(expr.Expression, v)
	}
}

func walkOnStatements(list []Statement, v Visitor) {
	for _, stmt := range list {
		doWalk(stmt, v)
	}
}
