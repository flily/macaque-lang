package ast

type Node interface {
	Children() []Node
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
