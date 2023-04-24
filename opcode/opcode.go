package opcode

import (
	"fmt"
)

const (
	IInvalid  = iota
	INOP      // No operation.
	ILoadInt  // Load an integer to the top of the stack.
	ILoadNull // Load a NULL to the top of the stack.
	ILoadBool // Load a boolean to the top of the stack.
	ILoadBind // Load a variable from function bound varaible to the top of the stack.
	ILoad     // Load a variable from data segment to the top of the stack.
	IPop      // Pop the top of the stack.
	ISLoad    // Load a variable from stack frame to the top of the stack.
	ISStore   // Store TOS to a local variable
	IBinOp    // Binary operation.
	IUniOp    // Unary operation.
	IMakeList // Make a list.
	IMakeHash // Make a hash.
	IMakeFunc // Make a function.
	IIndex    // Get item of a list or a hash.
	IJump     // Jump to a position.
	IJumpIf   // Jump to a position if TOS is false
	IJumpFWD  // Jump forward.
	ISDUP     // Duplicate the top of the stack.
	ICall     // Call a function.
	IClean    // Clean the stack.
	IReturn   // Return from a function.
	IHalt     // Halt the VM.
	ILastInst // Last instruction, no use.
)

var codeNames = [...]string{
	IInvalid:  "INVALID",
	INOP:      "NOP",
	ILoadInt:  "LOADINT",
	ILoadNull: "LOADNULL",
	ILoadBool: "LOADBOOL",
	ILoadBind: "LOADBIND",
	ILoad:     "LOAD",
	IPop:      "POP",
	ISLoad:    "SLOAD",
	ISStore:   "SSTORE",
	IBinOp:    "BINOP",
	IUniOp:    "UNIOP",
	IMakeList: "MAKELIST",
	IMakeHash: "MAKEHASH",
	IMakeFunc: "MAKEFUNC",
	IIndex:    "INDEX",
	IJump:     "JUMP",
	IJumpFWD:  "JUMFWD",
	IJumpIf:   "JUMPIF",
	ISDUP:     "SDUP",
	ICall:     "CALL",
	IClean:    "CLEAN",
	IReturn:   "RETURN",
	IHalt:     "HALT",
	ILastInst: "LASTINST",
}

func CodeName(code int) string {
	if code < 0 || code >= len(codeNames) {
		return "INVALID"
	}

	return codeNames[code]
}

// A temporary struct to hold the opcode and operands.
// It should be replaced by a more efficient type, may be uint64.
type Opcode struct {
	Name     int
	Operand0 int
	Operand1 int
}

func Code(name int, ops ...int) Opcode {
	i := Opcode{
		Name: name,
	}

	switch len(ops) {
	case 2:
		i.Operand1 = ops[1]
		fallthrough

	case 1:
		i.Operand0 = ops[0]
	}

	return i
}

func (i Opcode) String() string {
	name := CodeName(i.Name)
	return fmt.Sprintf("%s %d %d", name, i.Operand0, i.Operand1)
}
