package vm

import (
	"fmt"
	"strings"

	"github.com/flily/macaque-lang/compiler"
	"github.com/flily/macaque-lang/errors"
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
}

type NaiveVM struct {
	Code []opcode.Opcode
	Data []object.Object

	ip uint64 // instruction pointer
	sp uint64 // stack pointer
	sb uint64 // stack base pointer
	bp uint64 // base pointer

	Stack     []object.Object
	callStack []callStackInfo
	csi       uint64
	Functions []compiler.FunctionInfo

	AX int64
}

func NewNaiveVM() *NaiveVM {
	m := &NaiveVM{
		Code:      make([]opcode.Opcode, 0),
		Data:      make([]object.Object, DefaultDataSize),
		Stack:     make([]object.Object, DefaultStackSize),
		callStack: make([]callStackInfo, DefaultStackSize),
	}

	return m
}

func (m *NaiveVM) stackPush(o object.Object) {
	m.Stack[m.sp] = o
	m.sp++
}

func (m *NaiveVM) stackPop() object.Object {
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

func (m *NaiveVM) stackPopN(n uint64) {
	for m.sp > 0 && n > 0 {
		n--
		m.sp--
		m.Stack[m.sp] = nil
	}
}

func (m *NaiveVM) stackPushN(objects []object.Object) {
	for _, o := range objects {
		m.stackPush(o)
	}
}

func (m *NaiveVM) stackPopNWithValue(n int) []object.Object {
	r := make([]object.Object, n)
	for i := 0; i < n; i++ {
		r[n-i-1] = m.stackPop()
	}

	return r
}

func (m *NaiveVM) incrIP(n uint64) {
	m.ip += n
}

func (m *NaiveVM) localBind(i int64, o object.Object) {
	m.Stack[int64(m.bp)+i] = o
}

func (m *NaiveVM) localRead(i int64) object.Object {
	return m.Stack[int64(m.bp)+i]
}

func (m *NaiveVM) fetchOp() opcode.Opcode {
	r := m.Code[m.ip]
	m.ip++
	return r
}

func (m *NaiveVM) refData(i uint64) object.Object {
	return m.Data[i]
}

func (m *NaiveVM) pushCallInfo() {
	m.callStack[m.csi].bp = m.bp
	m.callStack[m.csi].sp = m.sp
	m.callStack[m.csi].sb = m.sb
	m.callStack[m.csi].ip = m.ip
	m.csi++
}

func (m *NaiveVM) popCallInfo() {
	m.csi--
	m.bp = m.callStack[m.csi].bp
	m.sp = m.callStack[m.csi].sp
	m.sb = m.callStack[m.csi].sb
	m.ip = m.callStack[m.csi].ip
}

func (m *NaiveVM) initCallStack(frameSize int) {
	m.bp = m.sp - 1
	for i := 0; i < frameSize; i++ {
		m.stackPush(null)
	}
	m.sb = m.sp
}

func (m *NaiveVM) Top() object.Object {
	return m.Stack[m.sp-1]
}

func (m *NaiveVM) LoadCodePage(page *compiler.CodePage) {
	m.LoadFunctions(page)
	m.LoadCode(page)
	m.LoadData(page)
}

func (m *NaiveVM) LoadFunctions(page *compiler.CodePage) {
	m.Functions = make([]compiler.FunctionInfo, len(page.Functions))
	copy(m.Functions, page.Functions)
}

func (m *NaiveVM) LoadCode(page *compiler.CodePage) {
	m.Code = make([]opcode.Opcode, len(page.Codes))
	copy(m.Code, page.Codes)
	if m.Code[len(m.Code)-1].Name != opcode.IHalt {
		m.Code = append(m.Code, opcode.Code(opcode.IHalt))
	}
}

func (m *NaiveVM) LoadData(c *compiler.CodePage) {
	m.Data = make([]object.Object, len(c.Data))
	copy(m.Data, c.Data)
}

func (m *NaiveVM) SetEntry(entry *object.FunctionObject) {
	m.stackPush(entry)
	m.ip = entry.IP
	for i := 0; i < entry.FrameSize; i++ {
		m.stackPush(null)
	}
	m.sb = m.sp
}

func (m *NaiveVM) GetFunctionInfo(i int) (compiler.FunctionInfo, bool) {
	if i < 0 || i >= len(m.Functions) {
		return compiler.FunctionInfo{}, false
	}

	return m.Functions[i], true
}

func (m *NaiveVM) ExecOpcode(op opcode.Opcode) error {
	var e error

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
			e = errors.NewError(errors.ErrCodeRuntimeError,
				"%s %s %s is not accepted", left.Type(), operator, right.Type())
			break
		}
		m.stackPush(o)

	case opcode.IUniOp:
		operator := token.Token(op.Operand0)
		operand := m.stackPop()
		o, ok := operand.OnPrefix(operator)
		if !ok {
			e = errors.NewError(errors.ErrCodeRuntimeError,
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
			e = errors.NewError(errors.ErrCodeRuntimeError,
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
			e = errors.NewError(errors.ErrCodeRuntimeError,
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
			e = errors.NewError(errors.ErrCodeRuntimeError,
				"%s is not callable", f.Type())
			break
		}

		fn := f.(*object.FunctionObject)
		m.pushCallInfo()
		m.initCallStack(fn.FrameSize)
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
		break
	}

	return e
}

func (m *NaiveVM) Run(entry *object.FunctionObject) error {
	m.SetEntry(entry)

	codeSize := uint64(len(m.Code))
	var e error
	for m.ip < codeSize && e == nil {
		op := m.fetchOp()
		e = m.ExecOpcode(op)
	}

	return e
}

func (m *NaiveVM) GetAllStack() []object.Object {
	return m.Stack[:m.sp]
}

func (m *NaiveVM) InspectStack() (string, string) {
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
