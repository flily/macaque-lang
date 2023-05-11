package parser

import (
	"github.com/flily/macaque-lang/ast"
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
	token.LBracket:   true,
	token.LBrace:     true,
	token.If:         true,
	token.Fn:         true,
	token.LastToken:  false,
}

func isExpressionFirstSet(token token.Token) bool {
	return expressionFirstSet[token]
}

const (
	RuleLetStatement        = "let statement"
	RuleExpressionStatement = "expression statement"
	RuleReturnStatement     = "return statement"
	RuleBlockStatement      = "block statement"
	RuleFunctionLiteral     = "function literal"
	RuleIdentifierList      = "identifier list"
	RuleIdentifier          = "identifier"
	RuleExpressionList      = "expression list"
	RuleExpression          = "expression"
	RuleArrayLiteral        = "array literal"
	RuleHashLiteral         = "hash literal"
	RuleCallExpression      = "call expression"
	RuleIndexExpression     = "index expression"
	RuleGroupedExpression   = "grouped expression"
	RuleIfExpression        = "if expression"
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

func (p *LLParser) expect(token token.Token, rule string) error {
	current := p.container.Current()
	if current.Token != token {
		return p.makeSyntaxError("expect token %s IN %s, but got %s",
			token, rule, current.Token)
	}

	return nil
}

func (p *LLParser) current() *token.TokenContext {
	return p.container.Current()
}

func (p *LLParser) currentToken() token.Token {
	current := p.container.Current()
	return current.Token
}

func (p *LLParser) DebugLine() string {
	current := p.container.Current()
	return current.Position.MakeLineHighlight()
}

// func (p *LLParser) peek(offset int) *lex.LexicalElement {
// 	return p.container.Peek(offset)
// }

func (p *LLParser) nextToken() *token.TokenContext {
	return p.container.Next()
}

func (p *LLParser) skipToken(token token.Token, rule string) (*token.TokenContext, error) {
	if err := p.expect(token, rule); err != nil {
		return nil, err
	}

	elem := p.container.current()
	p.container.Next()
	return elem, nil
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

func (p *LLParser) makeSyntaxError(format string, args ...interface{}) *token.SyntaxNError {
	current := p.container.Current()
	return current.NewSyntaxError(format, args...)
}

func (p *LLParser) parseProgram() (*ast.Program, error) {
	program := ast.NewEmptyProgram()
	current := p.current()

	p.DebugLine()

	for current != nil && current.Token != token.EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		program.AddStatement(stmt)

		next := p.current()
		if next == current {
			return nil, p.makeSyntaxError("parser does not shift any token")
		}

		current = next
	}

	return program, nil
}

// Parse statements

func (p *LLParser) parseStatement() (ast.Statement, error) {
	var stmt ast.Statement
	var err error

	current := p.current()
	switch current.Token {
	case token.Let:
		stmt, err = p.parseLetStatement()

	case token.Comment:
		// skip comment
		p.skipComment()

	case token.Null, token.False, token.True, token.Integer, token.Float, token.String,
		token.Identifier, token.Minus, token.Bang, token.LParen, token.LBracket, token.LBrace,
		token.If, token.Fn:
		stmt, err = p.parseExpressionStatement()

	case token.Return:
		stmt, err = p.parseReturnStatement()

	default:
		return nil, p.makeSyntaxError("unexpected token in PROGRAM: %s", current.Token)
	}

	if err != nil {
		return nil, err
	}

	return stmt, nil
}

// let-stmt => "let" identifier-list "=" expression-list ";"
func (p *LLParser) parseLetStatement() (*ast.LetStatement, error) {
	var sLet, sAssign, sSemicolon *token.TokenContext
	sLet, _ = p.skipToken(token.Let, RuleLetStatement)

	idList, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}

	if sAssign, err = p.skipToken(token.Assign, RuleLetStatement); err != nil {
		return nil, err
	}

	exprList, err := p.parseExpressionList(ExprListMustHave)
	if err != nil {
		return nil, err
	}

	if sSemicolon, err = p.skipToken(token.Semicolon, RuleLetStatement); err != nil {
		return nil, err
	}

	stmt := &ast.LetStatement{
		Let:         sLet,
		Identifiers: idList,
		Assign:      sAssign,
		Expressions: exprList,
		Semicolon:   sSemicolon,
	}

	return stmt, nil
}

