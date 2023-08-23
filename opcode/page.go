package opcode

import (
	"fmt"
	"strings"

	"github.com/flily/macaque-lang/ast"
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

type DataContainer struct {
	Data  object.Object
	Index uint64
}

func NewDataContainer(data object.Object, index uint64) *DataContainer {
	c := &DataContainer{
		Data:  data,
		Index: index,
	}

	return c
}

// IL of a statement of an expression.
type CodeBlock struct {
	Codes      []ILContext
	Values     int
	Determined bool
}

func NewCodeBlock() *CodeBlock {
	b := &CodeBlock{
		Determined: true,
	}
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
	if !block.Determined {
		b.Determined = false
	}
	return b
}

func (b *CodeBlock) Append(block *CodeBlock, err error) error {
	if err != nil {
		return err
	}

	b.Block(block)
	return nil
}

func (b *CodeBlock) PrependIL(ctx *token.Context, code int, ops ...interface{}) *CodeBlock {
	context := ILContext{
		IL:      NewIL(code, ops...),
		Context: ctx,
	}

	b.Codes = append([]ILContext{context}, b.Codes...)
	return b
}

func (b *CodeBlock) Prepend(block *CodeBlock, err error) error {
	if err != nil {
		return err
	}

	b.Codes = append(block.Codes, b.Codes...)
	b.Values += block.Values
	return nil
}

func (b *CodeBlock) SetValues(values int) *CodeBlock {
	b.Values = values
	return b
}

func (b *CodeBlock) Undetermined() *CodeBlock {
	b.Determined = false
	return b
}

func (b *CodeBlock) CleanStack() *CodeBlock {
	b.Values = 0
	b.Determined = true
	return b
}

func (b *CodeBlock) Length() int {
	return len(b.Codes)
}

func (b *CodeBlock) String() string {
	lines := make([]string, len(b.Codes))
	maxLength := len(fmt.Sprintf("%d", len(b.Codes))) + 1
	format := fmt.Sprintf("%%%dd    %%s", maxLength)

	for i, c := range b.Codes {
		lines[i] = fmt.Sprintf(format, i, c.IL.String())
	}

	return fmt.Sprintf("CodeBlock[%d]{\n%s\n}", b.Length(), strings.Join(lines, "\n"))
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
		f.Relink(postfix...)
	}

	return f.Opcodes
}

func (f *Function) Relink(postfix ...Opcode) []Opcode {
	codes, debug := f.Codes.Link(postfix...)
	f.Opcodes = codes
	f.DebugInfo = debug

	return codes
}

func (f *Function) Append(block *CodeBlock) {
	f.Codes.Block(block)
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
	Root      ast.Node
	Functions []*Function
	Data      []*DataContainer
}

func NewModule(filename string, root ast.Node) *Module {
	m := &Module{
		File:      token.NewFileInfo(filename),
		Root:      root,
		Functions: make([]*Function, 0),
	}

	return m
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
	dataIndex     map[interface{}]uint64
}

func BuildCodePage(modules []*Module) *CodePage {
	moduleList := make([]*Module, len(modules))
	copy(moduleList, modules)
	page := &CodePage{
		NativeModules: moduleList,
		dataIndex:     make(map[interface{}]uint64),
	}

	return page
}

func NewCodePage() *CodePage {
	p := &CodePage{}
	p.ModuleNameMap = make(map[string]*Module)
	return p
}

func (p *CodePage) MainModule() *Module {
	return p.NativeModules[0]
}

func (p *CodePage) Main() *Function {
	return p.Functions[0]
}

func (p *CodePage) UpsertData(data object.Object) uint64 {
	index, ok := p.dataIndex[data]
	if ok {
		return index
	}

	index = uint64(len(p.Data))
	p.Data = append(p.Data, data)
	p.dataIndex[data] = index

	return index
}

func (p *CodePage) LinkModule(m *Module) {
	for _, data := range m.Data {
		index := p.UpsertData(data.Data)
		data.Index = index
	}

	p.Functions = append(p.Functions, m.Functions...)
}

func (p *CodePage) LinkModules() {
	for _, m := range p.NativeModules {
		p.LinkModule(m)
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
