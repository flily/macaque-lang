package compiler

import (
	"github.com/flily/macaque-lang/ast"
	"github.com/flily/macaque-lang/errors"
	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/opcode"
)

type Compiler struct {
	Context *CompilerContext
}

func NewCompiler() *Compiler {
	c := &Compiler{
		Context: NewCompilerContext(),
	}

	return c
}

func (c *Compiler) Compile(p *ast.Program) (int, error) {
	return c.compileCode(p)
}

func (c *Compiler) writeCode(name int, ops ...int) {
	inst := opcode.Inst(name, ops...)
	c.Context.Instructions = append(c.Context.Instructions, inst)
}

func (c *Compiler) compileCode(node ast.Node) (int, error) {
	var r int
	var e error

CompileSwitch:
	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			r, e = c.compileCode(stmt)
			if e != nil {
				break CompileSwitch
			}
		}

	case *ast.ExpressionStatement:
		r, e = c.compileCode(n.Expressions)

	case *ast.ExpressionList:
		l := len(n.Expressions)
		for _, expr := range n.Expressions {
			r, e = c.compileCode(expr)
			if e != nil {
				break CompileSwitch
			}
		}

		c.writeCode(opcode.ISetAX, l)

	case *ast.Identifier:
		ref, kind := c.Context.Variable.Reference(n.Value)
		switch kind {
		case VariableKindGlobal, VariableKindModule:
			c.writeCode(opcode.ILoad, ref.Offset)

		case VariableKindBinding:
			c.writeCode(opcode.ILoad, ref.Offset)

		case VariableKindLocal:
			c.writeCode(opcode.ISLoad, ref.Offset)

		default:
			ctx := n.Position.MakeContext()
			e = errors.NewSyntaxError(ctx, "variable %s undefined", n.Value)
			break CompileSwitch
		}

	case *ast.IntegerLiteral:
		c.writeCode(opcode.ILoadInt, int(n.Value))

	case *ast.StringLiteral:
		s := n.Value
		i, ok := c.Context.Literal.Lookup(s)
		if !ok {
			i = c.Context.Literal.Add(s, object.NewString(s))
		}

		c.writeCode(opcode.ILoadStr, int(i))

	case *ast.InfixExpression:
		r, e = c.compileCode(n.LeftOperand)
		if e != nil {
			break CompileSwitch
		}

		r, e = c.compileCode(n.RightOperand)
		if e != nil {
			break CompileSwitch
		}

		c.writeCode(opcode.IBinOp, int(n.Operator))

	case *ast.LetStatement:
		index := make([]int, n.Identifiers.Length())
		for i, v := range n.Identifiers.Identifiers {
			j, ok := c.Context.Variable.DefineVariable(v.Value, v.Position)
			if !ok {
				ctx := v.Position.MakeContext()
				e = errors.NewSyntaxError(ctx, "variable %s redefined", v.Value)
				break CompileSwitch
			}
			index[i] = j
		}

		r, e = c.compileCode(n.Expressions)
		if e != nil {
			break CompileSwitch
		}

		for k := len(index) - 1; k >= 0; k-- {
			c.writeCode(opcode.ISStore, index[k])
		}
	}

	return r, e
}
