package opcode

import (
	"github.com/flily/macaque-lang/token"
)

// IL of a statement of an expression.
type CodeBlock struct {
	Context *token.Context
	Codes   []IL
}

func (b *CodeBlock) Length() int {
	return len(b.Codes)
}

func (b *CodeBlock) Link() []Opcode {
	result := make([]Opcode, b.Length())

	for i, c := range b.Codes {
		result[i] = c.GetOpcode()
	}

	return result
}

type Function struct {
	Index        int
	FrameSize    int
	Arguments    int
	ReturnValues int
	IP           int64
	Codes        []*CodeBlock
	Contexts     []*token.Context
}

func (f *Function) Link() []Opcode {
	l := 0
	for _, c := range f.Codes {
		l += c.Length()
	}

	result := make([]Opcode, 0, l)
	contexts := make([]*token.Context, 0, l)

	for _, c := range f.Codes {
		result = append(result, c.Link()...)
		for i := 0; i < c.Length(); i++ {
			contexts = append(contexts, c.Context)
		}
	}

	f.Contexts = contexts
	return result
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
