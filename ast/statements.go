package ast

import (
	"fmt"
	"strings"

	"github.com/flily/macaque-lang/token"
)

type StatementBase struct {
	LeadingComments *token.Context
}

func (s *StatementBase) SetLeadingComments(c *token.Context) {
	s.LeadingComments = c
}

func (s *StatementBase) GetLeadingComments() *token.Context {
	return s.LeadingComments
}

type LetStatement struct {
	StatementBase

	Let         *token.TokenContext
	Identifiers *IdentifierList
	Assign      *token.TokenContext
	Expressions *ExpressionList
	Semicolon   *token.TokenContext
}

func (s *LetStatement) statementNode()     {}
func (s *LetStatement) lineStatementNode() {}

func (s *LetStatement) CanonicalCode() string {
	result := fmt.Sprintf("let %s = %s;",
		s.Identifiers.CanonicalCode(),
		s.Expressions.CanonicalCode(),
	)

	return result
}

func (s *LetStatement) GetContext() *token.Context {
	c := token.JoinContext(
		s.Let.ToContext(),
		s.Identifiers.GetContext(),
		s.Assign.ToContext(),
		s.Expressions.GetContext(),
		s.Semicolon.ToContext(),
	)

	return c
}

func (s *LetStatement) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *LetStatement:
		result = s.Identifiers.EqualTo(n.Identifiers) &&
			s.Expressions.EqualTo(n.Expressions)
	}

	return result
}

type ReturnStatement struct {
	StatementBase

	Return      *token.TokenContext
	Expressions *ExpressionList
	Semicolon   *token.TokenContext
}

func (s *ReturnStatement) statementNode()     {}
func (s *ReturnStatement) lineStatementNode() {}

func (s *ReturnStatement) CanonicalCode() string {
	result := fmt.Sprintf("return %s;", s.Expressions.CanonicalCode())
	return result
}

func (s *ReturnStatement) GetContext() *token.Context {
	c := token.JoinContext(
		s.Return.ToContext(),
		s.Expressions.GetContext(),
		s.Semicolon.ToContext(),
	)

	return c
}

func (s *ReturnStatement) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *ReturnStatement:
		result = s.Expressions.EqualTo(n.Expressions)
	}

	return result
}

type IfStatement struct {
	StatementBase

	Expression *IfExpression
}

func (s *IfStatement) statementNode()      {}
func (s *IfStatement) blockStatementNode() {}

func (s *IfStatement) CanonicalCode() string {
	return s.Expression.CanonicalCode()
}

func (s *IfStatement) GetContext() *token.Context {
	return s.Expression.GetContext()
}

func (s *IfStatement) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *IfStatement:
		result = s.Expression.EqualTo(n.Expression)
	}

	return result
}

type BlockStatement struct {
	StatementBase

	LBrace     *token.TokenContext
	Statements []Statement
	RBrace     *token.TokenContext
}

func (s *BlockStatement) statementNode()      {}
func (s *BlockStatement) blockStatementNode() {}

func (s *BlockStatement) CanonicalCode() string {
	length := len(s.Statements)
	result := make([]string, length+2)

	result[0] = "{"
	for i, stmt := range s.Statements {
		result[i+1] = stmt.CanonicalCode()
	}
	result[length-1] = "}"

	return strings.Join(result, "\n")
}

func (s *BlockStatement) GetContext() *token.Context {
	ctxList := make([]*token.Context, len(s.Statements)+2)
	ctxList[0] = s.LBrace.ToContext()
	for i, stmt := range s.Statements {
		ctxList[i+1] = stmt.GetContext()
	}
	ctxList[len(ctxList)-1] = s.RBrace.ToContext()

	return token.JoinContext(ctxList...)
}

func (s *BlockStatement) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *BlockStatement:
		result = len(s.Statements) == len(n.Statements)
		for i, stmt := range s.Statements {
			result = result && stmt.EqualTo(n.Statements[i])
		}
	}

	return result
}

func (s *BlockStatement) AddStatement(stmt Statement) {
	s.Statements = append(s.Statements, stmt)
}

type ExpressionStatement struct {
	StatementBase

	Expressions *ExpressionList
	Semicolon   *token.TokenContext
}

func (s *ExpressionStatement) statementNode()     {}
func (s *ExpressionStatement) lineStatementNode() {}

func (s *ExpressionStatement) CanonicalCode() string {
	return fmt.Sprintf("%s;", s.Expressions.CanonicalCode())
}

func (s *ExpressionStatement) GetContext() *token.Context {
	c := token.JoinContext(
		s.Expressions.GetContext(),
		s.Semicolon.ToContext(),
	)

	return c
}

func (s *ExpressionStatement) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *ExpressionStatement:
		result = s.Expressions.EqualTo(n.Expressions)
	}

	return result
}

// ImportStatement is not determined yet.
type ImportStatement struct {
	StatementBase

	Import *token.TokenContext
	Target *token.TokenContext
}

func (s *ImportStatement) statementNode()     {}
func (s *ImportStatement) lineStatementNode() {}

func (s *ImportStatement) CanonicalCode() string {
	return "import;"
}

func (s *ImportStatement) GetContext() *token.Context {
	c := token.NewContext(s.Import, s.Target)

	return c
}

func (s *ImportStatement) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *ImportStatement:
		if s.Target != nil && n.Target != nil {
			result = s.Target.Content == n.Target.Content
		}
	}

	return result
}
