package vm

import (
	"fmt"
	"strings"

	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/opcode"
	"github.com/flily/macaque-lang/token"
)

const (
	DefaultStackSize = 65536
	DefaultDataSize  = 65536
)

var (
	null = object.NewNull()
)

type callStackInfo struct {
	bp uint64
	sp uint64
	sb uint64
	ip uint64
	fi uint64
	fp uint64
}

type VM interface {
	LoadCodePage(page *opcode.CodePage)
	GetSP() uint64
	GetStackObject(i int) object.Object
	GetRegister(name string) uint64
	Run(entry *object.FunctionObject) error
}

type NaiveVMBase struct {
	Data []object.Object

	ip uint64 // instruction pointer
	sp uint64 // stack pointer
	sb uint64 // stack base pointer
	bp uint64 // base pointer
	fi uint64 // function index
	fp uint64 // function pointer

	Stack     []object.Object
	callStack []callStackInfo
	csi       uint64
	Functions []*opcode.Function

	AX int64
}

func NewNaiveVMBase() *NaiveVMBase {
	m := &NaiveVMBase{
		Data:      make([]object.Object, DefaultDataSize),
		Stack:     make([]object.Object, DefaultStackSize),
		callStack: make([]callStackInfo, DefaultStackSize),
	}

	return m
}

func (m *NaiveVMBase) execOne() {
	m.ip++
}

func (m *NaiveVMBase) stackPush(o object.Object) {
	m.Stack[m.sp] = o
	m.sp++
}

func (m *NaiveVMBase) stackPop() object.Object {
	var r object.Object
	if m.sp > m.sb {
		m.sp--
		r = m.Stack[m.sp]
		m.Stack[m.sp] = nil
	} else {
		r = null
	}
	return r
}

func (m *NaiveVMBase) stackPopN(n uint64) {
	for m.sp > 0 && n > 0 {
		n--
		m.sp--
		m.Stack[m.sp] = nil
	}
}

func (m *NaiveVMBase) stackPushN(objects []object.Object) {
	for _, o := range objects {
		m.stackPush(o)
	}
}

func (m *NaiveVMBase) stackPopNWithValue(n int) []object.Object {
	r := make([]object.Object, n)
	for i := 0; i < n; i++ {
		r[n-i-1] = m.stackPop()
	}

	return r
}

func (m *NaiveVMBase) incrIP(n uint64) {
	m.ip += n
}

func (m *NaiveVMBase) localBind(i int64, o object.Object) {
	m.Stack[int64(m.bp)+i] = o
}

func (m *NaiveVMBase) localRead(i int64) object.Object {
	return m.Stack[int64(m.bp)+i]
}

func (m *NaiveVMBase) refData(i uint64) object.Object {
	return m.Data[i]
}

func (m *NaiveVMBase) pushCallInfo() {
	m.callStack[m.csi].bp = m.bp
	m.callStack[m.csi].sp = m.sp
	m.callStack[m.csi].sb = m.sb
	m.callStack[m.csi].ip = m.ip
	m.callStack[m.csi].fi = m.fi
	m.callStack[m.csi].fp = m.fp
	m.csi++
}

func (m *NaiveVMBase) popCallInfo() {
	m.csi--
	m.bp = m.callStack[m.csi].bp
	m.sp = m.callStack[m.csi].sp
	m.sb = m.callStack[m.csi].sb
	m.ip = m.callStack[m.csi].ip
	m.fi = m.callStack[m.csi].fi
	m.fp = m.callStack[m.csi].fp
}

func (m *NaiveVMBase) initCallStack(frameSize int) {
	m.bp = m.sp - 1
	for i := 0; i < frameSize; i++ {
		m.stackPush(null)
	}
	m.sb = m.sp
}

func (m *NaiveVMBase) Top() object.Object {
	return m.Stack[m.sp-1]
}

func (m *NaiveVMBase) GetRegister(name string) uint64 {
	var r uint64
	switch name {
	case "ip":
		r = m.ip

	case "sp":
		r = m.sp

	case "sb":
		r = m.sb

	case "bp":
		r = m.bp

	case "fi":
		r = m.fi

	case "fp":
		r = m.fp
	}

	return r
}

func (m *NaiveVMBase) GetSP() uint64 {
	return m.sp
}

func (m *NaiveVMBase) GetStackObject(i int) object.Object {
	if i < 0 || i >= int(m.sp) {
		return nil
	}

	return m.Stack[i]
}

func (m *NaiveVMBase) SetEntry(entry *object.FunctionObject) {
	m.stackPush(entry)
	m.ip = entry.IP
	for i := 0; i < entry.FrameSize; i++ {
		m.stackPush(null)
	}
	m.sb = m.sp
}

func (m *NaiveVMBase) GetFunctionInfo(i int) (*opcode.Function, bool) {
	if i < 0 || i >= len(m.Functions) {
		return nil, false
	}

	return m.Functions[i], true
}