// expression-stmt => expression-list ";"
func (p *LLParser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	var sSemicolon *token.TokenContext
	exprList, err := p.parseExpressionList(ExprListMustHave)
	if err != nil {
		return nil, err
	}

	sSemicolon, _ = p.skipToken(token.Semicolon, RuleExpressionStatement)

	stmt := &ast.ExpressionStatement{
		Expressions: exprList,
		Semicolon:   sSemicolon,
	}

	return stmt, nil
}

// return-stmt => "return" [ expression-list ] ";"
func (p *LLParser) parseReturnStatement() (*ast.ReturnStatement, error) {
	var sReturn, sSemicolon *token.TokenContext
	sReturn, _ = p.skipToken(token.Return, RuleReturnStatement)

	exprList, err := p.parseExpressionList(ExprListCanBeEmpty)
	if err != nil {
		return nil, err
	}

	sSemicolon, _ = p.skipToken(token.Semicolon, RuleReturnStatement)

	stmt := &ast.ReturnStatement{
		Return:      sReturn,
		Expressions: exprList,
		Semicolon:   sSemicolon,
	}

	return stmt, nil
}

// block-stmt => "{" *statement "}"
func (p *LLParser) parseBlockStatement(context string) (*ast.BlockStatement, error) {
	var sLBrace, sRBrace *token.TokenContext
	var err error
	sLBrace, _ = p.skipToken(token.LBrace, context)

	block := &ast.BlockStatement{
		LBrace: sLBrace,
	}
	current := p.current()
	for current != nil && current.Token != token.RBrace && current.Token != token.EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		block.AddStatement(stmt)

		next := p.current()
		if next == current {
			return nil, p.makeSyntaxError("parser does not shift any token")
		}

		current = next
	}

	if sRBrace, err = p.skipToken(token.RBrace, context); err != nil {
		return nil, err
	}

	block.RBrace = sRBrace
	return block, nil
}

func (p *LLParser) parseIfStatement() (*ast.IfStatement, error) {
	expr, err := p.parseIfExpression()
	if err != nil {
		return nil, err
	}

	stmt := &ast.IfStatement{
		Expression: expr,
	}

	return stmt, nil
}

// Parse terminal symbols, include identifiers and literals.

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
			id = ast.NewIdentifier(current.Content, current)
			p.nextToken()
			found = true

		case token.Comment:
			p.skipComment()

		default:
			found = true
			err = p.expect(token.Identifier, RuleIdentifier)
		}
	}

	return id, err
}

// identifier-list => identifier *( "," identifier ) [","]
func (p *LLParser) parseIdentifierList() (*ast.IdentifierList, error) {
	list := &ast.IdentifierList{}

	for {
		id, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}

		comma, err := p.skipToken(token.Comma, RuleIdentifierList)
		list.AddIdentifier(id, comma)
		if err != nil {
			break
		}

		if p.expect(token.Identifier, RuleIdentifierList) != nil {
			break
		}
	}

	return list, nil
}

// literals
// => null-literal
// => boolean-literal
// => integer-literal
// => float-literal
// => string-literal
// => array-literal
// => hash-literal
// => function-literal
func (p *LLParser) parseLiteral() (ast.Expression, error) {
	var expr ast.Expression
	var err error

	current := p.current()
	switch current.Token {
	case token.Integer:
		expr = newInteger(current)
		p.nextToken()

	case token.Float:
		expr = newFloat(current)
		p.nextToken()

	case token.String:
		expr = newString(current)
		p.nextToken()

	case token.True, token.False:
		expr = newBoolean(current)
		p.nextToken()

	case token.Null:
		expr = newNull(current)
		p.nextToken()

	case token.LBracket:
		expr, err = p.parseArrayLiteral()

	case token.LBrace:
		expr, err = p.parseHashLiteral()

	case token.Fn:
		expr, err = p.parseFunctionLiteral()
	}

	return expr, err
}

