package ast

import (
	"fmt"
	"strings"
)

type LetStatement struct {
	Identifiers *IdentifierList
	Expressions *ExpressionList
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
	Expressions *ExpressionList
}

func (s *ReturnStatement) statementNode()     {}
func (s *ReturnStatement) lineStatementNode() {}

func (s *ReturnStatement) CanonicalCode() string {
	result := fmt.Sprintf("return %s;", s.Expressions.CanonicalCode())
	return result
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
	Expression *IfExpression
}

func (s *IfStatement) statementNode()      {}
func (s *IfStatement) blockStatementNode() {}

func (s *IfStatement) CanonicalCode() string {
	return s.Expression.CanonicalCode()
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
	Statements []Statement
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
	Expressions *ExpressionList
}

func (s *ExpressionStatement) statementNode()     {}
func (s *ExpressionStatement) lineStatementNode() {}

func (s *ExpressionStatement) CanonicalCode() string {
	return fmt.Sprintf("%s;", s.Expressions.CanonicalCode())
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
}

func (s *ImportStatement) statementNode()     {}
func (s *ImportStatement) lineStatementNode() {}

func (s *ImportStatement) CanonicalCode() string {
	return "import;"
}

func (s *ImportStatement) EqualTo(node Node) bool {
	result := false
	switch node.(type) {
	case *ImportStatement:
		result = true
	}

	return result
}
