package ast

type Node interface {
	Children() []Node
	CanonicalCode() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}
