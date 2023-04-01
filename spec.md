Macaque Programming Language Specification
===========================================


Introdution
------------

Macaque programming language is a dialect of the [Monkey programming language](https://monkeylang.org/),
which is an example language for the book [Writing An Interpreter In Go](https://interpreterbook.com/)
and [Writing A Compiler In Go](https://compilerbook.com/) written by Thorsten Ball.

But a detailed specification of the Monkey programming language is missing, and the implementation
on the book is not complete and clean. So I decided to extend the Monkey programming language, as
the Macaque programming language, and write a complete specification of it.


Requirement and Notation
-------------------------

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT",
"RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in
[RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

The syntax is specified using BNF defined in [RFC 5234](https://www.rfc-editor.org/rfc/rfc5234).
But special, all characters are case-sensitive.


Source code encoding
---------------------

Lexical elements
-----------------

Types
------

Macaque has 9 basic data types:
  - `null`: the only value of this type is `null`.
  - `boolean`: the only values of this type are `true` and `false`.
  - `integer`: a signed 64-bit integer.
  - `float`: a double-precision floating point number.
  - `string`: a sequence of characters, encoded in raw.
  - `array`: a sequence of values.
  - `hash`: a sequence of key-value pairs.
  - `function`: a function.
  - `user value`: a value defined by the user.


Expressions
------------

### Operands

### Index expression

### Function literal

### Function call expression


Statements
-----------

### Let statement

### Assignment

### Blocks

### Return statement

### Import statement


Packages
---------

References
-----------


Complete syntax of Macaque
---------------------------
```abnf
program = *statement
statement = let-stmt
          / return-stmt
          / import-stmt       ; not determined yet
          / if-stmt
          / block-stmt
          / expression-stmt

let-stmt = "let" identifier-list "=" expression-list ";"

return-stmt = "return" [expression-list] ";"

import-stmt = "import" string-literal ";"  ; not determined yet

expression-stmt = expression-list ";"

expression-list = expression *( "," expression ) [","]

expression = literals
           / identifier
           / prefix-expression
           / infix-expression
           / index-expression
           / call-expression

literals = null-literal
         / boolean-literal
         / integer-literal
         / float-literal
         / string-literal
         / array-literal
         / hash-literal
         / function-literal

null-literal = "null"

boolean-literal = "true" / "false"

integer-literal = ( ["-"] 1* (DIGIT / "_") ) / ( "0x" 1*( HEXDIG / "_" ) )

float-literal = 1*( DIGIT / "_" ) "." 1*( DIGIT / "_" )

DQUOTE = %x22
        ; " (Double Quote)
        ; predefined name in RFC 5234

string-literal = DQUOTE *string-chars DQUOTE

string-chars = %x20-21 / %x23-5B / %x5D-10FFFF
               ; any Unicode character except double quote (") and backslash (\)
             / escape-sequence

escape-sequence = "\" ( DQUOTE / "\" / "n" / "r" / "t" / "x" *2HEXDIG )

array-literal = "[" expression-list "]"

hash-literal = "{" hash-pair *( "," hash-pair ) [","] "}"

hash-pair = expression ":" expression

function-literal = "fn" "(" [identifier-list] ")" block-stmt

identifier-list = identifier *( "," identifier ) [","]

identifier = [identifier-prefix] ( ALPHA / "_" ) *( ALPHA / DIGIT / "_" ) [identifier-suffix]

identifier-prefix = "@" / "$"

identifier-suffix = "!" / "?"

index-expression = ( expression "[" expression "]" )
                 / ( expression "." identifier )

call-expression = expression [ ":" identifier ] "(" [expression-list] ")"

ALPHA = %x41-5A / %x61-7A  
        ; A-Z / a-z
        ; predefined name in RFC 5234

DIGIT = %x30-39
        ; 0-9
        ; predefined name in RFC 5234

```