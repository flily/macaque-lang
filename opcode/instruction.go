package opcode

const (
	IInvalid  = iota
	INOP      // No operation.
	ILoadInt  // Load an integer to the top of the stack.
	INeg      // Negate the top of the stack.
	INot      // Logical not the top of the stack.
	IInvert   // Bitwise not the top of the stack.
	IBinOp    // Binary operation.
	IHalt     // Halt the VM.
	ILastInst // Last instruction, no use.
)

var instNames = [...]string{
	IInvalid:  "INVALID",
	INOP:      "NOP",
	ILoadInt:  "LOADINT",
	INeg:      "NEG",
	INot:      "NOT",
	IInvert:   "INVERT",
	IBinOp:    "BINOP",
	IHalt:     "HALT",
	ILastInst: "LASTINST",
}

const (
	XInstInvalid    = iota
	XLoadIntLiteral // Load integer from instruction.
	XLoadIntData    // Load integer from data section.
)

func InstName(code int) string {
	if code < 0 || code >= len(instNames) {
		return "INVALID"
	}

	return instNames[code]
}

type Instruction struct {
	Name     int
	Operand0 int
	Operand1 int
	Operand2 int
}

func Inst(name int, ops ...int) Instruction {
	i := Instruction{
		Name: name,
	}

	switch len(ops) {
	case 3:
		i.Operand2 = ops[2]
		fallthrough

	case 2:
		i.Operand1 = ops[1]
		fallthrough

	case 1:
		i.Operand0 = ops[0]
	}

	return i
}
