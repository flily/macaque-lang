package compiler

import (
	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/opcode"
	"github.com/flily/macaque-lang/token"
)

const (
	VariableScopeGlobal   = 1
	VariableScopeModule   = 2
	VariableScopeFunction = 3
	VariableScopeBlock    = 4

	VariableKindNotFound = 0
	VariableKindGlobal   = 1
	VariableKindModule   = 2
	VariableKindLocal    = 3
	VariableKindBinding  = 4
)

type VariableInfo struct {
	Name     string
	Offset   int
	Position *token.TokenInfo
}

type VariableScopeContext struct {
	Scope     int
	Variables map[string]VariableInfo
	outer     *VariableScopeContext
	arguments int
	variables int
	FrameSize int
}

func NewVariableScopeContext() *VariableScopeContext {
	c := &VariableScopeContext{
		Scope:     VariableKindGlobal,
		Variables: make(map[string]VariableInfo),
	}

	return c
}

func (c *VariableScopeContext) IsRoot() bool {
	return c.outer == nil
}

func (c *VariableScopeContext) IsDefined(name string) bool {
	_, ok := c.Variables[name]
	return ok
}

func (c *VariableScopeContext) DefineArgument(name string, pos *token.TokenInfo) (int, bool) {
	if c.IsDefined(name) {
		return 0, false
	}

	n := c.arguments + 1
	c.arguments = -n
	c.Variables[name] = VariableInfo{
		Name:     name,
		Offset:   -n,
		Position: pos,
	}

	return n, true
}

func (c *VariableScopeContext) DefineVariable(name string, pos *token.TokenInfo) (int, bool) {
	if c.IsDefined(name) {
		return 0, false
	}

	n := c.variables + 1
	c.variables = n
	c.Variables[name] = VariableInfo{
		Name:     name,
		Offset:   n,
		Position: pos,
	}

	return n, true
}

func (c *VariableScopeContext) currentVariableKind() int {
	switch c.Scope {
	case VariableScopeGlobal:
		return VariableKindGlobal

	case VariableScopeModule:
		return VariableKindModule

	case VariableScopeFunction, VariableScopeBlock:
		return VariableKindLocal
	}

	return 0
}

func (c *VariableScopeContext) Reference(name string) (VariableInfo, int) {
	v, ok := c.Variables[name]
	if ok {
		return v, c.currentVariableKind()
	}

	if c.IsRoot() {
		return VariableInfo{}, VariableKindNotFound
	}

	info, kind := c.outer.Reference(name)
	if kind != VariableKindLocal {
		return info, kind
	} else {
		if c.Scope == VariableScopeFunction {
			return info, VariableKindBinding
		} else {
			return info, VariableKindLocal
		}
	}
}

func (c *VariableScopeContext) EnterScope(scope int) *VariableScopeContext {
	s := &VariableScopeContext{
		Scope:     scope,
		Variables: make(map[string]VariableInfo),
		outer:     c,
	}

	return s
}

func (c *VariableScopeContext) UpdateFrameSize(count int) int {
	n := c.variables + count
	r := 0
	switch c.Scope {
	case VariableScopeBlock:
		r = c.outer.UpdateFrameSize(n)

	case VariableScopeFunction:
		r = n
		if n > c.FrameSize {
			c.FrameSize = n
		}
	}

	return r
}

func (c *VariableScopeContext) LeaveScope() *VariableScopeContext {
	c.UpdateFrameSize(0)
	return c.outer
}

type VariableContext struct {
	root *VariableScopeContext
	top  *VariableScopeContext
}

func NewVariableContext() *VariableContext {
	root := NewVariableScopeContext()
	c := &VariableContext{
		root: root,
		top:  root,
	}

	c.EnterScope(VariableScopeModule) // For current file module
	// c.EnterScope(VariableScopeFunction) // For main function
	return c
}

func (c *VariableContext) EnterScope(scope int) {
	c.top = c.top.EnterScope(scope)
}

func (c *VariableContext) LeaveScope() {
	c.top = c.top.LeaveScope()
}

func (c *VariableContext) DefineArgument(name string, pos *token.TokenInfo) (int, bool) {
	return c.top.DefineArgument(name, pos)
}

func (c *VariableContext) DefineVariable(name string, pos *token.TokenInfo) (int, bool) {
	return c.top.DefineVariable(name, pos)
}

func (c *VariableContext) Reference(name string) (VariableInfo, int) {
	return c.top.Reference(name)
}

func (c *VariableContext) CurrentFrameSize() int {
	return c.top.UpdateFrameSize(0)
}

type LiteralContext struct {
	index  map[interface{}]uint64
	Values []object.Object
	counts int
}

func NewLiteralContext() *LiteralContext {
	c := &LiteralContext{
		index:  make(map[interface{}]uint64),
		Values: make([]object.Object, 0),
	}

	return c
}

func (c *LiteralContext) Add(v interface{}, o object.Object) uint64 {
	n := c.counts
	c.index[c] = uint64(n)
	c.Values = append(c.Values, o)
	c.counts = n + 1
	return uint64(n)
}

func (c *LiteralContext) Lookup(literal interface{}) (uint64, bool) {
	n, ok := c.index[literal]
	if ok {
		return uint64(n), true
	}

	return 0, false
}

type CodeBuffer struct {
	Code []opcode.Opcode
}

func (b *CodeBuffer) Write(code int, operands ...int) {
	c := opcode.Code(code, operands...)
	b.Code = append(b.Code, c)
}

func (b *CodeBuffer) Append(buf *CodeBuffer) {
	b.Code = append(b.Code, buf.Code...)
}

func (b *CodeBuffer) AppendCode(buf *CodeBuffer, err error) error {
	if err != nil {
		return err
	}

	b.Append(buf)
	return nil
}

func NewCodeBuffer() *CodeBuffer {
	b := &CodeBuffer{
		Code: make([]opcode.Opcode, 0),
	}

	return b
}

type CompilerContext struct {
	Variable *VariableContext
	Literal  *LiteralContext
	Code     *CodeBuffer
}

func NewCompilerContext() *CompilerContext {
	c := &CompilerContext{
		Variable: NewVariableContext(),
		Literal:  NewLiteralContext(),
		Code:     NewCodeBuffer(),
	}

	return c
}
