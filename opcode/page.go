package opcode

import (
	"github.com/flily/macaque-lang/object"
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
	b.Values += block.Values
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

func (b *CodeBlock) Link(postfixes ...Opcode) ([]Opcode, []*token.Context) {
	codes := make([]Opcode, b.Length()+len(postfixes))
	debug := make([]*token.Context, b.Length())

	for i, c := range b.Codes {
		codes[i] = c.IL.GetOpcode()
		debug[i] = c.Context
	}

	copy(codes[b.Length():], postfixes)
	return codes, debug
}

type Function struct {
	ModuleIndex  uint64
	GlobalIndex  uint64
	FrameSize    int
	Arguments    int
	ReturnValues int
	IP           uint64
	Codes        *CodeBlock
	Opcodes      []Opcode
	DebugInfo    []*token.Context
}

func (f *Function) Func(bounds []object.Object) *object.FunctionObject {
	function := &object.FunctionObject{
		Index:     f.GlobalIndex,
		FrameSize: f.FrameSize,
		Arguments: f.Arguments,
		IP:        f.IP,
		Bounds:    bounds,
	}

	return function
}

func (f *Function) IsLink() bool {
	return f.Opcodes != nil
}

func (f *Function) Link(postfix ...Opcode) []Opcode {
	if !f.IsLink() {
		codes, debug := f.Codes.Link(postfix...)
		f.Opcodes = codes
		f.DebugInfo = debug
	}

	return f.Opcodes
}

const (
	ModuleTypeNoSet = iota
	ModuleTypeSystem
	ModuleTypeImported
	ModuleTypeUser
)

// Codes of a module, which codes in a single file.
type Module struct {
	Name      string
	Canonical string
	Type      int
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

// Collection of all modules, and link to an executable file.
type CodePage struct {
	NativeModules []*Module
	ModuleNameMap map[string]*Module
	Functions     []*Function
	Data          []object.Object
}

func NewCodePage() *CodePage {
	p := &CodePage{}
	p.ModuleNameMap = make(map[string]*Module)
	return p
}

func (p *CodePage) Main() *Function {
	return p.Functions[0]
}

func (p *CodePage) LinkFunctions() {
	p.Functions = make([]*Function, 1)

	for _, m := range p.NativeModules {
		for _, f := range m.Functions {
			f.GlobalIndex = uint64(len(p.Functions))
			p.Functions = append(p.Functions, f)
		}
	}
}

func (p *CodePage) LinkCode() []Opcode {
	var post []Opcode
	if len(p.Functions) > 1 {
		post = []Opcode{Code(IHalt)}
	}

	code := make([]Opcode, 0)

	for _, f := range p.Functions {
		f.IP = uint64(len(code))
		fc := f.Link(post...)
		code = append(code, fc...)
	}

	return code
}
