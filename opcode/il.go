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

type ilCodeOp2 struct {
	Code     int
	Operand0 int
	Operand1 int
}

func (i ilCodeOp2) GetCode() int {
	return i.Code
}

func (i ilCodeOp2) GetOpcode() Opcode {
	return Code(i.Code, i.Operand0, i.Operand1)
}

type ilCodeMakeFunc struct {
}
