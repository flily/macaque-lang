package ast

import (
	"github.com/flily/macaque-lang/token"
)

func GetContext(node Node) *token.Context {
	if node == nil {
		return nil
	}

	return node.GetContext()
}

func ContextGroup(nodes ...Node) *token.Context {
	ctxList := make([]*token.Context, len(nodes))

	for i, n := range nodes {
		ctxList[i] = GetContext(n)
	}

	return token.JoinContext(ctxList...)
}
