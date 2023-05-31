package opcode

import (
	"github.com/flily/macaque-lang/token"
)

type CodeBlockItem interface {
	elemCodeBlock()
}

type ILContext struct {
	IL      IL
	Context *token.Context
}

// IL of a statement of an expression.
type CodeBlock struct {
	Codes  []ILContext
	Values int
}

func NewCodeBlock() *CodeBlock {
	b := &CodeBlock{}
	return b
}

func (b *CodeBlock) elemCodeBlock() {}

func (b *CodeBlock) IL(ctx *token.Context, code int, ops ...interface{}) *CodeBlock {
	context := ILContext{
		IL:      NewIL(code, ops...),
		Context: ctx,
	}

	b.Codes = append(b.Codes, context)
	return b
}

func (b *CodeBlock) Block(block *CodeBlock) *CodeBlock {
	b.Codes = append(b.Codes, block.Codes...)
	return b
}

func (b *CodeBlock) Append(block *CodeBlock, err error) error {
	if err != nil {
		return err
	}

	b.Block(block)
	return nil
}

func (b *CodeBlock) SetValues(values int) *CodeBlock {
	b.Values = values
	return b
}

func (b *CodeBlock) Length() int {
	return len(b.Codes)
}

func (b *CodeBlock) Link() ([]Opcode, []*token.Context) {
	codes := make([]Opcode, b.Length())
	debug := make([]*token.Context, b.Length())

	for i, c := range b.Codes {
		codes[i] = c.IL.GetOpcode()
		debug[i] = c.Context
	}

	return codes, debug
}

type Function struct {
	Index        int
	FrameSize    int
	Arguments    int
	ReturnValues int
	IP           uint64
	Codes        *CodeBlock
}

func (f *Function) Link() []Opcode {
	codes, _ := f.Codes.Link()
	return codes
}

// Codes of a module, which codes in a single file.
type Module struct {
	Name      string
	Canonical string
	Index     int
	File      *token.FileInfo
	Functions []*Function
}

func (m *Module) Link() []Opcode {
	result := make([]Opcode, 0)

	for _, f := range m.Functions {
		result = append(result, f.Link()...)
	}

	return result
}

func (m *Module) AddFunction(f *Function) {
	m.Functions = append(m.Functions, f)
}

func (m *Module) Main() *Function {
	return m.Functions[0]
}

// Collection of modules, and link to an executable file.
type Program struct {
	Modules   []*Module
	Functions []*Function
}

func (p *Program) LinkFunctions() {
	p.Functions = make([]*Function, 1)

	for _, m := range p.Modules {
		for _, f := range m.Functions {
			f.Index = len(p.Functions)
			p.Functions = append(p.Functions, f)
		}
	}
}
