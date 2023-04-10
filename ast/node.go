package ast

type Node interface {
	Children() []Node
	FormalCode() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}
