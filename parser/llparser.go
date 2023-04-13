package parser

import (
	"github.com/flily/macaque-lang/ast"
	"github.com/flily/macaque-lang/errors"
	"github.com/flily/macaque-lang/lex"
	"github.com/flily/macaque-lang/token"
)

var expressionFirstSet = [...]bool{
	token.Identifier: true,
	token.Integer:    true,
	token.Float:      true,
	token.String:     true,
	token.True:       true,
	token.False:      true,
	token.Null:       true,
	token.LParen:     true,
	token.Minus:      true,
	token.Bang:       true,
}

func isExpressionFirstSet(token token.Token) bool {
	return expressionFirstSet[token]
}

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

func (p *LLParser) currentToken() token.Token {
	current := p.container.Current()
	return current.Token
}

// func (p *LLParser) peek(offset int) *lex.LexicalElement {
// 	return p.container.Peek(offset)
// }

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

		case token.Null, token.False, token.True, token.Integer, token.Float, token.String,
			token.Identifier, token.Minus, token.Bang, token.LParen:
			stmt, err = p.parseExpressionStatement()

		default:
			return nil, p.makeSyntaxError("unexpected token in PROGRAM: %s", current.Token)
		}

		if err != nil {
			return nil, err
		}

		program.AddStatement(stmt)
		current = p.current()
	}

	return program, nil
}

// let-stmt
// => "let" identifier-list "=" expression-list ";"
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

	exprList, err := p.parseExpressionList()
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

// expression-stmt
// => expression-list ";"
func (p *LLParser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	exprList, err := p.parseExpressionList()
	if err != nil {
		return nil, err
	}

	_ = p.skipToken(token.Semicolon)

	stmt := &ast.ExpressionStatement{
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

		if p.expect(token.Identifier) != nil {
			break
		}
	}

	return list, err
}

// identifier
// => [identifier-prefix] ( ALPHA / "_" ) *( ALPHA / DIGIT / "_" ) [identifier-suffix]
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

// expression-list
// => expression *( "," expression ) [","]
func (p *LLParser) parseExpressionList() (*ast.ExpressionList, error) {
	var err error
	list := &ast.ExpressionList{}

	for isExpressionFirstSet(p.currentToken()) {
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

	case token.Identifier:
		left, err = p.parseIdentifier()

	case token.Bang, token.Minus:
		left, err = p.parsePrefixExpression()

	default:
		err = p.makeSyntaxError("unexpected token in EXPRESSION: %s", current.Token)
	}

	if err != nil {
		return nil, err
	}

	current = p.current()
	if current != nil && current.Token.IsOperator() {
		return p.parseInfixExpression(left, precedence)
	}

	return left, nil
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

func (p *LLParser) parseNextInfixExpression(expr ast.Expression, precedence int) (ast.Expression, error) {
	currentOperator := p.current()
	if currentOperator != nil && currentOperator.Token.IsOperator() {
		return p.parseInfixExpression(expr, precedence)
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

	switch operator.Token {
	case token.Colon:
		return p.parseCallExpression(left, true)

	case token.LParen:
		return p.parseCallExpression(left, false)

	case token.LBracket:
		return p.parseIndexExpressionBracketNotaion(left)

	case token.Period:
		return p.parseIndexExpressionPeriodNotation(left)
	}

	right, err := p.parseExpression(operatorPrecedence)
	if err != nil {
		return nil, err
	}

	expr := &ast.InfixExpression{
		LeftOperand:  left,
		Operator:     operator.Token,
		RightOperand: right,
	}

	return p.parseNextInfixExpression(expr, precedence)
}

// call-expression
// => expression [ ":" identifier ] "(" [expression-list] ")"
func (p *LLParser) parseCallExpression(callable ast.Expression, findMethod bool) (ast.Expression, error) {
	var err error
	var member *ast.Identifier
	if findMethod {
		member, err = p.parseIdentifier()
		if err != nil {
			return nil, err
		}

		if err := p.skipToken(token.LParen); err != nil {
			return nil, err
		}
	}

	arguments, err := p.parseExpressionList()
	if err != nil {
		return nil, err
	}

	if err := p.skipToken(token.RParen); err != nil {
		return nil, err
	}

	expr := &ast.CallExpression{
		Callable: callable,
		Member:   member,
		Args:     arguments,
	}

	return p.parseNextInfixExpression(expr, PrecedenceCall)
}

// index-expression
// => expression "[" expression "]"
func (p *LLParser) parseIndexExpressionBracketNotaion(base ast.Expression) (ast.Expression, error) {
	var err error
	var index ast.Expression

	index, err = p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	if err := p.skipToken(token.RBracket); err != nil {
		return nil, err
	}

	expr := &ast.IndexExpression{
		Base:     base,
		Operator: token.LBracket,
		Index:    index,
	}

	return p.parseNextInfixExpression(expr, PrecedenceLowest)
}

// index-expression
// => expression "." identifier
func (p *LLParser) parseIndexExpressionPeriodNotation(base ast.Expression) (ast.Expression, error) {
	var err error
	var index ast.Expression

	index, err = p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	expr := &ast.IndexExpression{
		Base:     base,
		Operator: token.Period,
		Index:    index,
	}

	return p.parseNextInfixExpression(expr, PrecedenceLowest)
}

func (p *LLParser) parseGroupExpression() (ast.Expression, error) {
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

	return expr, nil
}
