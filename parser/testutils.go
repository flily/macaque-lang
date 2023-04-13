package parser

import (
	"strings"

	"github.com/flily/macaque-lang/ast"
	"github.com/flily/macaque-lang/token"
)

func makeMultilines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func makeProgram(statements ...ast.Statement) *ast.Program {
	program := &ast.Program{
		Statements: statements,
	}

	return program
}

func makeExpressionStatement(expressions ...ast.Expression) *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Expressions: &ast.ExpressionList{
			Expressions: expressions,
		},
	}

	return stmt
}

func makeLetStatement(identifers *ast.IdentifierList, expressions *ast.ExpressionList) *ast.LetStatement {
	stmt := &ast.LetStatement{
		Identifiers: identifers,
		Expressions: expressions,
	}

	return stmt
}

func makeIdentifierList(names ...string) *ast.IdentifierList {
	list := &ast.IdentifierList{}
	for _, name := range names {
		id := &ast.Identifier{
			Value: name,
		}
		list.Identifiers = append(list.Identifiers, id)
	}

	return list
}

func makeExpressionList(expressions ...ast.Expression) *ast.ExpressionList {
	list := &ast.ExpressionList{
		Expressions: expressions,
	}

	return list
}

func makePrefixExpression(operator string, right ast.Expression) *ast.PrefixExpression {
	expr := &ast.PrefixExpression{
		PrefixOperator: token.CheckOperatorToken(operator),
		Operand:        right,
	}

	return expr
}

func makeInfixExpression(operator string, left ast.Expression, right ast.Expression) *ast.InfixExpression {
	expr := &ast.InfixExpression{
		LeftOperand:  left,
		Operator:     token.CheckOperatorToken(operator),
		RightOperand: right,
	}

	return expr
}
