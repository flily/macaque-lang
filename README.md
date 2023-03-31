macaque-lang
=============

A dialect of monkey language, with brand new implement from the book and feature improvements.


Specification
--------------

Specification of the Macaque programming language is in the [spec.md](spec.md) file


Improvements from Monkey
-------------------------

Macaque makes some improvements from Monkey:
  - Add `null` type and keyword.
  - Add `float` type for calculation.
  - Add `import` statement, to write more complex program with multiple modules files.
    - Use `return` statement in out most scope to export module, just like lua.
  - Remove built-in functions, but introduct standard library instead.
  - Support regular format `[a-zA-Z_][a-zA-Z0-9_]*` for identifiers.
  - Implement more readable error and warning messages.
  - Support `else if` statement.


Missing features in Monkey, but not decided to add in Macaque yet:
  - Loop statement, like `while` and `for`, but it can be implemented by recursion.
    + When using recursion, tail call optimization is required.
    + Utility functions like `first`, `rest`, `last` are required.
    + Slice of array and hash is required to optimize performance.
  - Local and global variables.
    + For global variables is harmful, make all variables are local variables may be better.
    + Use naming convention to distinguish local and global variables, variables start with `_` is
      local.
  - Add support for object-like use of hash, with following features:
    + Reference item key with identifier, like `hash.key`;
    + Call function item directly with key reference, like `hash.key()` or `hask:key()` like lua.
  - Error handling mechanism.
    + Use `try`, `catch`, `finally` and `throw` like Java.
    + Use `ON ERROR` trap like BASIC.
    + Use `pcall` and `xpcall` like lua.
    + Use `recover()` with `defer` like go, but it sucks.
    + Directly return error by every function call, like go and lua.
      * Can not handle critical errors like panic.
    + Use `Result` type to wrap return value, like rust, but it is not elegant.
      * And I don't want to implement generics type system.
      * It is not elegant that a language behaviour is highly dependent on a specific type.
  - Get variable type.
    + Use `typeof()` or `type()` to get type of variable.
    + Type representation.
      * Use type type like `int`, `string` to represent type.
      * Use string like `"INT"`, `"STRING"` to represent type, like lua.
      * Use type variable like `Int` or `std.Int` to represent type, like Java.
  - `return` statement can return multiple values.
  - Detect variable redeclaration.
    + May be not a bug. Perhaps the author use let statement to modify variable.
      * But, the code `let a = 1; let a = a + 1;` may crush in compiler but not in interpreter.
    + Fix it and make variable immutable.
  - Make variable can be modified after declaration.
    + Use `var` keyword to declare variable and `let` keyword to declare immutable variable, like
      rust. The keyword `mut` used in rust is not elegant.
    + Make all variables immutable, like erlang.
      * Some mechanism like pattern matching may be required.
  - Add variable parameter list.
    + Use `...` to represent variable parameter list, like lua.
    + Use `*` to represent variable parameter list, like python.
  - Add debuggging support.
    + Support step trace debugging and breakpoint.
  - Array and hash modification.
    + In offical implement, monkey-lang can ONLY modify array, append element to the end, via
      builtin function `push`. And there is no way to modify hash.
    + Write more external builtin functions to modify array and hash, but it is not elegant.
    + Write native monkey-lang code to modify, which in the way like erlang, build a new hash or
      array in functional programming way.
