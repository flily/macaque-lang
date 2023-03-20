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


Missing features in Monkey, but not decided to add in Macaque yet:
  - Loop statement, like `while` and `for`, but it can be implemented by recursion.
    + When using recursion, tail call optimization is required.
  - Local and global variables.
    + For global variables is harmful, make all variables are local variables may be better.
    + Use naming convention to distinguish local and global variables, variables start with `_` is
      local.
  - Add support for object-like use of hash, with following features:
    + Reference item key with identifier, like `hash.key`;
    + Call function item directly with key reference, like `hash.key()` or `hask:key()` like lua.
