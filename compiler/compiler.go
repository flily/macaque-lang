package compiler

import (
	"github.com/flily/macaque-lang/ast"
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
		for _, expr := range n.Expressions {
			r, e = c.compileCode(expr)
			if e != nil {
				break CompileSwitch
			}
		}

	case *ast.IntegerLiteral:
		c.writeCode(opcode.ILoadInt, opcode.XLoadIntLiteral, int(n.Value))

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
	}

	return r, e
}
