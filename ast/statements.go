package ast

import (
	"fmt"
	"strings"
)

type LineStatementNode interface {
	Statement
	lineStatementNode()
}

type BlockStatementNode interface {
	Statement
	blockStatementNode()
}

type LetStatement struct {
	Identifiers *IdentifierList
	Expressions *ExpressionList
}

func (s *LetStatement) statementNode() {}

func (s *LetStatement) Children() []Node {
	idLength := s.Identifiers.Length()
	count := idLength + s.Expressions.Length()
	nodes := make([]Node, count)

	for i, id := range s.Identifiers.Identifiers {
		nodes[i] = id
	}

	for i, expr := range s.Expressions.Expressions {
		nodes[idLength+i] = expr
	}

	return nodes
}

func (s *LetStatement) CanonicalCode() string {
	result := fmt.Sprintf("let %s = %s;",
		s.Identifiers.CanonicalCode(),
		s.Expressions.CanonicalCode(),
	)

	return result
}

func (s *LetStatement) lineStatementNode() {}

type ReturnStatement struct {
	Expressions *ExpressionList
}

func (s *ReturnStatement) statementNode() {}

func (s *ReturnStatement) Children() []Node {
	nodes := make([]Node, s.Expressions.Length())

	for i, expr := range s.Expressions.Expressions {
		nodes[i] = expr
	}

	return nodes
}

func (s *ReturnStatement) CanonicalCode() string {
	result := fmt.Sprintf("return %s;", s.Expressions.CanonicalCode())
	return result
}

func (s *ReturnStatement) lineStatementNode() {}

type IfStatement struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative BlockStatementNode
}

func (s *IfStatement) statementNode() {}

func (s *IfStatement) Children() []Node {
	nodes := []Node{
		s.Condition,
		s.Consequence,
		s.Alternative,
	}

	return nodes
}

func (s *IfStatement) blockStatementNode() {}

type BlockStatement struct {
	Statements []Statement
}

func (s *BlockStatement) statementNode() {}

func (s *BlockStatement) Children() []Node {
	nodes := make([]Node, len(s.Statements))

	for i, stmt := range s.Statements {
		nodes[i] = stmt
	}

	return nodes
}

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

func (s *BlockStatement) blockStatementNode() {}

type ExpressionStatement struct {
	Expressions *ExpressionList
}

func (s *ExpressionStatement) statementNode() {}

func (s *ExpressionStatement) Children() []Node {
	return s.Expressions.Children()
}

func (s *ExpressionStatement) CanonicalCode() string {
	return s.Expressions.CanonicalCode()
}

func (s *ExpressionStatement) lineStatementNode() {}

// ImportStatement is not determined yet.
type ImportStatement struct {
}

func (s *ImportStatement) statementNode() {}

func (s *ImportStatement) Children() []Node {
	return nil
}

func (s *ImportStatement) CanonicalCode() string {
	return "import;"
}
