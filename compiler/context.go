package compiler

import (
	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/opcode"
	"github.com/flily/macaque-lang/token"
)

const (
	FrameScopeGlobal    FrameScope = 1
	FrameScopeModule    FrameScope = 2
	FrameScopeFunction  FrameScope = 3
	FrameScopeBlock     FrameScope = 4
	VariableKindMiss    VarKind    = 0
	VariableKindGlobal  VarKind    = 1
	VariableKindModule  VarKind    = 2
	VariableKindLocal   VarKind    = 3
	VariableKindBinding VarKind    = 4
)

type FrameScope int

var scopeNames = [...]string{
	FrameScopeGlobal:   "global",
	FrameScopeModule:   "module",
	FrameScopeFunction: "function",
	FrameScopeBlock:    "block",
}

func (s FrameScope) String() string {
	if s >= 0 && s <= FrameScopeBlock {
		return scopeNames[s]
	}

	return "unknown"
}

type VarKind int

var kindNames = [...]string{
	VariableKindMiss:    "miss",
	VariableKindGlobal:  "global",
	VariableKindModule:  "module",
	VariableKindLocal:   "local",
	VariableKindBinding: "binding",
}

func (k VarKind) String() string {
	if k >= 0 && k <= VariableKindBinding {
		return kindNames[k]
	}

	return "unknown"
}

type VariableInfo struct {
	Name     string
	Offset   int
	Kind     VarKind
	Position *token.TokenInfo
}

type VariableScopeContext struct {
	Level     int
	Scope     FrameScope
	Variables map[string]VariableInfo
	Bindings  map[string]VariableInfo
	outer     *VariableScopeContext
	arguments int
	variables int
	FrameSize int
}

func NewVariableScopeContext() *VariableScopeContext {
	c := &VariableScopeContext{
		Level:     0,
		Scope:     FrameScopeGlobal,
		Variables: make(map[string]VariableInfo),
		Bindings:  make(map[string]VariableInfo),
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

func (c *VariableScopeContext) FrameOffset() int {
	size := c.variables
	switch c.Scope {
	case FrameScopeFunction:
		return size

	case FrameScopeBlock:
		return size + c.outer.FrameOffset()

	default:
		return 0
	}
}

func (c *VariableScopeContext) DefineArgument(name string, pos *token.TokenInfo) (int, bool) {
	if c.IsDefined(name) {
		return 0, false
	}

	n := c.arguments + 1
	c.arguments = n
	c.Variables[name] = VariableInfo{
		Name:     name,
		Offset:   -n,
		Kind:     VariableKindLocal,
		Position: pos,
	}

	return n, true
}

func (c *VariableScopeContext) DefineVariable(name string, pos *token.TokenInfo) (int, bool) {
	if c.IsDefined(name) {
		return 0, false
	}

	n := c.FrameOffset() + 1
	c.variables += 1
	c.Variables[name] = VariableInfo{
		Name:     name,
		Offset:   n,
		Kind:     VariableKindLocal,
		Position: pos,
	}

	return n, true
}

func (c *VariableScopeContext) AddBinding(name string, info VariableInfo) VariableInfo {
	if v, ok := c.Bindings[name]; ok {
		return v
	}

	offset := len(c.Bindings)
	info.Offset = offset
	c.Bindings[name] = info
	return info
}

func (c *VariableScopeContext) currentVariableKind() VarKind {
	switch c.Scope {
	case FrameScopeGlobal:
		return VariableKindGlobal

	case FrameScopeModule:
		return VariableKindModule

	case FrameScopeFunction, FrameScopeBlock:
		return VariableKindLocal
	}

	return 0
}

func (c *VariableScopeContext) Reference(name string) (VariableInfo, VarKind) {
	v, ok := c.Variables[name]
	if ok {
		return v, c.currentVariableKind()
	}

	v, ok = c.Bindings[name]
	if ok {
		return v, VariableKindBinding
	}

	if c.IsRoot() {
		return VariableInfo{}, VariableKindMiss
	}

	info, kind := c.outer.Reference(name)
	if kind != VariableKindLocal {
		if kind == VariableKindBinding {
			c.AddBinding(name, info)
		}

		return info, kind
	} else {
		if c.Scope == FrameScopeFunction {
			binding := c.AddBinding(name, info)
			return binding, VariableKindBinding
		} else {
			return info, VariableKindLocal
		}
	}
}

func (c *VariableScopeContext) EnterScope(scope FrameScope) *VariableScopeContext {
	s := &VariableScopeContext{
		Level:     c.Level + 1,
		Scope:     scope,
		Variables: make(map[string]VariableInfo),
		Bindings:  make(map[string]VariableInfo),
		outer:     c,
	}

	return s
}

func (c *VariableScopeContext) UpdateFrameSize(count int) int {
	n := c.variables + count
	r := 0
	switch c.Scope {
	case FrameScopeBlock:
		r = c.outer.UpdateFrameSize(n)

	case FrameScopeFunction:
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

	c.EnterScope(FrameScopeModule) // For current file module
	// c.EnterScope(VariableScopeFunction) // For main function
	return c
}

func (c *VariableContext) CurrentScope() *VariableScopeContext {
	return c.top
}

func (c *VariableContext) EnterScope(scope FrameScope) {
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

func (c *VariableContext) Reference(name string) (VariableInfo, VarKind) {
	info, kind := c.top.Reference(name)
	return info, kind
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
	c.index[v] = uint64(n)
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

func (c *LiteralContext) ReferenceString(s string) uint64 {
	if n, ok := c.Lookup(s); ok {
		return n
	}

	o := object.NewString(s)
	return c.Add(s, o)
}

type FunctionContext struct {
	FunctionInfo
	Code *CodeBuffer
}

type CompilerContext struct {
	Variable  *VariableContext
	Literal   *LiteralContext
	Functions []*FunctionContext
}

func (c *CompilerContext) LinkCodePage(main *CodeBuffer) *CodePage {
	code := NewCodeBuffer()
	links := make([]FunctionInfo, len(c.Functions))

	mainInfo := FunctionInfo{
		IP:        0,
		FrameSize: c.Variable.CurrentFrameSize(),
	}

	code.Append(main)
	links[0] = mainInfo

	if len(c.Functions) > 1 {
		code.Write(opcode.IHalt)
		for i, f := range c.Functions {
			if i == 0 {
				continue
			}

			offset := code.Length()

			info := f.FunctionInfo
			info.IP = uint64(offset)

			code.Append(f.Code)
			code.Write(opcode.IHalt)
			links[i] = info
		}
	}

	data := make([]object.Object, len(c.Literal.Values))
	copy(data, c.Literal.Values)

	page := &CodePage{
		Codes:     code.Code,
		Functions: links,
		Data:      data,
	}

	return page
}

func (c *CompilerContext) AddFunction(f *FunctionContext) int {
	n := len(c.Functions)
	f.Index = n
	c.Functions = append(c.Functions, f)
	return n
}

func NewCompilerContext() *CompilerContext {
	c := &CompilerContext{
		Variable: NewVariableContext(),
		Literal:  NewLiteralContext(),
		Functions: []*FunctionContext{
			nil, // reserve for main function
		},
	}

	return c
}
