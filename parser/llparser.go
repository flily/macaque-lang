package parser

import (
	"github.com/flily/macaque-lang/ast"
	"github.com/flily/macaque-lang/errors"
	"github.com/flily/macaque-lang/lex"
	"github.com/flily/macaque-lang/token"
)

type LLParser struct {
	container *CodeContainer
}

func NewLLParser(scanner lex.Scanner) *LLParser {
	p := &LLParser{
		container: NewContainer(scanner),
	}

	return p
}

func (p *LLParser) ReadTokens() error {
	return p.container.ReadTokens()
}

func (p *LLParser) Parse() (*ast.Program, error) {
	return p.parseProgram()
}

func (p *LLParser) expect(token token.Token) error {
	current := p.container.Current()
	if current.Token != token {
		return p.makeSyntaxError("expected token %s, but got %s", token, current.Token)
	}

	return nil
}

func (p *LLParser) current() *lex.LexicalElement {
	return p.container.Current()
}

func (p *LLParser) peek(offset int) *lex.LexicalElement {
	return p.container.Peek(offset)
}

func (p *LLParser) nextToken() *lex.LexicalElement {
	return p.container.Next()
}

func (p *LLParser) skipToken(token token.Token) error {
	if err := p.expect(token); err != nil {
		return err
	}

	p.container.Next()
	return nil
}

func (p *LLParser) skipComment() {
	for {
		current := p.container.Current()
		if current == nil || current.Token != token.Comment {
			break
		}

		p.container.Next()
	}
}

func (p *LLParser) makeSyntaxError(format string, args ...interface{}) *errors.SyntaxError {
	current := p.container.Current()
	ctx := current.Position.MakeContext()
	return ctx.NewSyntaxError(format, args...)
}

func (p *LLParser) parseProgram() (*ast.Program, error) {
	program := ast.NewEmptyProgram()
	current := p.current()

	for current != nil && current.Token != token.EOF {
		var stmt ast.Statement
		var err error

		switch current.Token {
		case token.Let:
			stmt, err = p.parseLetStatement()

		case token.Comment:
			// skip comment
			p.skipComment()
			continue

		default:
			return nil, p.makeSyntaxError("unexpected token: %s", current.Token)
		}

		if err != nil {
			return nil, err
		}

		program.AddStatement(stmt)
		current = p.current()
	}

	return program, nil
}

func (p *LLParser) parseLetStatement() (*ast.LetStatement, error) {
	if err := p.skipToken(token.Let); err != nil {
		return nil, err
	}

	idList, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}

	if err := p.skipToken(token.Assign); err != nil {
		return nil, err
	}

	exprList, err := p.parseExpessionList()
	if err != nil {
		return nil, err
	}

	_ = p.skipToken(token.Semicolon)

	stmt := &ast.LetStatement{
		Identifiers: idList,
		Expressions: exprList,
	}

	return stmt, nil
}

func (p *LLParser) parseIdentifierList() (*ast.IdentifierList, error) {
	var err error
	list := &ast.IdentifierList{}

	for {
		id, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}

		list.AddIdentifier(id)

		if err := p.skipToken(token.Comma); err != nil {
			break
		}
	}

	return list, err
}

// identifier = [identifier-prefix] ( ALPHA / "_" ) *( ALPHA / DIGIT / "_" ) [identifier-suffix]
func (p *LLParser) parseIdentifier() (*ast.Identifier, error) {
	var id *ast.Identifier
	var err error

	found := false
	for !found {
		current := p.current()
		switch current.Token {
		case token.Identifier:
			id = ast.NewIdentifier(current.Content, current.Position)
			p.nextToken()
			found = true

		case token.Comment:
			p.skipComment()

		default:
			found = true
			err = p.expect(token.Identifier)
		}
	}

	return id, err
}

func (p *LLParser) parseExpessionList() (*ast.ExpressionList, error) {
	var err error
	list := &ast.ExpressionList{}

	for {
		exp, err := p.parseExpression(PrecedenceLowest)
		if err != nil {
			return nil, err
		}

		list.AddExpression(exp)

		if err := p.skipToken(token.Comma); err != nil {
			break
		}
	}

	return list, err
}

func (p *LLParser) parseExpression(precedence int) (ast.Expression, error) {
	var left ast.Expression
	var err error

	current := p.current()

	switch current.Token {
	case token.LParen:
		left, err = p.parseGroupExpression()

	case token.Integer, token.Float, token.String, token.True, token.False, token.Null:
		left, err = p.parseLiteral()

	case token.Bang, token.Minus:
		left, err = p.parsePrefixExpression()
	}

	current = p.current()
	if current != nil && current.Token.IsOperator() {
		return p.parseInfixExpression(left, precedence)
	}

	return left, err
}

func (p *LLParser) parseLiteral() (ast.Expression, error) {
	var expr ast.Expression
	var err error

	current := p.current()
	switch current.Token {
	case token.Integer:
		expr = NewInteger(current)

	case token.Float:
		expr = NewFloat(current)

	case token.String:
		expr = NewString(current)

	case token.True, token.False:
		expr = NewBoolean(current)

	case token.Null:
		expr = NewNull(current)

	}

	p.nextToken()
	return expr, err
}

func (p *LLParser) parsePrefixExpression() (*ast.PrefixExpression, error) {
	operator := p.current()
	p.nextToken()

	operand, err := p.parseExpression(PrecedencePrefix)
	if err != nil {
		return nil, err
	}

	expr := &ast.PrefixExpression{
		PrefixOperator: operator.Token,
		Operand:        operand,
	}

	return expr, nil

}

func (p *LLParser) parseInfixExpression(left ast.Expression, precedence int) (ast.Expression, error) {
	operator := p.current()
	operatorPrecedence := GetPrecedence(operator.Token)
	if operatorPrecedence <= precedence {
		return left, nil
	}

	p.nextToken()
	right, err := p.parseExpression(operatorPrecedence)
	if err != nil {
		return nil, err
	}

	expr := &ast.InfixExpression{
		LeftOperand:  left,
		Operator:     operator.Token,
		RightOperand: right,
	}

	return expr, nil
}

func (p *LLParser) parseCallExpression() (*ast.CallExpression, error) {
	return nil, nil
}

func (p *LLParser) parseIndexExpression() (*ast.IndexExpression, error) {
	return nil, nil
}

func (p *LLParser) parseGroupExpression() (*ast.GroupExpression, error) {
	if err := p.skipToken(token.LParen); err != nil {
		return nil, err
	}

	expr, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	if err := p.skipToken(token.RParen); err != nil {
		return nil, err
	}

	group := &ast.GroupExpression{
		Expression: expr,
	}

	return group, nil
}