func (m *NaiveVMBase) ExecOpcode(op opcode.Opcode) (error, bool) {
	var e error
	var isHalt bool

	// fmt.Printf("OP %s\n", op)
	// vs, vv := m.InspectStack()
	// fmt.Printf("STACK %s\n", vs)
	// fmt.Printf("      %s\n", vv)

	switch op.Name {
	case opcode.ILoadNull:
		m.stackPush(null)

	case opcode.ILoadBool:
		o := object.NewBoolean(op.Operand0 != 0)
		m.stackPush(o)

	case opcode.ILoadInt:
		o := object.NewInteger(int64(op.Operand0))
		m.stackPush(o)

	case opcode.ILoadBind:
		o := m.localRead(0)
		f := o.(*object.FunctionObject)
		i := int64(op.Operand0)
		m.stackPush(f.Bounds[i])

	case opcode.ISStore:
		o := m.stackPop()
		offset := op.Operand0
		m.localBind(int64(offset), o)

	case opcode.ISLoad:
		offset := op.Operand0
		o := m.localRead(int64(offset))
		m.stackPush(o)

	case opcode.ILoad:
		o := m.refData(uint64(op.Operand0))
		m.stackPush(o)

	case opcode.IPop:
		m.stackPopN(uint64(op.Operand0))

	case opcode.IBinOp:
		operator := token.Token(op.Operand0)
		right := m.stackPop()
		left := m.stackPop()
		o, ok := left.OnInfix(operator, right)
		if !ok {
			e = NewRuntimeError(
				"%s %s %s is not accepted", left.Type(), operator, right.Type())
			break
		}
		m.stackPush(o)

	case opcode.IUniOp:
		operator := token.Token(op.Operand0)
		operand := m.stackPop()
		o, ok := operand.OnPrefix(operator)
		if !ok {
			e = NewRuntimeError(
				"%s %s is not accepted", operator, operand.Type())
			break
		}
		m.stackPush(o)

	case opcode.IMakeList:
		n := op.Operand0
		array := make([]object.Object, n)
		for i := 0; i < n; i++ {
			array[n-1-i] = m.stackPop()
		}
		o := object.NewArray(array)
		m.stackPush(o)

	case opcode.IMakeHash:
		n := op.Operand0
		hash := make([]object.HashPair, n)
		for i := 0; i < n; i++ {
			value := m.stackPop()
			key := m.stackPop()
			item := object.HashPair{
				Key:   key,
				Value: value,
			}
			hash[n-1-i] = item
		}

		o := object.NewHash(hash)
		m.stackPush(o)

	case opcode.IMakeFunc:
		index := op.Operand0
		info, ok := m.GetFunctionInfo(index)
		if !ok {
			e = NewRuntimeError(
				"function %d not found", index)
			break
		}

		n := op.Operand1
		bounds := make([]object.Object, n)
		for i := 0; i < n; i++ {
			bounds[n-1-i] = m.stackPop()
		}

		o := info.Func(bounds)
		m.stackPush(o)

	case opcode.IIndex:
		index := m.stackPop()
		base := m.stackPop()
		o, ok := base.OnIndex(index)
		if !ok {
			e = NewRuntimeError(
				"%s[%s] is not accepted", base.Type(), index.Type())
			break
		}
		m.stackPush(o)

	case opcode.IJumpFWD:
		m.incrIP(uint64(op.Operand0))

	case opcode.IJumpIf:
		o := m.stackPop()
		notJump := false
		switch obj := o.(type) {
		case *object.NullObject:
			notJump = false
		case *object.BooleanObject:
			notJump = obj.Value
		default:
			notJump = true
		}
		if !notJump {
			m.incrIP(uint64(op.Operand0))
		}

	case opcode.ISDUP:
		top := m.Top()
		m.stackPush(top)

	case opcode.ICall:
		f := m.Top()
		if f.Type() != object.ObjectTypeFunction {
			e = NewRuntimeError(
				"%s is not callable", f.Type())
			break
		}

		fn := f.(*object.FunctionObject)
		m.pushCallInfo()
		m.initCallStack(fn.FrameSize)
		m.fi = fn.Index
		m.fp = fn.IP
		m.ip = fn.IP

	case opcode.IClean:
		n := m.sp - m.sb
		m.stackPopN(n)

	case opcode.IReturn:
		n := int(m.sp - m.sb)
		returnValues := m.stackPopNWithValue(n)
		m.popCallInfo()

		f := m.stackPop()
		fn := f.(*object.FunctionObject)
		args := fn.Arguments
		m.stackPopN(uint64(args))
		m.stackPushN(returnValues)

	case opcode.IHalt:
		isHalt = true
	}

	return e, isHalt
}

