package ast

import (
	"strings"
)

type Node interface {
	CanonicalCode() string
	EqualTo(Node) bool
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type LiteralValue interface {
	Expression
	literalValue()
}

type LineStatementNode interface {
	Statement
	lineStatementNode()
}

type BlockStatementNode interface {
	Statement
	blockStatementNode()
}

type Program struct {
	Statements []Statement
}

func NewEmptyProgram() *Program {
	p := &Program{
		Statements: make([]Statement, 0),
	}

	return p
}

func (p *Program) AddStatement(stmt Statement) {
	p.Statements = append(p.Statements, stmt)
}

func (p *Program) EqualTo(node Node) bool {
	result := false

TypeSwitch:
	switch n := node.(type) {
	case *Program:
		if len(p.Statements) != len(n.Statements) {
			break
		}

		for i, stmt := range p.Statements {
			if !stmt.EqualTo(n.Statements[i]) {
				break TypeSwitch
			}
		}

		result = true
	}

	return result
}

func (p *Program) CanonicalCode() string {
	lines := make([]string, len(p.Statements))
	for i, stmt := range p.Statements {
		lines[i] = stmt.CanonicalCode()
	}

	return strings.Join(lines, "\n")
}
