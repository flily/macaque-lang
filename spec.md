Macaque Programming Language Specification
===========================================


Introdution
------------

Macaque programming language is a dialect of the
[Monkey programming language](https://monkeylang.org/), which is an example
language for the book
[Writing An Interpreter In Go](https://interpreterbook.com/) and
[Writing A Compiler In Go](https://compilerbook.com/) written by Thorsten Ball.

But a detailed specification of the Monkey programming language is missing, and
the implementation on the book is not complete and clean. So I decided to extend
the Monkey programming language, as the Macaque programming language, and write
a complete specification of it.


Requirement and Notation
-------------------------

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be
interpreted as described in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

The syntax is specified using BNF defined in
[RFC 5234](https://www.rfc-editor.org/rfc/rfc5234).
But special, all characters are case-sensitive.


Source code encoding
---------------------

While macaque language should be a lightweight scripting language. All lexical
elements of language are is ASCII. But, actually the parser and compiler are
written in Go, which process the source code in UTF-8, it is OK to write UNICODE
characters in source code, where it is allowed. ONLY ASCII characters are
accepted for all language keywords, operators and identifiers.

The following positions are allowed to write UNICODE characters:
  - Comments. Comment content can be composed with any UNICODE characters,
    including emojis. But the comment commands, if one day I decide to implement
    it, MUST be ASCII characters.
  - String literals. Macaque language supports ONLY raw string bytes, which they
    will be recognized as is.

To avoid mojibake (which called 乱码 in Chinese), you
**MUST use ONLY ASCII or UTF-8** to encode the source code files.


Lexical elements
-----------------

### Comments

Macaque language DO support comments, which the original Monkey language does
not. Use `//` to start a single-line comment, and stop at the end of the line.
Multiple-line comments are not determined to design yet.

### Keywords

Macauqe has 9 keywords:
  - `let`: declare a symbol to represent a value.
  - `fn`: start a function, or a lambda, literal.
  - `return`: return a value from a function.
  - `if` and `else`: basic control flow.
  - `import`: import a module from a file.
  - `null`: a special value that represents nothing.
  - `true` and `false`: boolean values literals.

### Operators and punctuation

#### In original Monkey language
Unary operators:
  - `!`: bang operator, logical NOT.
  - `-`: nagative operator, arithmetic negation.

Binary operators:
  - `+`, `-`, `*`, `/`: arithmetic operators.
  - `==`, `!=`, `<`, `>`: comparison operators.

Other punctuation:
  - `(`, `)`: parentheses, use for function calls and grouping.
  - `{`, `}`: curly braces, use for blocks and hash literals.
  - `[`, `]`: square brackets, use for indexing and array literals.
  - `=`: assignment, which is only used in `let` statement.
  - `;`: statement terminator.
  - `,`: delimiter for list of values in parameters list in function calls,
    elements of array and hash literals.
  - `:`: delimiter for key-value pairs in hash literals.
 
#### New in Macaque language
Following operators and punctuation are new introduced in Macaque language, to
make them easier to understand, I choose C-style operators and punctuation.
  - `&&`, `||`: logical AND and OR.
  - `~`, `&`, `|`, `^`: bitwise NOT, AND, OR and XOR.
  - `%`: modulus.
  - `<=`, `>=`: less than or equal to, greater than or equal to.
  - '.': access member of a hash.

### Keyword literals

`null` is a value representing nothing or empty or any other invalid value.
`true` and `false` are boolean values.

### Number literals

The original Monkey language only supports integer literals, which are 64-bit
signed integers. For read world usage, Macaque language add support for floating
numbers, which are double-precision floating point numbers, and integer literals
in hexadecimal format, like `0xDEADBEEF`.

In very large numbers, you can use `_` to separate the digits, like `1_000_000`
and `3.14_15_92_65_36`.

### Identifiers

The original Monkey language only supports alphabetic characters in identifiers,
and numbers are not accepted. For example, `answer42` is not a valid identifier
in Monkey language. But in Macaque language, numbers and underscore `_` are
accepted in identifiers, just like most other programming languages.

Even though Macaque language supports UTF-8 encoding for source code,
**UNICODE characters in identifiers ARE NOT ACCEPTED**.

### String literals

The original Monkey language does not describe strings in detail. In Macaque
language, strings are sequences of bytes, which are encoded in raw, or we can
say they have no encoding. And the **strings are immutable object**.

C-style escape sequences are supported in string literals. The following escape
sequences can be used:
  - `"\""`: double quote.
  - `"\\"`: backslash itself.
  - `"\n"`: newline.
  - `"\r"`: carriage return.
  - `"\t"`: tab character.
  - `"\xHH"`: hexadecimal byte, where `HH` is a hexadecimal number,
    between `00` and `FF`.

Types
------

Macaque has 9 basic data types:
  - `null`: the only value of this type is `null`, represent for nothing.
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

expression-stmt = expression ";"

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

integer-literal = decimal-integer / hexdecimal-integer

decimal-integer = NONZERO-DIGIT *( DIGIT / "_" )

hexdecimal-integer = "0x" 1*( HEXDIGIT / "_" )

float-literal = DIGIT *( DIGIT / "_" ) "." DIGIT *( DIGIT / "_" )

DQUOTE = %x22
        ; " (Double Quote)
        ; predefined name in RFC 5234

string-literal = DQUOTE *string-chars DQUOTE

string-chars = %x20-21 / %x23-5B / %x5D-10FFFF
               ; any Unicode character except double quote (") and backslash (\)
             / escape-sequence

escape-sequence = "\" ( DQUOTE / "\" / "n" / "r" / "t" / "x" *2HEXDIGIT )

array-literal = "[" expression-list "]"

hash-literal = "{" hash-pair *( "," hash-pair ) [","] "}"

hash-pair = expression ":" expression

function-literal = "fn" "(" [identifier-list] ")" block-stmt

identifier-list = identifier *( "," identifier ) [","]

identifier = [identifier-prefix] ( ALPHA / "_" ) *( ALPHA / DIGIT / "_" ) [identifier-suffix]

identifier-prefix = "@" / "$"  ; not implemented for now

identifier-suffix = "!" / "?"  ; not implemented for now

index-expression = ( expression "[" expression "]" )
                 / ( expression "." identifier )

call-expression = expression [ ":" identifier ] "(" [expression-list] ")"

prefix-expression = prefix-operator expression

prefix-operator = "!" / "-"    ; supported in official Monkey implementation
                / "~"          ; add to support bitwise NOT.

infix-expression = expression infix-operator expression

infix-operator = "+" / "-" / "*" / "/" / "==" / "!=" / "<" / ">"
                        ; supported in official Monkey implementation
               / "%"    ; modulo
               / "<="   ; less than or equal to
               / ">="   ; greater than or equal to
               / "&&"   ; logical AND
               / "||"   ; logical OR
               / "&"    ; bitwise AND
               / "|"    ; bitwise OR
               / "^"    ; bitwise XOR

if-stmt = "if" "(" expression ")" block-stmt [ "else" ( block-stmt / if-stmt ) ]
          ; support 'else if' statement

block-stmt = "{" *statement "}"

ALPHA = %x41-5A / %x61-7A  
        ; A-Z / a-z
        ; predefined name in RFC 5234

DIGIT = %x30-39
        ; 0-9
        ; predefined name in RFC 5234

HEXDIGIT = DIGIT / "A" / "B" / "C" / "D" / "E" / "F"
                 / "a" / "b" / "c" / "d" / "e" / "f"

NONZERO-DIGIT = %x31-39
                ; 1-9

```