func (m *NaiveVMBase) GetAllStack() []object.Object {
	return m.Stack[:m.sp]
}

func (m *NaiveVMBase) InspectStack() (string, string) {
	items := m.GetAllStack()
	data := make([]string, len(items))
	view := make([]string, len(items))
	for i, item := range items {
		o := fmt.Sprintf("%-6s", item.Inspect())
		data[i] = "| " + o + " "

		var p string
		switch uint64(i) {
		case m.sp:
			p = "SP"

		case m.sb:
			p = "SB"

		case m.bp:
			p = "BP"

		default:
			p = fmt.Sprintf("%d", i)
		}

		view[i] = "| " + p + strings.Repeat(" ", len(o)-len(p)) + " "
	}

	return strings.Join(data, ""), strings.Join(view, "")
}

type NaiveVM struct {
	NaiveVMBase
	Code []opcode.Opcode
}

func NewNaiveVM() *NaiveVM {
	m := &NaiveVM{
		NaiveVMBase: *NewNaiveVMBase(),
		Code:        make([]opcode.Opcode, 0),
	}

	return m
}

func (m *NaiveVM) fetchOp() opcode.Opcode {
	r := m.Code[m.ip]
	m.execOne()
	return r
}

func (m *NaiveVM) Run(entry *object.FunctionObject) error {
	m.SetEntry(entry)

	codeSize := uint64(len(m.Code))
	var e error
	var isHalt bool
	for m.ip < codeSize && e == nil && !isHalt {
		op := m.fetchOp()
		// f := m.fi
		// info := m.Functions[f].DebugInfo[m.fi]
		// fmt.Printf("%s\n", info.Message("%s", op))

		e, isHalt = m.ExecOpcode(op)
	}

	return e
}

func (m *NaiveVM) loadFunctions(page *opcode.CodePage) {
	m.Functions = make([]*opcode.Function, len(page.Functions))
	copy(m.Functions, page.Functions)
}

func (m *NaiveVM) loadCode(page *opcode.CodePage) {
	codes := page.LinkCode()
	m.Code = make([]opcode.Opcode, len(codes))
	copy(m.Code, codes)
	if m.Code[len(m.Code)-1].Name != opcode.IHalt {
		m.Code = append(m.Code, opcode.Code(opcode.IHalt))
	}
}

func (m *NaiveVM) loadData(c *opcode.CodePage) {
	m.Data = make([]object.Object, len(c.Data))
	copy(m.Data, c.Data)
}

func (m *NaiveVM) LoadCodePage(page *opcode.CodePage) {
	m.loadFunctions(page)
	m.loadCode(page)
	m.loadData(page)
}

type NaiveVMInterpreter struct {
	NaiveVMBase
	CodePage *opcode.CodePage
}

func NewNaiveVMInterpreter() *NaiveVMInterpreter {
	m := &NaiveVMInterpreter{
		NaiveVMBase: *NewNaiveVMBase(),
	}

	return m
}

func (i *NaiveVMInterpreter) LoadCodePage(page *opcode.CodePage) {
	i.CodePage = page
	i.Data = page.Data
}

func (i *NaiveVMInterpreter) getFunction(o object.Object) (*opcode.Function, error) {
	f, ok := o.(*object.FunctionObject)
	if !ok {
		return nil, NewRuntimeError(
			"%s is not callable", f.Type())
	}

	fn, ok := i.GetFunctionInfo(int(f.Index))
	if !ok {
		return nil, NewRuntimeError(
			"function %d not found", f.Index)
	}

	return fn, nil
}

func (i *NaiveVMInterpreter) runFunction(f *opcode.Function) (error, bool) {
	var e error
	var isHalt bool
	var breakAndReturn bool

	length := len(f.Opcodes)

	for i.ip-i.fp < uint64(length) && e == nil && !isHalt {
		j := int(i.ip - i.fp)
		i.execOne()

		code := f.Opcodes[j]
		// info := f.DebugInfo[j]
		// fmt.Printf("%s\n", info.Message("%s", code))
		top := i.Top()
		e, isHalt = i.ExecOpcode(code)
		if e != nil {
			break
		}

		if isHalt {
			break
		}

		switch code.Name {
		case opcode.ICall:
			fn, err := i.getFunction(top)
			if err != nil {
				e = err
				break
			}

			e, isHalt = i.runFunction(fn)

		case opcode.IReturn:
			breakAndReturn = true
		}

		if breakAndReturn {
			break
		}
	}

	return e, isHalt
}

func (i *NaiveVMInterpreter) Run(entry *object.FunctionObject) error {
	i.SetEntry(entry)

	index := int(entry.Index)
	if index < 0 || index >= len(i.CodePage.Functions) {
		return NewRuntimeError("function %d not found", entry.Index)
	}

	fn := i.CodePage.Functions[index]
	err, _ := i.runFunction(fn)

	return err
}
