Macaque Language Virtual Machine Specification
===============================================

Overview
---------

Macaque language uses a stack-based virtual machine, and 32-bit unsigned
integers bytecodes.


Bytecodes bit layout
---------------------

Every bytecode is 32-bit unsigned integer, actually can be process as 4 8-bit
integers, in big-endian byte order. The first byte is the opcode, and the rest
bytes are operands.

| BYTE 0 | BYTE 1 | BYTE 2 | BYTE 3 |
|--------|--------|--------|--------|
| OPCODE | OPERAND | OPERAND | OPERAND |

Bytecode can extend to more than 4 bytes, but it MUST be integer times of 4.


Memory layout
--------------

Macaque VM use 32-bit instruction, and 64-bit data element in memory.

Macaque VM has 5 memory segments:
  1. **Code segment**: contains all the bytecode of the program.
  2. **Data segment**: contains static data and constants.
  3. **Stack segment**: contains local variable data elements.
  4. 


Opcode References
------------------

Letters `NBWTbwt` are used to represent the type of operands. Upper case letter
represents an unsigned integer, and lower case represents a signed integer.
  + `N`: nil, ALL-ZERO byte.
  + `B`, `b`: byte, an 8-bit integer.
  + `W`, `w`: word, a 16-bit integer.
  + `T`, `t`: tri-byte word, a 24-bit integer.

The mnemonic of opcode MUST NOT be longer than 8 characters.

| Opcode (Hex) | Mnemonic | Operands | Description                 |
|:------------:|:--------:|:--------:|:----------------------------|
|    0x00      | NOP      |   NNN    | No operation
|    0x01      | JUMPA    |   T      | Jump to an absolute address
|    0x02      | JUMPR    |   t      | Jump to an relative address
|    0x10      | PUTI     |   t      | Put an integer from operand to stack
|    0x11      | PUTN     |   NNN    | Put a null value to stack
|    0xff      | HALT     |   NNN    | Halt the VM
