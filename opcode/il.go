package opcode

import (
	"fmt"
)

type IL interface {
	GetCode() int
	GetOpcode() Opcode
	String() string

	elemCodeBlock()
}

func NewIL(code int, ops ...interface{}) IL {
	var r IL
	var err string

	switch code {
	case INOP, ILoadNull, IIndex, IClean, IReturn, IHalt, IScopeIn, IStackRev:
		r = ilCodeOp0(code)
		if len(ops) > 0 {
			err = fmt.Sprintf("code %s(%d) MUST NOT have operands", CodeName(code), code)
		}

	case IMakeFunc:
		var function, bindings int
		var ok bool
		if len(ops) != 2 {
			err = fmt.Sprintf("code %s(%d) MUST have 2 operands", CodeName(code), code)
			break
		}

		if function, ok = ops[0].(int); !ok {
			err = fmt.Sprintf("operand 0 MUST be int, got %T", ops[0])
			break
		}

		if bindings, ok = ops[1].(int); !ok {
			err = fmt.Sprintf("operand 1 MUST be int, got %T", ops[1])
			break
		}

		r = ilCodeMakeFunc{
			Function: function,
			Bindings: bindings,
		}

	case ILoad:
		if len(ops) != 1 {
			err = fmt.Sprintf("code %s(%d) MUST have 1 operand", CodeName(code), code)
			break
		}

		operand := ops[0]
		switch o := operand.(type) {
		case int:
			r = ilCodeOp1{
				Code:    code,
				Operand: o,
			}

		case uint64:
			r = ilCodeOp1{
				Code:    code,
				Operand: int(o),
			}

		case *DataContainer:
			r = &ilCodeLoad{
				Data: o,
			}

		default:
			err = fmt.Sprintf("operand MUST be int or *DataContainer, got %T", ops[0])
		}

	default:
		if len(ops) != 1 {
			err = fmt.Sprintf("code %s(%d) MUST have 1 operand", CodeName(code), code)
			break
		}

		operand, ok := ops[0].(int)
		if !ok {
			err = fmt.Sprintf("operand MUST be int, got %T", ops[0])
			break
		}

		r = ilCodeOp1{
			Code:    code,
			Operand: operand,
		}
	}

	if err != "" {
		panic(err)
	}

	return r
}

type ilCodeOp0 int

func (i ilCodeOp0) GetCode() int {
	return int(i)
}

func (i ilCodeOp0) GetOpcode() Opcode {
	return Code(int(i))
}

func (i ilCodeOp0) String() string {
	if i < ILastInst {
		return CodeName(int(i))
	}

	return fmt.Sprintf("INVALID[%d]", i)
}

func (i ilCodeOp0) elemCodeBlock() {}

type ilCodeOp1 struct {
	Code    int
	Operand int
}

func (i ilCodeOp1) GetCode() int {
	return i.Code
}

func (i ilCodeOp1) GetOpcode() Opcode {
	return Code(i.Code, i.Operand)
}

func (i ilCodeOp1) String() string {
	if i.Code < ILastInst {
		return fmt.Sprintf("%s %d", CodeName(i.Code), i.Operand)
	}

	return fmt.Sprintf("INVALID[%d]", i.Code)
}

func (i ilCodeOp1) elemCodeBlock() {}

type ilCodeMakeFunc struct {
	Module   string
	Function int
	Bindings int
	Info     *Function
}

func (i ilCodeMakeFunc) GetCode() int {
	return IMakeFunc
}

func (i ilCodeMakeFunc) GetOpcode() Opcode {
	return Code(IMakeFunc, i.Function, i.Bindings)
}

func (i ilCodeMakeFunc) String() string {
	return fmt.Sprintf("%s %d %d", CodeName(IMakeFunc), i.Function, i.Bindings)
}

func (i ilCodeMakeFunc) elemCodeBlock() {}

type ilCodeLoad struct {
	Data *DataContainer
}

func (i *ilCodeLoad) GetCode() int {
	return ILoad
}

func (i *ilCodeLoad) GetOpcode() Opcode {
	return Code(ILoad, int(i.Data.Index))
}

func (i *ilCodeLoad) String() string {
	return fmt.Sprintf("%s %d", CodeName(ILoad), i.Data.Index)
}

func (i *ilCodeLoad) elemCodeBlock() {}
