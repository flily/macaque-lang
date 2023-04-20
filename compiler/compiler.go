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
	code, err := c.compileCode(p)
	if err != nil {
		return 0, err
	}

	c.Context.Code.Append(code)
	return len(c.Context.Code.Code), nil
}

func (c *Compiler) GetMain() *object.FunctionObject {
	main := &object.FunctionObject{
		StackSize: c.Context.Variable.CurrentFrameSize(),
		IP:        0,
		Bounds:    nil,
	}

	return main
}

func (c *Compiler) compileCode(node ast.Node) (*CodeBuffer, error) {
	b := NewCodeBuffer()
	var e error

CompileSwitch:
	switch n := node.(type) {
	case *ast.Program:
		c.Context.Variable.EnterScope(VariableScopeFunction)
		for _, stmt := range n.Statements {
			if e = b.AppendCode(c.compileCode(stmt)); e != nil {
				break CompileSwitch
			}
		}

	case *ast.ExpressionStatement:
		e = b.AppendCode(c.compileCode(n.Expressions))

	case *ast.ExpressionList:
		l := len(n.Expressions)
		for _, expr := range n.Expressions {
			if e = b.AppendCode(c.compileCode(expr)); e != nil {
				break CompileSwitch
			}
		}

		b.Write(opcode.ISetAX, l)

	case *ast.Identifier:
		ref, kind := c.Context.Variable.Reference(n.Value)
		switch kind {
		case VariableKindGlobal, VariableKindModule:
			b.Write(opcode.ILoad, ref.Offset)

		case VariableKindBinding:
			b.Write(opcode.ILoad, ref.Offset)

		case VariableKindLocal:
			b.Write(opcode.ISLoad, ref.Offset)

		default:
			ctx := n.Position.MakeContext()
			e = errors.NewSyntaxError(ctx, "variable %s undefined", n.Value)
			break CompileSwitch
		}

	case *ast.IntegerLiteral:
		b.Write(opcode.ILoadInt, int(n.Value))

	case *ast.StringLiteral:
		s := n.Value
		i, ok := c.Context.Literal.Lookup(s)
		if !ok {
			i = c.Context.Literal.Add(s, object.NewString(s))
		}

		b.Write(opcode.ILoad, int(i))

	case *ast.InfixExpression:
		if e = b.AppendCode(c.compileCode(n.LeftOperand)); e != nil {
			break CompileSwitch
		}

		if e = b.AppendCode(c.compileCode(n.RightOperand)); e != nil {
			break CompileSwitch
		}

		b.Write(opcode.IBinOp, int(n.Operator))

	case *ast.LetStatement:
		index := make([]int, n.Identifiers.Length())
		for i, v := range n.Identifiers.Identifiers {
			j, ok := c.Context.Variable.DefineVariable(v.Value, v.Position)
			if !ok {
				ctx1, _ := c.Context.Variable.Reference(v.Value)
				ctx2 := v.Position.MakeContext()

				e = ctx2.NewCompilationError("variable %s redeclared", v.Value).
					WithInfo(ctx1.Position.MakeContext(),
						"variable %s is already declared here", v.Value)

				break CompileSwitch
			}
			index[i] = j
		}

		if e = b.AppendCode(c.compileCode(n.Expressions)); e != nil {
			break CompileSwitch
		}

		for k := len(index) - 1; k >= 0; k-- {
			b.Write(opcode.ISStore, index[k])
		}
	}

	return b, e
}