// array-literal
// => "[" expression-list "]"
func (p *LLParser) parseArrayLiteral() (*ast.ArrayLiteral, error) {
	var sLBracket, sRBracket *token.TokenContext
	var err error
	sLBracket, _ = p.skipToken(token.LBracket, RuleArrayLiteral)

	if sRBracket, err = p.skipToken(token.RBracket, RuleArrayLiteral); err == nil {
		return array(), nil
	}

	list, err := p.parseExpressionList(ExprListCanBeEmpty)
	if err != nil {
		return nil, err
	}

	if sRBracket, err = p.skipToken(token.RBracket, RuleArrayLiteral); err != nil {
		return nil, err
	}

	array := &ast.ArrayLiteral{
		LBracket:    sLBracket,
		Expressions: list,
		RBracket:    sRBracket,
	}

	return array, nil
}

// hash-literal => "{" hash-pair *( "," hash-pair ) [","] "}"
// hash-pair => expression ":" expression
func (p *LLParser) parseHashLiteral() (*ast.HashLiteral, error) {
	var sLBrace, sRBrace *token.TokenContext
	var err error

	sLBrace, _ = p.skipToken(token.LBrace, RuleHashLiteral)

	hash := &ast.HashLiteral{
		LBrace: sLBrace,
	}

	for {
		current := p.current()
		if current.Token == token.RBrace {
			break
		}

		var colon, comma *token.TokenContext
		key, err := p.parseExpression(PrecedenceLowest)
		if err != nil {
			return nil, err
		}

		if colon, err = p.skipToken(token.Colon, RuleHashLiteral); err != nil {
			return nil, err
		}

		value, err := p.parseExpression(PrecedenceLowest)
		if err != nil {
			return nil, err
		}

		comma, err = p.skipToken(token.Comma, RuleHashLiteral)
		hash.AddPair(key, colon, value, comma)

		if err != nil {
			break
		}
	}

	if sRBrace, err = p.skipToken(token.RBrace, RuleHashLiteral); err != nil {
		return nil, err
	}

	hash.RBrace = sRBrace
	return hash, nil
}

func (p *LLParser) parseFunctionLiteral() (ast.Expression, error) {
	var sFunction, sLParen, sRParen *token.TokenContext
	var err error
	sFunction, _ = p.skipToken(token.Fn, RuleFunctionLiteral)

	if sLParen, err = p.skipToken(token.LParen, RuleFunctionLiteral); err != nil {
		return nil, err
	}

	args, err := p.parseExpressionList(ExprListCanBeEmpty)
	if err != nil {
		return nil, err
	}

	if sRParen, err = p.skipToken(token.RParen, RuleFunctionLiteral); err != nil {
		return nil, err
	}

	current := p.current()

	if p.currentToken() != token.LBrace {
		callExpr := &ast.CallExpression{
			Base:      nil,
			Token:     sFunction,
			LParen:    sLParen,
			Args:      args,
			RParen:    sRParen,
			Recursion: true,
		}

		return callExpr, nil
	}

	if !args.IsIdentifierList() {
		err := current.Position.MakeContext().NewSyntaxError(
			"recursion function call MUST NOT follow by a block statement",
		)
		return nil, err
	}

	if err := p.expect(token.LBrace, RuleFunctionLiteral); err != nil {
		return nil, err
	}

	body, err := p.parseBlockStatement(RuleFunctionLiteral)
	if err != nil {
		return nil, err
	}

	literal := &ast.FunctionLiteral{
		Function:     sFunction,
		LParen:       sLParen,
		Arguments:    args.ToIdentifierList(),
		RParen:       sRParen,
		Body:         body,
		ReturnValues: -1,
	}

	return literal, nil
}

// Parse expressions

