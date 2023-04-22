Macaque Language Naive Virtual Machine Specification
=====================================================

Overview
---------

Macaque language uses a stack-based virtual machine. We decide to implement two
different VMs. One is a naive VM, executes instructions with all features
implemented, to validated the lanaugae design, but not focus on performance.
The other is a turbo VM engine with performance optimization, will be designed
after the simple VM is fully implemented.


Memory layout
--------------

Macaque VM comes with 4 memory segments:
  1. **Code segment**: contains all the bytecode of the program. All compiled
     codes are stored in this segment in a sequence. At the end of each function,
     a `HALT` instruction is inserted.
  2. **Data segment**: contains constants and some static data.
  3. **Stack segment**: contains local variable data.
  4. **Call stack**: contains the return address and info of each function call.
     Call stack is a VM managed stack, and can not be accessed directly.


Function call
--------------

For example, we have a function `f(a, b)`:
```monkey
let c, d = 5, 7
let result = f(a, b) {
    let e = a - c
    let f = b + d
    return e + f
}
```
To this function, there are 2 parameters `a` and `b`, 2 bound variables `c` and
`d` from outer scope, and 2 local variables `e` and `f`.


Function calls `f(a, b)` are performed with following steps:
  1. Push protection object onto the stack.
     ```
                           SP
     +---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+
     | . | . | . | . | P |   |   |   |   |   |   |   |   |   |   |   |   |   |
     +---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+
           BP
     ```
  2. Evaluation each arguments of the function call.
  3. Push each value of arguments onto the stack, from right to left.
     ```
                                   SP
     +---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+
     | . | . | . | . | P | b | a |   |   |   |   |   |   |   |   |   |   |   |
     +---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+
           BP
     ```
  4. Push function (closure) object onto the stack.
     ```
                                       SP
     +---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+
     | . | . | . | . | P | b | a | f |   |   |   |   |   |   |   |   |   |   |
     +---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+
           BP
     ```
  5. Execute instruction `CALL f`, operand means the number of arguments.
     1. Push base pointer (BP) and stack pointer (SP) into the call stack.
     2. Set base pointer (BP) to the top of the stack (SP-1).
     3. Read function bound variables and push them onto the stack.
     4. Shift stack pointer (SP) to make room for local variables.
     5. Push return address onto the call stack.
     ```
                                                       SP
     +---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+
     | . | . | . | . | P | b | a | f | c | d | e | f |   |   |   |   |   |   |
     +---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+
                                   BP
     +----+----+----+----+----+----+----+
     | BP | SP | IP |    |    |    |    | call stack
     +----+----+----+----+----+----+----+
     ```


Instructions
-------------

Letters `NBbWwDd` are used to represent the type of operands. Upper case letter
represents an unsigned integer, and lower case represents a signed integer.
  + `N`: nil, ALL-ZERO byte.
  + `B`, `b`: byte, an 8-bit integer.
  + `W`, `w`: word, a 16-bit integer.
  + `D`, `d`: quad-byte word, a 32-bit integer.

| Mnemonic | Operands | Description                                           |
|:--------:|:--------:|:------------------------------------------------------|
| NOP      |   NNN    | No operation
| LOADINT  |   D      | Load an integer constant onto the stack
| LOADNULL |   NNN    | Load a null object onto the stack
| LOADBOOL |   B      | Load a boolean object onto the stack
| LOADBIND |   D      | Load a bound variable onto the stack
| LOAD     |   D      | Load object from data segment onto the stack
| POP      |   D      | Pop D objects from the stack
| SLOAD    |   D      | Load element of stack with Base to stack slot
| SSTORE   |   D      | Store top of stack to local variable
| BINOP    |   W      | Perform a binary operation to the top 2 values on the stack
| UNIOP    |   W      | Perform a unary operation to the top value on the stack
| MAKELIST |   D      | Make an array object with the top D values on the stack
| MAKEHASH |   D      | Make a hash object with the top 2 * D values on the stack
| MAKEFUNC |   DD     | Make a function object with the top D values on the stack
| INDEX    |   NNN    | Get index item TOP from base object TOP-1
| JUMP     |   D      | Jump to the instruction at the given address
| JUMPFWD  |   D      | Jump forward D instructions
| JUMPIF   |   D      | Jump forward D instructions if the top value on the stack is false
| CALL     |   D      | Call the function with the top D values as arguments
| TAILCALL |   D      | Tail call the function with the top D values as arguments
| RETURN   |   D      | Return function call, pop D values from the stack as return values
| HALT     |   NNN    | Halt the VM
