macaque-lang
=============

![GitHub](https://img.shields.io/github/license/flily/macaque-lang)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/flily/macaque-lang)
![GitHub top language](https://img.shields.io/github/languages/top/flily/macaque-lang)

[![CI](https://github.com/flily/macaque-lang/actions/workflows/ci.yaml/badge.svg)](https://github.com/flily/macaque-lang/actions/workflows/ci.yaml)
[![Build Status](https://dev.azure.com/flily/macaque-lang/_apis/build/status%2Fflily.macaque-lang?branchName=main)](https://dev.azure.com/flily/macaque-lang/_build/latest?definitionId=1&branchName=main)
[![codecov](https://codecov.io/gh/flily/macaque-lang/branch/main/graph/badge.svg?token=DzOEyayucW)](https://codecov.io/gh/flily/macaque-lang)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/87a4caf4cbfd403fb143db2c501dba90)](https://app.codacy.com/gh/flily/macaque-lang/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/flily/macaque-lang)](https://goreportcard.com/report/github.com/flily/macaque-lang)


A dialect of monkey language, with brand new implement from the book and feature improvements.


Specification
--------------

Specification of the Macaque programming language is in the [spec.md](spec.md) file


Improvements from Monkey    
-------------------------

Principles to design:
  - It must be familiar to original Monkey language.
  - Only add necessary keywords.
  - Main purpose is to be a lightweight scripting language.
  - Try to add some modern features without too big changes.

Macaque makes some improvements from Monkey:
  - More human readable error messages.
  - Add `null` type and keyword.
  - Add `float` type for calculation.
  - Add `import` statement, to write more complex program with multiple modules files.
    + Use `return` statement in out most scope to export module, just like lua.
    + A implicit `return` will be added in out most scope if it is not exist.
  - Remove built-in functions, but introduct standard library instead.
  - Support regular format `[a-zA-Z_][a-zA-Z0-9_]*` for identifiers.
  - Implement more readable error and warning messages.
  - Support `else if` statement.
  - Fully support of c-style operators.
    + add `%` for modulo.
    + add `>=` and `<=` for comparison.
    + add `&&` and `||` for logical AND and OR.
    + use `!` for logical NOT and `~` for bitwise NOT.
    + add `&`, `|` and `^` for bitwise AND, OR and XOR.
  - Add one-line comments, leading with `//`.
  - Recursive call in closure by call `fn(a, b, c)` without function body.
    + Recursive call with function name is deprecated, because the original design makes bugs.
  - Add support for object-like use of hash, with following features:
    + Reference item key with identifier, like `hash.key`;
    + Use `hash::key()` to call a member function of a object, like this-call, with the object
      itself as the first parameter.
  - `return` statement can return multiple values.
  - `return` statement will act differently in following situations:
    + Works like lua.
    + In top-level scope, it will return the value as a module for importing.
    + In function body, it will return the value as the result of the function, and terminate the
      execution of rest code in the function.
  - The top-level statements are considered as the main function, and always have a return value.

Missing features in Monkey, but not decided to add in Macaque yet:
  - Loop statement, like `while` and `for`, but it can be implemented by recursion.
    + When using recursion, tail call optimization is required.
    + Utility functions like `first`, `rest`, `last` are required.
    + Slice of array and hash is required to optimize performance.
  - Local and global variables.
    + For global variables is harmful, make all variables are local variables may be better.
    + Use naming convention to distinguish local and global variables, variables start with `_` is
      local.
    + Add module level variables, for imported modules.
  - Make all types object.
    + Use `::function()` to make class-call, like lua.
    + `int` is object, has native methods and can be called on literals, `5::times()` like ruby.
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
  - Strings are raw strings, binary data. Do not support unicodes.
    + An unicode support library may be introduced.
    + Unicode string can be processed as array of integers.
