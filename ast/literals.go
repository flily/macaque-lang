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

func (i *Identifier) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *Identifier:
		result = n.Value == i.Value
	}

	return result
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

func (i *IntegerLiteral) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *IntegerLiteral:
		result = n.Content == i.Content
	}

	return result
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

func (f *FloatLiteral) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *FloatLiteral:
		result = n.Content == f.Content
	}

	return result
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

func (s *StringLiteral) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *FloatLiteral:
		result = n.Content == s.Content
	}

	return result
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

func (b *BooleanLiteral) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *BooleanLiteral:
		result = n.Value == b.Value
	}

	return result
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

func (n *NullLiteral) EqualTo(node Node) bool {
	result := false
	switch node.(type) {
	case *NullLiteral:
		result = true
	}

	return result
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

func (a *ArrayLiteral) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *ArrayLiteral:
		if len(n.Elements) == len(a.Elements) {
			result = true
			for i, e := range a.Elements {
				if !e.EqualTo(n.Elements[i]) {
					result = false
					break
				}
			}
		}
	}

	return result
}

type HashItem struct {
	Key   Expression
	Value Expression
}

func (h *HashItem) EqualTo(item *HashItem) bool {
	if ok := h.Key.EqualTo(item.Key); !ok {
		return false
	}

	if ok := h.Value.EqualTo(item.Value); !ok {
		return false
	}

	return true
}

type HashLiteral struct {
	Pairs    []*HashItem
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

func (h *HashLiteral) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *HashLiteral:
		if len(n.Pairs) == len(h.Pairs) {
			result = true
			for i, p := range h.Pairs {
				if !p.EqualTo(n.Pairs[i]) {
					result = false
					break
				}
			}
		}
	}

	return result
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

func (f *FunctionLiteral) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *FunctionLiteral:
		if ok := f.Arguments.EqualTo(n.Arguments); !ok {
			break
		}

		if ok := f.Body.EqualTo(n.Body); !ok {
			break
		}

		result = true
	}

	return result
}
