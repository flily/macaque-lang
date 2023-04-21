package compiler

import (
	"github.com/flily/macaque-lang/ast"
	"github.com/flily/macaque-lang/errors"
	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/opcode"
)

const (
	FlagNone      = 0x0
	FlagPackValue = 0x1
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

func (c *Compiler) Compile(p *ast.Program) (*CodePage, error) {
	r, err := c.compileCode(p, FlagNone)
	if err != nil {
		return nil, err
	}

	page := c.Context.LinkCode(r.Code)
	return page, nil
}

func (c *Compiler) GetMain() *object.FunctionObject {
	main := &object.FunctionObject{
		StackSize: c.Context.Variable.CurrentFrameSize(),
		IP:        0,
		Bounds:    nil,
	}

	return main
}

func (c *Compiler) compileCode(node ast.Node, flag uint64) (*CompileResult, error) {
	r := NewCompileResult()
	var e error

CompileSwitch:
	switch n := node.(type) {
	case *ast.Program:
		c.Context.Variable.EnterScope(VariableScopeFunction)
		if e = r.Append(c.compileStatements(n.Statements, false)); e != nil {
			break CompileSwitch
		}

	case *ast.NullLiteral:
		r.Write(opcode.ILoadNull)
		r.Values = 1

	case *ast.BooleanLiteral:
		if n.Value {
			r.Write(opcode.ILoadBool, 1)
		} else {
			r.Write(opcode.ILoadBool, 0)
		}
		r.Values = 1

	case *ast.IntegerLiteral:
		r.Write(opcode.ILoadInt, int(n.Value))
		r.Values = 1

	case *ast.StringLiteral:
		s := n.Value
		i, ok := c.Context.Literal.Lookup(s)
		if !ok {
			i = c.Context.Literal.Add(s, object.NewString(s))
		}

		r.Write(opcode.ILoad, int(i))
		r.Values = 1

	case *ast.ArrayLiteral:
		f := flag
		if len(n.Elements) > 1 {
			f |= FlagPackValue
		}

		for _, expr := range n.Elements {
			if e = r.Append(c.compileCode(expr, uint64(f))); e != nil {
				break CompileSwitch
			}
		}

		r.Write(opcode.IMakeList, len(n.Elements))
		r.Values = 1

	case *ast.Identifier:
		p := c.compileIdentifierReference(n.Value, r)
		if p <= 0 {
			ctx := n.Position.MakeContext()
			e = errors.NewSyntaxError(ctx, "variable %s undefined", n.Value)
			break CompileSwitch
		}

		r.Values = 1

	case *ast.FunctionLiteral:
		if e = r.Append(c.compileFunctionLiteral(n)); e != nil {
			break CompileSwitch
		}

	case *ast.ExpressionList:
		isList := n.Length() > 1
		f := flag
		if isList {
			f |= FlagPackValue
		}
		for _, expr := range n.Expressions {
			if e = r.Append(c.compileCode(expr, f)); e != nil {
				break CompileSwitch
			}
		}

		if flag&FlagPackValue != 0 {
			r.Write(opcode.IMakeList, n.Length())
			r.Values = 1
		}

	case *ast.InfixExpression:
		if e = r.Append(c.compileCode(n.LeftOperand, flag|FlagPackValue)); e != nil {
			break CompileSwitch
		}

		if e = r.Append(c.compileCode(n.RightOperand, flag|FlagPackValue)); e != nil {
			break CompileSwitch
		}

		r.Write(opcode.IBinOp, int(n.Operator))
		r.Values = 1

	case *ast.PrefixExpression:
		if e = r.Append(c.compileCode(n.Operand, flag|FlagPackValue)); e != nil {
			break CompileSwitch
		}

		r.Write(opcode.IUniOp, int(n.PrefixOperator))
		r.Values = 1

	case *ast.IfExpression:
		if e = r.Append(c.compileIfExpression(n)); e != nil {
			break CompileSwitch
		}

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

		if e = r.Append(c.compileCode(n.Expressions, flag)); e != nil {
			break CompileSwitch
		}

		for k := len(index) - 1; k >= 0; k-- {
			r.Write(opcode.ISStore, index[k])
		}
		r.Values = 0

	case *ast.IfStatement:
		if e = r.Append(c.compileIfExpression(n.Expression)); e != nil {
			break CompileSwitch
		}

	case *ast.ExpressionStatement:
		if e = r.Append(c.compileCode(n.Expressions, flag)); e != nil {
			break CompileSwitch
		}

	case *ast.BlockStatement:
		if e = r.Append(c.compileStatements(n.Statements, false)); e != nil {
			break CompileSwitch
		}
	}

	return r, e
}

func (c *Compiler) compileIdentifierReference(name string, r *CompileResult) int {
	ref, kind := c.Context.Variable.Reference(name)
	n := 1
	switch kind {
	case VariableKindGlobal, VariableKindModule:
		r.Write(opcode.ILoad, ref.Offset)

	case VariableKindBinding:
		r.Write(opcode.ILoadBind, ref.Offset)

	case VariableKindLocal:
		r.Write(opcode.ISLoad, ref.Offset)

	default:
		n = 0
	}

	return n
}

func (c *Compiler) compileIfExpression(n *ast.IfExpression) (*CompileResult, error) {
	code, err := c.compileCode(n.Condition, FlagNone)
	if err != nil {
		return nil, err
	}

	consequence, err := c.compileCode(n.Consequence, FlagNone)
	if err != nil {
		return nil, err
	}

	if n.Alternative == nil {
		code.Write(opcode.IJumpIf, consequence.Code.Length())
		code.AppendCode(consequence)
		return code, nil
	}

	alternative, err := c.compileCode(n.Alternative, FlagNone)
	if err != nil {
		return nil, err
	}

	consequence.Write(opcode.IJumpFWD, alternative.Code.Length())
	code.Write(opcode.IJumpIf, consequence.Code.Length())
	code.AppendCode(consequence)
	code.AppendCode(alternative)

	return code, nil
}

func (c *Compiler) compileStatements(statements []ast.Statement, withReturn bool) (*CompileResult, error) {
	r := NewCompileResult()
	var e error

	var last ast.Statement
	for i, stmt := range statements {
		if i < len(statements)-1 {
			if e = r.Append(c.compileCode(stmt, FlagNone)); e != nil {
				break
			}
		} else {
			last = stmt
		}
	}

	if last == nil {
		r.Write(opcode.ILoadNull)
		r.Values = 1

	} else {
		count := r.Values
		if count > 0 {
			r.Write(opcode.IPop, count)
		}

		lastResult, err := c.compileCode(last, FlagNone)
		if err != nil {
			return nil, err
		}

		r.AppendCode(lastResult)
		r.Values = lastResult.Values
	}

	if withReturn {
		r.Write(opcode.IReturn, r.Values)
	}

	return r, e
}

func (c *Compiler) compileFunctionLiteral(f *ast.FunctionLiteral) (*CompileResult, error) {
	result := NewCompileResult()
	c.Context.Variable.EnterScope(VariableScopeFunction)

	for _, arg := range f.Arguments.Identifiers {
		c.Context.Variable.DefineArgument(arg.Value, arg.Position)
	}

	r, e := c.compileStatements(f.Body.Statements, true)
	if e != nil {
		return nil, e
	}

	functionContext := &FunctionContext{
		Code: r.Code,
	}

	id := c.Context.AddFunction(functionContext)
	scope := c.Context.Variable.CurrentScope()
	c.Context.Variable.LeaveScope()

	for _, arg := range scope.Bindings {
		c.compileIdentifierReference(arg.Name, result)
	}

	result.Write(opcode.IMakeFunc, id, len(scope.Bindings))
	return result, nil
}
