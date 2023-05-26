package opcode

import "fmt"

type IL interface {
	GetCode() int
	GetOpcode() Opcode
}

func NewIL(code int, ops ...interface{}) IL {
	var r IL
	var err string

	switch code {
	case INOP, ILoadNull, IIndex, IClean, IReturn, IHalt:
		r = ilCodeOp0(code)
		if len(ops) > 0 {
			err = fmt.Sprintf("code %d MUST NOT have operands", code)
		}

	case IMakeFunc:
		var function, bindings int
		var ok bool
		if len(ops) != 2 {
			err = fmt.Sprintf("code %d MUST have 2 operands", code)
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

	default:
		if len(ops) != 1 {
			err = fmt.Sprintf("code %d MUST have 1 operand", code)
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
