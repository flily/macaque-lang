package compiler

import (
	"github.com/flily/macaque-lang/ast"
	"github.com/flily/macaque-lang/errors"
	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/opcode"
	"github.com/flily/macaque-lang/token"
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
		c.Context.Variable.EnterScope(FrameScopeFunction)
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
		i := c.Context.Literal.ReferenceString(n.Value)
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

	case *ast.HashLiteral:
		l := len(n.Pairs)
		for _, pair := range n.Pairs {
			if e = r.Append(c.compileCode(pair.Key, FlagNone)); e != nil {
				break CompileSwitch
			}

			if e = r.Append(c.compileCode(pair.Value, FlagNone)); e != nil {
				break CompileSwitch
			}
		}

		r.Write(opcode.IMakeHash, l)
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

	case *ast.IndexExpression:
		if e = r.Append(c.compileIndexExpression(n)); e != nil {
			break CompileSwitch
		}

	case *ast.CallExpression:
		if e = r.Append(c.compileCallExpression(n, flag)); e != nil {
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

	c.Context.Variable.EnterScope(FrameScopeBlock)
	consequence, err := c.compileCode(n.Consequence, FlagNone)
	if err != nil {
		return nil, err
	}
	c.Context.Variable.LeaveScope()

	if n.Alternative == nil {
		code.Write(opcode.IJumpIf, consequence.Code.Length())
		code.AppendCode(consequence)
		return code, nil
	}

	c.Context.Variable.EnterScope(FrameScopeBlock)
	alternative, err := c.compileCode(n.Alternative, FlagNone)
	if err != nil {
		return nil, err
	}
	c.Context.Variable.LeaveScope()

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
	c.Context.Variable.EnterScope(FrameScopeFunction)

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

func (c *Compiler) compileCallExpression(expr *ast.CallExpression, flag uint64) (*CompileResult, error) {
	result := NewCompileResult()

	args := NewCompileResult()
	l := expr.Args.Length()
	for i := 0; i < l; i++ {
		a := expr.Args.Expressions[l-i-1]
		if e := args.Append(c.compileCode(a, FlagNone|FlagPackValue)); e != nil {
			return nil, e
		}
	}
	result.AppendCode(args)

	callable, err := c.compileCode(expr.Callable, flag)
	if err != nil {
		return nil, err
	}
	result.AppendCode(callable)

	result.Write(opcode.ICall, args.Values)
	return result, nil
}

func (c *Compiler) compileIndexExpression(expr *ast.IndexExpression) (*CompileResult, error) {
	result := NewCompileResult()

	base, err := c.compileCode(expr.Base, FlagNone|FlagPackValue)
	if err != nil {
		return nil, err
	}

	var index *CompileResult

	if expr.Operator == token.LBracket {
		index, err = c.compileCode(expr.Index, FlagNone|FlagPackValue)
		if err != nil {
			return nil, err
		}
	} else {
		key := expr.Index.(*ast.Identifier)
		i := c.Context.Literal.ReferenceString(key.Value)
		index = NewCompileResult()
		index.Write(opcode.ILoad, int(i))
		index.Values = 1
	}

	result.AppendCode(base)
	result.AppendCode(index)
	result.Write(opcode.IIndex)

	return result, nil
}
