package vm

import (
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

type NaiveVM struct {
	Code  []opcode.Opcode
	Stack []object.Object
	Data  []object.Object

	ip uint64
	sp uint64
	bp uint64

	AX int64
}

func NewNaiveVM() *NaiveVM {
	m := &NaiveVM{
		Code:  make([]opcode.Opcode, 0),
		Stack: make([]object.Object, DefaultStackSize),
		Data:  make([]object.Object, DefaultDataSize),
	}

	return m
}

func (m *NaiveVM) stackPush(o object.Object) {
	m.Stack[m.sp] = o
	m.sp++
}

func (m *NaiveVM) stackPop() object.Object {
	m.sp--
	r := m.Stack[m.sp]
	m.Stack[m.sp] = nil
	return r
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

func (m *NaiveVM) Top() object.Object {
	return m.Stack[m.sp-1]
}

func (m *NaiveVM) LoadCode(c *compiler.Compiler) {
	m.Code = append(m.Code, c.Context.Instructions...)
	m.Code = append(m.Code, opcode.Inst(opcode.IHalt))
}

func (m *NaiveVM) LoadData(c *compiler.Compiler) {
	copy(m.Data, c.Context.Literal.Values)
}

func (m *NaiveVM) SetEntry(ip uint64) {
	m.ip = ip
	// Use a null as placeholder for current function
	m.stackPush(object.NewNull())
}

func (m *NaiveVM) Run(ip uint64) error {
	m.SetEntry(ip)

	var e error
RunSwitch:
	for {
		op := m.fetchOp()

		switch op.Name {
		case opcode.ILoadInt:
			o := object.NewInteger(int64(op.Operand0))
			m.stackPush(o)

		case opcode.ILoadStr:
			o := m.refData(uint64(op.Operand0))
			m.stackPush(o)

		case opcode.IBinOp:
			operator := token.Token(op.Operand0)
			right := m.stackPop()
			left := m.stackPop()
			o, ok := left.OnInfix(operator, right)
			if !ok {
				e = errors.NewError(errors.ErrRuntimeError,
					"%s %s %s is not accepted", left.Type(), operator, right.Type())
				break RunSwitch
			}
			m.stackPush(o)

		case opcode.ISStore:
			o := m.stackPop()
			offset := op.Operand0
			m.localBind(int64(offset), o)

		case opcode.ISLoad:
			offset := op.Operand0
			o := m.localRead(int64(offset))
			m.stackPush(o)

		case opcode.ILoad:
			offset := op.Operand0
			o := m.Data[offset]
			m.stackPush(o)

		case opcode.IHalt:
			break RunSwitch
		}
	}

	return e
}
