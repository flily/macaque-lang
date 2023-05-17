package compiler

import (
	"github.com/flily/macaque-lang/ast"
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

	page := c.Context.LinkCodePage(r.Code)
	return page, nil
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
		if n.Length() > 1 {
			f |= FlagPackValue
		}

		for _, expr := range n.Expressions.Expressions {
			if e = r.Append(c.compileCode(expr.Expression, uint64(f))); e != nil {
				break CompileSwitch
			}
		}

		r.Write(opcode.IMakeList, n.Length())
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
			ctx := n.GetContext()
			e = NewSemanticError(ctx, "variable %s undefined", n.Value)
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
			if e = r.Append(c.compileCode(expr.Expression, f)); e != nil {
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

		r.Write(opcode.IBinOp, int(n.Operator.Token))
		r.Values = 1

	case *ast.PrefixExpression:
		if e = r.Append(c.compileCode(n.Operand, flag|FlagPackValue)); e != nil {
			break CompileSwitch
		}

		r.Write(opcode.IUniOp, int(n.Prefix.Token))
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
		for i, item := range n.Identifiers.Identifiers {
			v := item.Identifier
			j, ok := c.Context.Variable.DefineVariable(v.Value, v.Context)
			if !ok {
				ctx1, _ := c.Context.Variable.Reference(v.Value)
				ctx2 := v.Context.Tokens[0]

				e = NewSemanticError(ctx2.ToContext(), "variable %s redeclared", v.Value).
					WithInfo(ctx1.Context,
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

	case *ast.ReturnStatement:
		if e = r.Append(c.compileCode(n.Expressions, FlagNone)); e != nil {
			break CompileSwitch
		}

		r.Write(opcode.IReturn)
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

	var alternative *CompileResult
	c.Context.Variable.EnterScope(FrameScopeBlock)
	if n.Alternative != nil {
		alternative, err = c.compileCode(n.Alternative, FlagNone)
		if err != nil {
			return nil, err
		}

	} else {
		alternative = NewCompileResult()
		alternative.Write(opcode.ILoadNull)
		alternative.Values = 1
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
			// r.Write(opcode.IPop, count)
			r.Write(opcode.IClean)
		}

		lastResult, err := c.compileCode(last, FlagNone)
		if err != nil {
			return nil, err
		}

		r.AppendCode(lastResult)
		r.Values = lastResult.Values
	}

	if withReturn {
		r.Write(opcode.IReturn)
	}

	return r, e
}

func (c *Compiler) compileFunctionLiteral(f *ast.FunctionLiteral) (*CompileResult, error) {
	result := NewCompileResult()
	c.Context.Variable.EnterScope(FrameScopeFunction)

	for _, item := range f.Arguments.Identifiers {
		c.Context.Variable.DefineArgument(item.Identifier.Value, item.Identifier.GetContext())
	}

	r, e := c.compileStatements(f.Body.Statements, true)
	if e != nil {
		return nil, e
	}

	frameSize := c.Context.Variable.CurrentScope().UpdateFrameSize(0)
	functionContext := &FunctionContext{
		FunctionInfo: FunctionInfo{
			Arguments: f.Arguments.Length(),
			FrameSize: frameSize,
		},
		Code: r.Code,
	}

	id := c.Context.AddFunction(functionContext)
	scope := c.Context.Variable.CurrentScope()
	c.Context.Variable.LeaveScope()

	for _, arg := range scope.BindingOrder {
		c.compileIdentifierReference(arg.Name, result)
	}

	result.Write(opcode.IMakeFunc, id, len(scope.Bindings))
	result.Values = 1
	return result, nil
}

func (c *Compiler) compileCallExpression(expr *ast.CallExpression, flag uint64) (*CompileResult, error) {
	result := NewCompileResult()

	args := NewCompileResult()
	l := expr.Args.Length()
	for i := 0; i < l; i++ {
		a := expr.Args.Expressions[l-i-1]
		if e := args.Append(c.compileCode(a.Expression, FlagNone|FlagPackValue)); e != nil {
			return nil, e
		}
	}
	result.AppendCode(args)

	switch expr.Token.GetToken() {
	case token.Nil:
		callable, err := c.compileCode(expr.Base, flag)
		if err != nil {
			return nil, err
		}
		result.AppendCode(callable)

	case token.DualColon:
		callable, err := c.compileCode(expr.Base, flag)
		if err != nil {
			return nil, err
		}

		memberIndex := c.Context.Literal.ReferenceString(expr.Member.Value)
		callable.Write(opcode.ILoad, int(memberIndex))
		callable.Write(opcode.ISDUP)
		callable.Write(opcode.IIndex)
		result.AppendCode(callable)

	case token.Fn:
		result.Write(opcode.ISLoad, 0)
	}

	result.Write(opcode.ICall, args.Values)
	result.Values = 1
	return result, nil
}

func (c *Compiler) compileIndexExpression(expr *ast.IndexExpression) (*CompileResult, error) {
	result := NewCompileResult()

	base, err := c.compileCode(expr.Base, FlagNone|FlagPackValue)
	if err != nil {
		return nil, err
	}

	var index *CompileResult

	if expr.Operator.Token == token.LBracket {
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