// expression-list => expression *( "," expression ) [","]
func (p *LLParser) parseExpressionList(canBeEmpty bool) (*ast.ExpressionList, error) {
	var err error
	list := &ast.ExpressionList{}

	current := p.current()
	if !canBeEmpty && !isExpressionFirstSet(current.Token) {
		return nil, p.expect(token.Identifier, RuleExpressionList)
	}

	for isExpressionFirstSet(current.Token) {
		exp, err := p.parseExpression(PrecedenceLowest)
		if err != nil {
			return nil, err
		}

		comma, err := p.skipToken(token.Comma, RuleExpressionList)
		list.AddExpression(exp, comma)
		if err != nil {
			break
		}

		current = p.current()
	}

	return list, err
}

func (p *LLParser) parseExpression(precedence int) (ast.Expression, error) {
	var expr ast.Expression
	var err error

	current := p.current()

	switch current.Token {
	case token.LParen:
		expr, err = p.parseGroupExpression()

	case token.Integer, token.Float, token.String, token.True, token.False, token.Null,
		token.LBracket, token.LBrace, token.Fn:
		expr, err = p.parseLiteral()

	case token.Identifier:
		expr, err = p.parseIdentifier()

	case token.Bang, token.Minus:
		expr, err = p.parsePrefixExpression()

	case token.If:
		expr, err = p.parseIfExpression()

	default:
		err = p.makeSyntaxError("unexpected token in EXPRESSION: %s", current.Token)
	}

	if err != nil {
		return nil, err
	}

	return p.parseExpressionWithOperator(expr, precedence)
}

func (p *LLParser) parseExpressionWithOperator(expr ast.Expression, precedence int) (ast.Expression, error) {
	current := p.current()
	operator := current.Token
	currentPrecedent := GetPrecedence(operator)

	if currentPrecedent <= precedence {
		return expr, nil
	}

	var err error
	switch operator {
	case token.LParen:
		expr, err = p.parseCallExpression(expr, false)

	case token.DualColon:
		expr, err = p.parseCallExpression(expr, true)

	case token.LBracket:
		expr, err = p.parseIndexExpressionBracketNotaion(expr)

	case token.Period:
		expr, err = p.parseIndexExpressionPeriodNotation(expr)

	default:
		if IsInfixOperator(current.Token) {
			expr, err = p.parseInfixExpression(expr, precedence)

		} else {
			err = p.makeSyntaxError("unexpected token in EXPRESSION_OP: %s", operator)
		}
	}

	if err != nil {
		return nil, err
	}

	return p.parseExpressionWithOperator(expr, precedence)
}

// prefix-expression => prefix-operator expression
// prefix-operator
// => "!" / "-"    ; original design
// => "~" 		   ; extended operators
func (p *LLParser) parsePrefixExpression() (*ast.PrefixExpression, error) {
	operator := p.current()
	p.nextToken()

	operand, err := p.parseExpression(PrecedencePrefix)
	if err != nil {
		return nil, err
	}

	expr := &ast.PrefixExpression{
		Prefix:  operator,
		Operand: operand,
	}

	return expr, nil

}

// infix-expression => expression infix-operator expression
// infix-operator
// => "+" / "-" / "*" / "/" / "==" / "!=" / "<" / ">"     ; original design
// => "%" / "<=" / ">=" / "&&" / "||" / "&" / "|" / "^"   ; extended operators
func (p *LLParser) parseInfixExpression(left ast.Expression, precedence int) (ast.Expression, error) {
	operator := p.current()
	currentPrecedence := GetPrecedence(operator.Token)
	p.nextToken()

	right, err := p.parseExpression(currentPrecedence)
	if err != nil {
		return nil, err
	}

	expr := &ast.InfixExpression{
		LeftOperand:  left,
		Operator:     operator,
		RightOperand: right,
	}

	return p.parseExpressionWithOperator(expr, precedence)
}

