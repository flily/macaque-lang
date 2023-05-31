package ast

import (
	"fmt"
	"strings"

	"github.com/flily/macaque-lang/token"
)

type Identifier struct {
	Value   string
	Context *token.Context
}

func NewIdentifier(value string, ctx *token.TokenContext) *Identifier {
	i := &Identifier{
		Value:   value,
		Context: ctx.ToContext(),
	}

	return i
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) literalValue()   {}

func (i *Identifier) CanonicalCode() string {
	return i.Value
}

func (i *Identifier) GetContext() *token.Context {
	if i == nil {
		return nil
	}

	return i.Context
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
	Value   int64
	Content string
	Context *token.Context
}

func (i *IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) literalValue()   {}

func (i *IntegerLiteral) CanonicalCode() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *IntegerLiteral) GetContext() *token.Context {
	return i.Context
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
	Value   float64
	Content string
	Context *token.Context
}

func (f *FloatLiteral) expressionNode() {}
func (f *FloatLiteral) literalValue()   {}

func (f *FloatLiteral) CanonicalCode() string {
	return fmt.Sprintf("%f", f.Value)
}

func (f *FloatLiteral) GetContext() *token.Context {
	return f.Context
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
	Value   string
	Content string
	Context *token.Context
}

func (s *StringLiteral) expressionNode() {}
func (s *StringLiteral) literalValue()   {}

func (s *StringLiteral) CanonicalCode() string {
	return s.Content
}

func (s *StringLiteral) GetContext() *token.Context {
	return s.Context
}

func (s *StringLiteral) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *StringLiteral:
		result = n.Content == s.Content
	}

	return result
}

type BooleanLiteral struct {
	Value   bool
	Context *token.Context
}

func (b *BooleanLiteral) expressionNode() {}
func (b *BooleanLiteral) literalValue()   {}

func (b *BooleanLiteral) CanonicalCode() string {
	if b.Value {
		return token.STrue
	}

	return token.SFalse
}

func (b *BooleanLiteral) GetContext() *token.Context {
	return b.Context
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
	Context *token.Context
}

func (n *NullLiteral) expressionNode() {}
func (n *NullLiteral) literalValue()   {}

func (n *NullLiteral) Walk(v Visitor) int {
	return v.Visit(n)
}

func (n *NullLiteral) CanonicalCode() string {
	return token.SNull
}

func (n *NullLiteral) GetContext() *token.Context {
	return n.Context
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
	LBracket    *token.TokenContext
	Expressions *ExpressionList
	RBracket    *token.TokenContext
}

func (a *ArrayLiteral) expressionNode() {}
func (a *ArrayLiteral) literalValue()   {}

func (a *ArrayLiteral) CanonicalCode() string {
	code := a.Expressions.CanonicalCode()
	return fmt.Sprintf("[%s]", code)
}

func (a *ArrayLiteral) GetContext() *token.Context {
	c := token.JoinContext(
		a.LBracket.ToContext(),
		a.Expressions.GetContext(),
		a.RBracket.ToContext(),
	)

	return c
}

func (a *ArrayLiteral) EqualTo(node Node) bool {
	result := false
	switch n := node.(type) {
	case *ArrayLiteral:
		return a.Expressions.EqualTo(n.Expressions)
	}

	return result
}

func (a *ArrayLiteral) Length() int {
	return a.Expressions.Length()
}

type HashItem struct {
	Key   Expression
	Colon *token.TokenContext
	Value Expression
	Comma *token.TokenContext
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

func (h *HashItem) GetContext() *token.Context {
	c := token.JoinContext(
		h.Key.GetContext(),
		h.Colon.ToContext(),
		h.Value.GetContext(),
		h.Comma.ToContext(),
	)

	return c
}

type HashLiteral struct {
	LBrace *token.TokenContext
	Pairs  []*HashItem
	RBrace *token.TokenContext
}

func (h *HashLiteral) expressionNode() {}
func (h *HashLiteral) literalValue()   {}

func (h *HashLiteral) CanonicalCode() string {
	pairs := make([]string, len(h.Pairs))
	for i, p := range h.Pairs {
		pairs[i] = fmt.Sprintf("%s: %s", p.Key.CanonicalCode(), p.Value.CanonicalCode())
	}

	return fmt.Sprintf("{%s}", strings.Join(pairs, ", \n"))
}

func (h *HashLiteral) GetContext() *token.Context {
	ctxs := make([]*token.Context, len(h.Pairs)+2)
	ctxs[0] = h.LBrace.ToContext()
	for i, pair := range h.Pairs {
		ctxs[i+1] = pair.GetContext()
	}
	ctxs[len(ctxs)-1] = h.RBrace.ToContext()

	return token.JoinContext(ctxs...)
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

func (h *HashLiteral) AddPair(key Expression, colon *token.TokenContext, value Expression, comma *token.TokenContext) {
	item := &HashItem{
		Key:   key,
		Colon: colon,
		Value: value,
		Comma: comma,
	}
	h.Pairs = append(h.Pairs, item)
}

type FunctionLiteral struct {
	Function     *token.TokenContext
	LParen       *token.TokenContext
	Arguments    *IdentifierList
	RParen       *token.TokenContext
	Body         *BlockStatement
	ReturnValues int
}

func (f *FunctionLiteral) expressionNode() {}

func (f *FunctionLiteral) CanonicalCode() string {
	result := fmt.Sprintf("fn(%s) %s",
		f.Arguments.CanonicalCode(),
		f.Body.CanonicalCode(),
	)

	return result
}

func (f *FunctionLiteral) GetContext() *token.Context {
	c := token.JoinContext(
		f.Function.ToContext(),
		f.LParen.ToContext(),
		f.Arguments.GetContext(),
		f.RParen.ToContext(),
		f.Body.GetContext(),
	)

	return c
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
