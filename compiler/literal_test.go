package compiler

import (
	"testing"

	"github.com/flily/macaque-lang/opcode"
	"github.com/flily/macaque-lang/token"
)

func TestCompileListLiteral(t *testing.T) {
	tests := []testCompilerCase{
		{
			text(
				"[]",
			),
			code(
				inst(opcode.IMakeList, 0),
			),
			data(),
		},
		{
			text(
				"[1, 2, 3, 4]",
			),
			code(
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 3),
				inst(opcode.ILoadInt, 4),
				inst(opcode.IMakeList, 4),
			),
			data(),
		},
		{
			text(
				"[1, 2 + 3, 4,]",
			),
			code(
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 3),
				inst(opcode.IBinOp, int(token.Plus)),
				inst(opcode.ILoadInt, 4),
				inst(opcode.IMakeList, 3),
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}

func TestCompileHashLiteral(t *testing.T) {
	tests := []testCompilerCase{
		{
			text(
				"{}",
			),
			code(
				inst(opcode.IMakeHash, 0),
			),
			data(),
		},
		{
			text(
				"{1: 2, 3: 4}",
			),
			code(
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 3),
				inst(opcode.ILoadInt, 4),
				inst(opcode.IMakeHash, 2),
			),
			data(),
		},
		{
			text(
				"{1: 2 + 3, 4: 5 + 6}",
			),
			code(
				inst(opcode.ILoadInt, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ILoadInt, 3),
				inst(opcode.IBinOp, int(token.Plus)),
				inst(opcode.ILoadInt, 4),
				inst(opcode.ILoadInt, 5),
				inst(opcode.ILoadInt, 6),
				inst(opcode.IBinOp, int(token.Plus)),
				inst(opcode.IMakeHash, 2),
			),
			data(),
		},
		{
			text(
				"let one = 1;",
				"let two = 2;",
				"{one: 10, 2: two}",
			),
			code(
				inst(opcode.ILoadInt, 1),
				inst(opcode.ISStore, 1),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ISStore, 2),
				inst(opcode.ISLoad, 1),
				inst(opcode.ILoadInt, 10),
				inst(opcode.ILoadInt, 2),
				inst(opcode.ISLoad, 2),
				inst(opcode.IMakeHash, 2),
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}

func TestCompileFunctions(t *testing.T) {
	tests := []testCompilerCase{
		{
			text(
				"let answer = fn() { 42; };",
			),
			code(
				inst(opcode.IMakeFunc, 1, 0),
				inst(opcode.ISStore, 1),
				inst(opcode.IHalt),
				inst(opcode.ILoadInt, 42),
				inst(opcode.IReturn),
				inst(opcode.IHalt),
			),
			data(),
		},
		{
			text(
				"let answer = 42;",
				"let final_answer = fn() { answer; };",
				"final_answer",
			),
			code(
				inst(opcode.ILoadInt, 42),
				inst(opcode.ISStore, 1),
				inst(opcode.ISLoad, 1),
				inst(opcode.IMakeFunc, 1, 1),
				inst(opcode.ISStore, 2),
				inst(opcode.ISLoad, 2),
				inst(opcode.IHalt),
				inst(opcode.ILoadBind, 0),
				inst(opcode.IReturn),
				inst(opcode.IHalt),
			),
			data(),
		},
		{
			text(
				"let answer = 42;",
				"let f = fn(a) {",
				"    let b = 3;",
				"    if (a > b) {",
				"        let c = a + b;",
				"        c",
				"    } else {",
				"        let d = a - b;",
				"        let c = a + b;",
				"        d + c + answer",
				"    }",
				"};",
			),
			code(
				// let answer = 42
				inst(opcode.ILoadInt, 42), // 0
				inst(opcode.ISStore, 1),   // 1
				// let f = fn(a)
				inst(opcode.ISLoad, 1),       // 2
				inst(opcode.IMakeFunc, 1, 1), // 3
				inst(opcode.ISStore, 2),      // 4
				inst(opcode.IHalt),           // 5

				// fn(a) {
				//     let b = 3;
				inst(opcode.ILoadInt, 3), // 6
				inst(opcode.ISStore, 1),  // 7
				//     if (a > b) {
				inst(opcode.ISLoad, -1),            // 8
				inst(opcode.ISLoad, 1),             // 9
				inst(opcode.IBinOp, int(token.GT)), // 10
				inst(opcode.IJumpIf, 6),            // 11
				//         let c = a + b;
				inst(opcode.ISLoad, -1),              // 12
				inst(opcode.ISLoad, 1),               // 13
				inst(opcode.IBinOp, int(token.Plus)), // 14
				inst(opcode.ISStore, 2),              // 15
				//         c
				inst(opcode.ISLoad, 2),    // 16
				inst(opcode.IJumpFWD, 13), // 17
				//    } else {
				//         let d = a - b;
				inst(opcode.ISLoad, -1),               // 18
				inst(opcode.ISLoad, 1),                // 19
				inst(opcode.IBinOp, int(token.Minus)), // 20
				//         let e = a + b;
				inst(opcode.ISStore, 2),              // 21
				inst(opcode.ISLoad, -1),              // 22
				inst(opcode.ISLoad, 1),               // 23
				inst(opcode.IBinOp, int(token.Plus)), // 24
				//         d + e + answer
				inst(opcode.ISStore, 3),              // 25
				inst(opcode.ISLoad, 2),               // 26
				inst(opcode.ISLoad, 3),               // 27
				inst(opcode.IBinOp, int(token.Plus)), // 28
				inst(opcode.ILoadBind, 0),            // 29
				inst(opcode.IBinOp, int(token.Plus)), // 30
				inst(opcode.IReturn),                 // 31
				inst(opcode.IHalt),                   // 32
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}

func TestCompileComplexFunctions(t *testing.T) {
	tests := []testCompilerCase{
		{
			text(
				"let a = 1;",                            // 1
				"let newAdderOuter = fn(b) {",           // 2
				"    fn(c) {",                           // 3
				"        fn(d) { a + b + c + d };",      // 4
				"    };",                                // 5
				"};",                                    // 6
				"let newAdderInner = newAdderOuter(2);", // 7
				"let adder = newAdderInner(3);",         // 8
				"adder(8);",                             // 9
			),
			code(
				// let a = 1
				inst(opcode.ILoadInt, 1), // 0
				inst(opcode.ISStore, 1),  // 1
				// let newAdderOuter = fn(b)
				inst(opcode.ISLoad, 1),       // 2
				inst(opcode.IMakeFunc, 3, 1), // 3
				inst(opcode.ISStore, 2),      // 4
				// let newAdderInner = newAdderOuter(2)
				inst(opcode.ILoadInt, 2), // 5
				inst(opcode.ISLoad, 2),   // 6
				inst(opcode.ICall, 1),    // 7
				inst(opcode.ISStore, 3),  // 8
				// let adder = newAdderInner(3)
				inst(opcode.ILoadInt, 3), // 9
				inst(opcode.ISLoad, 3),   // 10
				inst(opcode.ICall, 1),    // 11
				inst(opcode.ISStore, 4),  // 12
				// adder(8)
				inst(opcode.ILoadInt, 8), // 13
				inst(opcode.ISLoad, 4),   // 14
				inst(opcode.ICall, 1),    // 15
				inst(opcode.IHalt),       // 16
				// fn(d) {
				//     a + b + c + d }
				inst(opcode.ILoadBind, 0),            // 17
				inst(opcode.ILoadBind, 1),            // 18
				inst(opcode.IBinOp, int(token.Plus)), // 19
				inst(opcode.ILoadBind, 2),            // 20
				inst(opcode.IBinOp, int(token.Plus)), // 21
				inst(opcode.ISLoad, -1),              // 22
				inst(opcode.IBinOp, int(token.Plus)), // 23
				inst(opcode.IReturn),                 // 24
				inst(opcode.IHalt),                   // 25
				// fn(c) {
				//     fn(d) { a + b + c + d } }
				inst(opcode.ILoadBind, 0),    // 26
				inst(opcode.ILoadBind, 1),    // 27
				inst(opcode.ISLoad, -1),      // 28
				inst(opcode.IMakeFunc, 1, 3), // 29
				inst(opcode.IReturn),         // 30
				inst(opcode.IHalt),           // 31
				// fn(b) {
				//     fn(c) { fn(d) { a + b + c + d } } }
				inst(opcode.ILoadBind, 0),    // 32
				inst(opcode.ISLoad, -1),      // 33
				inst(opcode.IMakeFunc, 2, 2), // 34
				inst(opcode.IReturn),         // 35
				inst(opcode.IHalt),           // 36
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}

func TestCompileRecursiveFunctions(t *testing.T) {
	tests := []testCompilerCase{
		{
			text(
				"let countDown = fn(x) {",
				"    if (x == 0) {",
				"        return 0;",
				"    } else {",
				"        fn(x - 1);",
				"    }",
				"};",
				"countDown(1);",
			),
			code(
				// let countDown = fn(x)
				inst(opcode.IMakeFunc, 1, 0), // 0
				inst(opcode.ISStore, 1),      // 1
				// countDown(1)
				inst(opcode.ILoadInt, 1), // 2
				inst(opcode.ISLoad, 1),   // 3
				inst(opcode.ICall, 1),    // 4
				inst(opcode.IHalt),       // 5
				// fn(x) {
				//     if (x == 0) {
				inst(opcode.ISLoad, -1),            // 6
				inst(opcode.ILoadInt, 0),           // 7
				inst(opcode.IBinOp, int(token.EQ)), // 8
				inst(opcode.IJumpIf, 3),            // 9
				//         return 0;
				inst(opcode.ILoadInt, 0), // 10
				inst(opcode.IReturn),     // 11
				inst(opcode.IJumpFWD, 5), // 12
				//     } else {
				//         fn(x - 1);
				inst(opcode.ISLoad, -1),               // 13
				inst(opcode.ILoadInt, 1),              // 14
				inst(opcode.IBinOp, int(token.Minus)), // 15
				inst(opcode.ISLoad, 0),                // 16
				inst(opcode.ICall, 1),                 // 17
				//    }
				inst(opcode.IReturn), // 18
				inst(opcode.IHalt),   // 19
			),
			data(),
		},
		{
			text(
				"let countDown = fn(x) {",
				"    if (x == 0) {",
				"        return 0;",
				"    } else {",
				"        fn(x - 1);",
				"    }",
				"};",
				"let wrapper = fn() {",
				"    countDown(1);",
				"};",
				"wrapper();",
			),
			code(
				// let countDown = fn(x)
				inst(opcode.IMakeFunc, 1, 0), // 0
				inst(opcode.ISStore, 1),      // 1
				// let wrapper = fn()
				inst(opcode.ISLoad, 1),       // 2
				inst(opcode.IMakeFunc, 2, 1), // 3
				inst(opcode.ISStore, 2),      // 4
				// wrapper()
				inst(opcode.ISLoad, 2), // 5
				inst(opcode.ICall, 0),  // 6
				inst(opcode.IHalt),     // 7
				// fn(x) {
				//     if (x == 0) {
				inst(opcode.ISLoad, -1),            // 8
				inst(opcode.ILoadInt, 0),           // 9
				inst(opcode.IBinOp, int(token.EQ)), // 10
				inst(opcode.IJumpIf, 3),            // 11
				//         return 0;
				inst(opcode.ILoadInt, 0), // 12
				inst(opcode.IReturn),     // 13
				inst(opcode.IJumpFWD, 5), // 14
				//     } else {
				//         fn(x - 1);
				inst(opcode.ISLoad, -1),               // 15
				inst(opcode.ILoadInt, 1),              // 16
				inst(opcode.IBinOp, int(token.Minus)), // 17
				inst(opcode.ISLoad, 0),                // 18
				inst(opcode.ICall, 1),                 // 19
				//    }
				inst(opcode.IReturn), // 20
				inst(opcode.IHalt),   // 21
				// fn()
				//    countDown(1);
				inst(opcode.ILoadInt, 1),  // 22
				inst(opcode.ILoadBind, 0), // 23
				inst(opcode.ICall, 1),     // 24
				inst(opcode.IReturn),      // 25
				inst(opcode.IHalt),        // 26
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}
