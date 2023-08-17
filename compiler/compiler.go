package compiler

import (
	"github.com/flily/macaque-lang/ast"
	"github.com/flily/macaque-lang/opcode"
	"github.com/flily/macaque-lang/token"
)

type Compiler struct {
	Context *CompilerContext
}

func NewCompiler() *Compiler {
	c := &Compiler{
		Context: NewCompilerContext(),
	}

	c.Context.Variable.EnterScope(FrameScopeFunction)
	return c
}

func (c *Compiler) Compile(node ast.Node) (*opcode.CodePage, error) {
	r, err := c.CompileCode(node)
	if err != nil {
		return nil, err
	}

	page := c.Link(r)
	return page, nil
}

func (c *Compiler) CompileCode(node ast.Node) (*opcode.CodeBlock, error) {
	return c.compileStatement(node, NewFlag(FlagWithReturn))
}

func (c *Compiler) CompileCodeSnippet(node ast.Node) (*opcode.CodeBlock, error) {
	return c.compileStatement(node, NewFlag(FlagNone))
}

func (c *Compiler) Link(block *opcode.CodeBlock) *opcode.CodePage {
	return c.Context.LinkCodePage(block)
}

func (c *Compiler) compileStatement(node ast.Node, flag CompilerFlag) (*opcode.CodeBlock, error) {
	r := opcode.NewCodeBlock()
	var e error
	ctx := node.GetContext()

	switch {
	case flag.Has(FlagCleanStack):
		r.IL(ctx, opcode.IClean)
	}

	nextFlag := flag
	nextFlag.ClearNonPassable()

CompileSwitch:
	switch n := node.(type) {
	case *ast.Program:
		if flag.Has(FlagWithReturn) {
			nextFlag.Set(FlagWithReturn)
		}

		if e = r.Append(c.compileStatements(n.GetContext(), n.Statements, nextFlag)); e != nil {
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

		if e = r.Append(c.compileExpression(n.Expressions, nextFlag)); e != nil {
			break CompileSwitch
		}

		for k := len(index) - 1; k >= 0; k-- {
			r.IL(ctx, opcode.ISStore, index[k])
		}
		r.SetValues(0)

	case *ast.IfStatement:
		if e = r.Append(c.compileIfExpression(n.Expression)); e != nil {
			break CompileSwitch
		}

	case *ast.ExpressionStatement:
		if e = r.Append(c.compileExpression(n.Expressions, nextFlag)); e != nil {
			break CompileSwitch
		}

	case *ast.BlockStatement:
		if e = r.Append(c.compileStatements(n.GetContext(), n.Statements, NewFlag(FlagNone))); e != nil {
			break CompileSwitch
		}

	case *ast.ReturnStatement:
		if e = r.Append(c.compileExpression(n.Expressions, NewFlag(FlagNone))); e != nil {
			break CompileSwitch
		}

		r.IL(ctx, opcode.IReturn)
	}

	return r, e
}

func (c *Compiler) compileExpression(expr ast.Expression, flag CompilerFlag) (*opcode.CodeBlock, error) {
	r := opcode.NewCodeBlock()
	var e error

	ctx := expr.GetContext()
CompileSwitch:
	switch n := expr.(type) {
	case *ast.NullLiteral:
		r.IL(ctx, opcode.ILoadNull).
			SetValues(1)

	case *ast.BooleanLiteral:
		if n.Value {
			r.IL(ctx, opcode.ILoadBool, 1)
		} else {
			r.IL(ctx, opcode.ILoadBool, 0)
		}

		r.SetValues(1)

	case *ast.IntegerLiteral:
		r.IL(ctx, opcode.ILoadInt, int(n.Value)).
			SetValues(1)

	case *ast.StringLiteral:
		i := c.Context.Literal.ReferenceString(n.Value)
		r.IL(ctx, opcode.ILoad, int(i)).
			SetValues(1)

	case *ast.ArrayLiteral:
		f := flag
		if n.Length() > 1 {
			f.Set(FlagPackValue)
		}

		for _, expr := range n.Expressions.Expressions {
			if e = r.Append(c.compileExpression(expr.Expression, f)); e != nil {
				break CompileSwitch
			}
		}

		r.IL(ctx, opcode.IMakeList, n.Length()).
			SetValues(1)

	case *ast.HashLiteral:
		l := len(n.Pairs)
		for _, pair := range n.Pairs {
			if e = r.Append(c.compileExpression(pair.Key, NewFlag(FlagNone))); e != nil {
				break CompileSwitch
			}

			if e = r.Append(c.compileExpression(pair.Value, NewFlag(FlagNone))); e != nil {
				break CompileSwitch
			}
		}

		r.IL(ctx, opcode.IMakeHash, l).
			SetValues(1)

	case *ast.Identifier:
		p := c.compileIdentifierReference(n.Value, ctx, r)
		if p <= 0 {
			ctx := n.GetContext()
			e = NewSemanticError(ctx, "variable %s undefined", n.Value)
			break CompileSwitch
		}

		r.SetValues(1)

	case *ast.FunctionLiteral:
		if e = r.Append(c.compileFunctionLiteral(n)); e != nil {
			break CompileSwitch
		}

	case *ast.ExpressionList:
		isList := n.Length() > 1
		f := flag
		if isList {
			f.Set(FlagPackValue)
		}

		for _, expr := range n.Expressions {
			if e = r.Append(c.compileExpression(expr.Expression, f)); e != nil {
				break CompileSwitch
			}
		}

		if flag.Has(FlagPackValue) {
			r.IL(ctx, opcode.IMakeList, n.Length()).
				SetValues(1)
		}

	case *ast.InfixExpression:
		if e = r.Append(c.compileExpression(n.LeftOperand, flag.With(FlagPackValue))); e != nil {
			break CompileSwitch
		}

		if e = r.Append(c.compileExpression(n.RightOperand, flag.With(FlagPackValue))); e != nil {
			break CompileSwitch
		}

		r.IL(ctx, opcode.IBinOp, int(n.Operator.Token)).
			SetValues(1)

	case *ast.PrefixExpression:
		if e = r.Append(c.compileExpression(n.Operand, flag.With(FlagPackValue))); e != nil {
			break CompileSwitch
		}

		r.IL(ctx, opcode.IUniOp, int(n.Prefix.Token)).
			SetValues(1)

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
	}

	return r, e
}

func (c *Compiler) compileIdentifierReference(name string, ctx *token.Context, r *opcode.CodeBlock) int {
	ref, kind := c.Context.Variable.Reference(name)
	n := 1
	switch kind {
	case VariableKindGlobal, VariableKindModule:
		r.IL(ctx, opcode.ILoad, ref.Offset)

	case VariableKindBinding:
		r.IL(ctx, opcode.ILoadBind, ref.Offset)

	case VariableKindLocal:
		r.IL(ctx, opcode.ISLoad, ref.Offset)

	default:
		n = 0
	}

	return n
}

func (c *Compiler) compileIfExpression(n *ast.IfExpression) (*opcode.CodeBlock, error) {
	code, err := c.compileExpression(n.Condition, NewFlag(FlagNone))
	if err != nil {
		return nil, err
	}

	c.Context.Variable.EnterScope(FrameScopeBlock)
	consequence, err := c.compileStatement(n.Consequence, NewFlag(FlagCleanStack))
	if err != nil {
		return nil, err
	}
	c.Context.Variable.LeaveScope()

	var alternative *opcode.CodeBlock
	c.Context.Variable.EnterScope(FrameScopeBlock)
	if n.Alternative != nil {
		alternative, err = c.compileStatement(n.Alternative, NewFlag(FlagCleanStack))
		if err != nil {
			return nil, err
		}

	} else {
		alternative = opcode.NewCodeBlock()
		alternative.IL(n.GetContext(), opcode.ILoadNull)
		alternative.Values = 1
	}
	c.Context.Variable.LeaveScope()

	consequence.IL(n.Consequence.GetContext(), opcode.IJumpFWD, alternative.Length())
	code.IL(n.Consequence.GetContext(), opcode.IJumpIf, consequence.Length())
	code.Block(consequence)
	code.Block(alternative)

	return code, nil
}

func (c *Compiler) compileEmptyStatementBlock(ctx *token.Context, flag CompilerFlag) (*opcode.CodeBlock, error) {
	r := opcode.NewCodeBlock()
	r.IL(ctx, opcode.ILoadNull)
	r.Values = 1

	if flag.Has(FlagWithReturn) {
		r.IL(ctx, opcode.IReturn)
	}

	return r, nil
}

func (c *Compiler) compileStatements(ctx *token.Context, statements []ast.Statement, flag CompilerFlag) (*opcode.CodeBlock, error) {
	r := opcode.NewCodeBlock()
	var e error

	if len(statements) <= 0 {
		return c.compileEmptyStatementBlock(ctx, flag)
	}

	var last ast.Statement
	for i, stmt := range statements {
		if i < len(statements)-1 {
			nextFlag := NewFlag()
			_, isReturn := stmt.(*ast.ReturnStatement)
			if isReturn && r.Values > 0 {
				nextFlag.Set(FlagCleanStack)
			}

			if e = r.Append(c.compileStatement(stmt, nextFlag)); e != nil {
				break
			}
		} else {
			last = stmt
		}
	}

	lastContext := last.GetContext()
	nextFlag := NewFlag()
	count := r.Values
	if flag.Has(FlagWithReturn) || count > 0 {
		nextFlag.Set(FlagCleanStack)
	}

	lastResult, err := c.compileStatement(last, nextFlag)
	if err != nil {
		return nil, err
	}

	r.Block(lastResult)
	r.Values = lastResult.Values
	if flag.Has(FlagWithReturn) {
		r.IL(lastContext, opcode.IReturn)
	}

	return r, e
}

func (c *Compiler) compileFunctionLiteral(f *ast.FunctionLiteral) (*opcode.CodeBlock, error) {
	result := opcode.NewCodeBlock()
	c.Context.Variable.EnterScope(FrameScopeFunction)

	for _, item := range f.Arguments.Identifiers {
		c.Context.Variable.DefineArgument(item.Identifier.Value, item.Identifier.GetContext())
	}

	r, e := c.compileStatements(f.Body.GetContext(), f.Body.Statements, NewFlag(FlagWithReturn))
	if e != nil {
		return nil, e
	}

	frameSize := c.Context.Variable.CurrentScope().UpdateFrameSize(0)
	functionContext := &opcode.Function{
		FrameSize: frameSize,
		Arguments: f.Arguments.Length(),
		Codes:     r,
	}

	id := c.Context.AddFunction(functionContext)
	scope := c.Context.Variable.CurrentScope()
	c.Context.Variable.LeaveScope()

	for _, arg := range scope.BindingOrder {
		c.compileIdentifierReference(arg.Name, arg.Context, result)
	}

	result.IL(f.GetContext(), opcode.IMakeFunc, id, len(scope.Bindings))
	result.Values = 1
	return result, nil
}

func (c *Compiler) compileCallExpression(expr *ast.CallExpression, flag CompilerFlag) (*opcode.CodeBlock, error) {
	result := opcode.NewCodeBlock()

	args := opcode.NewCodeBlock()
	l := expr.Args.Length()
	for i := 0; i < l; i++ {
		a := expr.Args.Expressions[l-i-1]
		if e := args.Append(c.compileExpression(a.Expression, NewFlag(FlagPackValue))); e != nil {
			return nil, e
		}
	}
	result.Block(args)

	switch expr.Token.GetToken() {
	case token.Nil:
		callable, err := c.compileExpression(expr.Base, flag)
		if err != nil {
			return nil, err
		}
		result.Block(callable)

	case token.DualColon:
		callable, err := c.compileExpression(expr.Base, flag)
		if err != nil {
			return nil, err
		}

		memberIndex := c.Context.Literal.ReferenceString(expr.Member.Value)
		callable.IL(expr.Member.GetContext(), opcode.ILoad, int(memberIndex))
		callable.IL(expr.Base.GetContext(), opcode.ISDUP)
		callable.IL(expr.Member.GetContext(), opcode.IIndex)
		result.Block(callable)

	case token.Fn:
		result.IL(expr.Token.ToContext(), opcode.ISLoad, 0)
	}

	result.IL(expr.GetContext(), opcode.ICall, args.Values)
	result.Values = 1
	return result, nil
}

func (c *Compiler) compileIndexExpression(expr *ast.IndexExpression) (*opcode.CodeBlock, error) {
	result := opcode.NewCodeBlock()

	base, err := c.compileExpression(expr.Base, NewFlag(FlagPackValue))
	if err != nil {
		return nil, err
	}

	var index *opcode.CodeBlock

	if expr.Operator.Token == token.LBracket {
		index, err = c.compileExpression(expr.Index, NewFlag(FlagPackValue))
		if err != nil {
			return nil, err
		}
	} else {
		key := expr.Index.(*ast.Identifier)
		i := c.Context.Literal.ReferenceString(key.Value)
		index = opcode.NewCodeBlock()
		index.IL(expr.Index.GetContext(), opcode.ILoad, int(i))
		index.Values = 1
	}

	result.Block(base)
	result.Block(index)
	result.IL(expr.Index.GetContext(), opcode.IIndex)

	return result, nil
}
