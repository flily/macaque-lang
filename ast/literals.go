package ast

import (
	"fmt"
	"strings"

	"github.com/flily/macaque-lang/token"
)

type Identifier struct {
	Value    string
	Position token.TokenInfo
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) Children() []Node {
	return nil
}

func (i *Identifier) CanonicalCode() string {
	return i.Value
}

type IntegerLiteral struct {
	Value    int64
	Content  string
	Position token.TokenInfo
}

func (i *IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) Children() []Node {
	return nil
}

func (i *IntegerLiteral) CanonicalCode() string {
	return fmt.Sprintf("%d", i.Value)
}

type FloatLiteral struct {
	Value    float64
	Content  string
	Position token.TokenInfo
}

func (f *FloatLiteral) expressionNode() {}

func (f *FloatLiteral) Children() []Node {
	return nil
}

func (f *FloatLiteral) CanonicalCode() string {
	return fmt.Sprintf("%f", f.Value)
}

type StringLiteral struct {
	Value    string
	Content  string
	Position token.TokenInfo
}

func (s *StringLiteral) expressionNode() {}

func (s *StringLiteral) Children() []Node {
	return nil
}

func (s *StringLiteral) CanonicalCode() string {
	return s.Content
}

type BooleanLiteral struct {
	Value    bool
	Position token.TokenInfo
}

func (b *BooleanLiteral) expressionNode() {}

func (b *BooleanLiteral) Children() []Node {
	return nil
}

func (b *BooleanLiteral) CanonicalCode() string {
	if b.Value {
		return token.STrue
	}

	return token.SFalse
}

type NullLiteral struct {
	Position token.TokenInfo
}

func (n *NullLiteral) expressionNode() {}

func (n *NullLiteral) Children() []Node {
	return nil
}

func (n *NullLiteral) CanonicalCode() string {
	return token.SNull
}

type ArrayLiteral struct {
	Elements []Expression
	Position token.TokenInfo
}

func (a *ArrayLiteral) expressionNode() {}

func (a *ArrayLiteral) Children() []Node {
	var nodes []Node

	for _, e := range a.Elements {
		nodes = append(nodes, e)
	}

	return nodes
}

func (a *ArrayLiteral) CanonicalCode() string {
	elements := make([]string, len(a.Elements))
	for i, e := range a.Elements {
		elements[i] = e.CanonicalCode()
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

type HashItem struct {
	Key   Expression
	Value Expression
}

type HashLiteral struct {
	Pairs    []HashItem
	Position token.TokenInfo
}

func (h *HashLiteral) expressionNode() {}

func (h *HashLiteral) Children() []Node {
	nodes := make([]Node, len(h.Pairs)*2)

	for i, p := range h.Pairs {
		nodes[i*2+0] = p.Key
		nodes[i*2+1] = p.Value
	}

	return nodes
}

func (h *HashLiteral) CanonicalCode() string {
	pairs := make([]string, len(h.Pairs))
	for i, p := range h.Pairs {
		pairs[i] = fmt.Sprintf("%s: %s", p.Key.CanonicalCode(), p.Value.CanonicalCode())
	}

	return fmt.Sprintf("{%s}", strings.Join(pairs, ", \n"))
}

type FunctionLiteral struct {
	Arguments *IdentifierList
	Body      *BlockStatement
}

func (f *FunctionLiteral) expressionNode() {}

func (f *FunctionLiteral) Children() []Node {
	nodes := make([]Node, f.Arguments.Length()+1)

	for i, arg := range f.Arguments.Identifiers {
		nodes[i] = arg
	}

	nodes[f.Arguments.Length()] = f.Body
	return nodes
}

func (f *FunctionLiteral) CanonicalCode() string {
	result := fmt.Sprintf("fn(%s) %s",
		f.Arguments.CanonicalCode(),
		f.Body.CanonicalCode(),
	)

	return result
}
