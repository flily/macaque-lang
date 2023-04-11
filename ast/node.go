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
