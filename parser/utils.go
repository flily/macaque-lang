package parser

import (
	"fmt"
	"strings"

	"github.com/flily/macaque-lang/ast"
	"github.com/flily/macaque-lang/token"
)

func makeMultilines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func program(statements ...ast.Statement) *ast.Program {
	program := &ast.Program{
		Statements: statements,
	}

	return program
}

func l(v interface{}) ast.Expression {
	var r ast.Expression

	if v == nil {
		return &ast.NullLiteral{}
	}

	switch v := v.(type) {
	case int:
		r = &ast.IntegerLiteral{
			Value:   int64(v),
			Content: fmt.Sprintf("%d", v),
		}

	case string:
		r = &ast.StringLiteral{
			Content: v,
			Value:   ConvertString(v),
		}

	case bool:
		r = &ast.BooleanLiteral{
			Value: v,
		}
	}

	return r
}

func float(s string) ast.Expression {
	return &ast.FloatLiteral{
		Value:   ConvertFloat(s),
		Content: s,
	}
}

// func integer(s string) ast.Expression {
// 	return &ast.IntegerLiteral{
// 		Value:   ConvertInteger(s),
// 		Content: s,
// 	}
// }

func punct(s token.Token) *token.TokenContext {
	t := &token.TokenContext{
		Token: s,
	}

	return t
}

func expr(expressions ...ast.Expression) *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Expressions: exprList(expressions...),
	}

	return stmt
}

func let(identifers *ast.IdentifierList, expressions *ast.ExpressionList) *ast.LetStatement {
	stmt := &ast.LetStatement{
		Identifiers: identifers,
		Expressions: expressions,
	}

	return stmt
}

func ret(expressions ...ast.Expression) *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Expressions: exprList(expressions...),
	}

	return stmt
}

func id(name string) *ast.Identifier {
	id := &ast.Identifier{
		Value: name,
	}

	return id
}

func idList(names ...string) *ast.IdentifierList {
	list := &ast.IdentifierList{}
	for _, name := range names {
		item := &ast.IdentifierListItem{
			Identifier: id(name),
		}
		list.Identifiers = append(list.Identifiers, item)
	}

	return list
}

func exprList(expressions ...ast.Expression) *ast.ExpressionList {
	e := make([]*ast.ExpressionListItem, len(expressions))
	for i, expr := range expressions {
		e[i] = &ast.ExpressionListItem{
			Expression: expr,
		}
	}

	list := &ast.ExpressionList{
		Expressions: e,
	}

	return list
}

func prefix(operator string, right ast.Expression) *ast.PrefixExpression {
	expr := &ast.PrefixExpression{
		Prefix: &token.TokenContext{
			Token: token.CheckOperatorToken(operator),
		},
		Operand: right,
	}

	return expr
}

func infix(operator string, left ast.Expression, right ast.Expression) *ast.InfixExpression {
	expr := &ast.InfixExpression{
		LeftOperand: left,
		Operator: &token.TokenContext{
			Token: token.CheckOperatorToken(operator),
		},
		RightOperand: right,
	}

	return expr
}

func array(elements ...ast.Expression) *ast.ArrayLiteral {
	expr := &ast.ArrayLiteral{
		Expressions: exprList(elements...),
	}

	return expr
}

func pair(key ast.Expression, value ast.Expression) *ast.HashItem {
	pair := &ast.HashItem{
		Key:   key,
		Value: value,
	}

	return pair
}

func hash(pairs ...*ast.HashItem) *ast.HashLiteral {
	expr := &ast.HashLiteral{
		Pairs: pairs,
	}

	return expr
}

func ifexp(condition ast.Expression, consequence *ast.BlockStatement, alternative ast.BlockStatementNode) *ast.IfExpression {
	expr := &ast.IfExpression{
		Condition:   condition,
		Consequence: consequence,
		Alternative: alternative,
	}

	return expr
}

func elseif(condition ast.Expression, consequence *ast.BlockStatement, alternative ast.BlockStatementNode) *ast.IfStatement {
	expr := &ast.IfStatement{
		Expression: ifexp(condition, consequence, alternative),
	}

	return expr
}

func block(statements ...ast.Statement) *ast.BlockStatement {
	stmt := &ast.BlockStatement{
		Statements: statements,
	}

	return stmt
}

func fn(args *ast.IdentifierList, body *ast.BlockStatement) *ast.FunctionLiteral {
	expr := &ast.FunctionLiteral{
		Arguments: args,
		Body:      body,
	}

	return expr
}
