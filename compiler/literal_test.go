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
				inst(opcode.IMakeFunc, 0, 0),
				inst(opcode.ISStore, 1),
				inst(opcode.IHalt),
				inst(opcode.ILoadInt, 42),
				inst(opcode.IReturn, 1),
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
				inst(opcode.IMakeFunc, 0, 1),
				inst(opcode.ISStore, 2),
				inst(opcode.ISLoad, 2),
				inst(opcode.IHalt),
				inst(opcode.ILoadBind, 0),
				inst(opcode.IReturn, 1),
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
				"}",
			),
			code(
				// let answer = 42
				inst(opcode.ILoadInt, 42), // 0
				inst(opcode.ISStore, 1),   // 1
				// let f = fn(a)
				inst(opcode.ISLoad, 1),       // 2
				inst(opcode.IMakeFunc, 0, 1), // 3
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
				inst(opcode.IReturn, 1),              // 31
				inst(opcode.IHalt),                   // 32
			),
			data(),
		},
	}

	runCompilerTestCases(t, tests)
}