// call-expression => expression [ "::" identifier ] "(" [expression-list] ")"
func (p *LLParser) parseCallExpression(callable ast.Expression, findMethod bool) (ast.Expression, error) {
	var sLParen, sRParen, sToken *token.TokenContext
	var err error
	var member *ast.Identifier

	if findMethod {
		dualColon, _ := p.skipToken(token.DualColon, RuleCallExpression)
		sToken = dualColon
		member, err = p.parseIdentifier()
		if err != nil {
			return nil, err
		}

		if sLParen, err = p.skipToken(token.LParen, RuleCallExpression); err != nil {
			return nil, err
		}

	} else {
		sLParen, _ = p.skipToken(token.LParen, RuleCallExpression)
	}

	arguments, err := p.parseExpressionList(ExprListCanBeEmpty)
	if err != nil {
		return nil, err
	}

	if sRParen, err = p.skipToken(token.RParen, RuleCallExpression); err != nil {
		return nil, err
	}

	expr := &ast.CallExpression{
		Base:   callable,
		Token:  sToken,
		Member: member,
		LParen: sLParen,
		Args:   arguments,
		RParen: sRParen,
	}

	return p.parseExpressionWithOperator(expr, PrecedenceCall)
}

// index-expression => expression "[" expression "]"
func (p *LLParser) parseIndexExpressionBracketNotaion(base ast.Expression) (ast.Expression, error) {
	var err error
	var index ast.Expression

	lb, _ := p.skipToken(token.LBracket, RuleIndexExpression)

	index, err = p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	var rb *token.TokenContext
	if rb, err = p.skipToken(token.RBracket, RuleIndexExpression); err != nil {
		return nil, err
	}

	expr := &ast.IndexExpression{
		Base:     base,
		Operator: lb,
		Index:    index,
		End:      rb,
	}

	return p.parseExpressionWithOperator(expr, PrecedenceIndex)
}

// index-expression => expression "." identifier
func (p *LLParser) parseIndexExpressionPeriodNotation(base ast.Expression) (ast.Expression, error) {
	var err error
	var index ast.Expression

	period, _ := p.skipToken(token.Period, RuleIndexExpression)

	index, err = p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	expr := &ast.IndexExpression{
		Base:     base,
		Operator: period,
		Index:    index,
	}

	return p.parseExpressionWithOperator(expr, PrecedenceIndex)
}

// group-expression => "(" expression ")"
func (p *LLParser) parseGroupExpression() (ast.Expression, error) {
	if _, err := p.skipToken(token.LParen, RuleGroupedExpression); err != nil {
		return nil, err
	}

	expr, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	if _, err := p.skipToken(token.RParen, RuleGroupedExpression); err != nil {
		return nil, err
	}

	return expr, nil
}

// if-stmt => if-expression
// if-expression
// => "if" "(" expression ")" block-statement
// => "if" "(" expression ")" block-statement "else" block-statement
// => "if" "(" expression ")" block-statement "else" if-stmt
func (p *LLParser) parseIfExpression() (*ast.IfExpression, error) {
	var sIf, sLParen, sRParen *token.TokenContext
	var err error

	sIf, _ = p.skipToken(token.If, RuleIfExpression)
	if sLParen, err = p.skipToken(token.LParen, RuleIfExpression); err != nil {
		return nil, err
	}

	condition, err := p.parseExpression(PrecedenceLowest)
	if err != nil {
		return nil, err
	}

	if sRParen, err = p.skipToken(token.RParen, RuleIfExpression); err != nil {
		return nil, err
	}

	consequence, err := p.parseBlockStatement(RuleIfExpression)
	if err != nil {
		return nil, err
	}

	stmt := &ast.IfExpression{
		If:          sIf,
		LParen:      sLParen,
		Condition:   condition,
		RParen:      sRParen,
		Consequence: consequence,
	}

	if sElse, err := p.skipToken(token.Else, RuleIfExpression); err == nil {
		stmt.Else = sElse
		var alternative ast.BlockStatementNode
		var err error

		current := p.current()
		switch current.Token {
		case token.If:
			alternative, err = p.parseIfStatement()

		case token.LBrace:
			alternative, err = p.parseBlockStatement(RuleIfExpression)

		default:
			return nil, p.makeSyntaxError("unexpected token in IF-ELSE: %s", current.Token)
		}

		if err != nil {
			return nil, err
		}

		stmt.Alternative = alternative
	}

	return stmt, nil
}
