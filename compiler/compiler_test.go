package compiler

import (
	"fmt"
	"testing"

	"strings"

	"github.com/flily/macaque-lang/lex"
	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/opcode"
	"github.com/flily/macaque-lang/parser"
)

func makeCodePrint(codes []opcode.Opcode) string {
	parts := make([]string, len(codes))
	for i, c := range codes {
		parts[i] = fmt.Sprintf("%2d %s", i, c.String())
	}

	return strings.Join(parts, "\n")
}

func testCompileCode(t *testing.T, code string) (*Compiler, error) {
	scanner := lex.NewRecursiveScanner("testcase")
	_ = scanner.SetContent([]byte(code))
	parser := parser.NewLLParser(scanner)
	err := parser.ReadTokens()
	if err != nil {
		t.Fatalf("parser error: %s", err)
	}

	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("parser error: %s", err)
	}

	compiler := NewCompiler()
	_, err = compiler.Compile(program)
	return compiler, err
}

func checkInstructions(t *testing.T, compiler *Compiler, expecteds []opcode.Opcode) {
	instructions := compiler.Context.Instructions
	if len(instructions) != len(expecteds) {
		t.Errorf("wrong instructions length. want=%d, got=%d", len(expecteds), len(instructions))
		t.Fatalf("want:\n%s\ngot:\n%s", makeCodePrint(expecteds), makeCodePrint(instructions))
	}

	for i, ins := range instructions {
		if ins != expecteds[i] {
			t.Errorf("wrong instruction at %d. want=%q, got=%q", i, expecteds[i], ins)
			t.Fatalf("want:\n%s\ngot:\n%s", makeCodePrint(expecteds), makeCodePrint(instructions))
		}
	}
}

func checkData(t *testing.T, compiler *Compiler, expecteds []object.Object) {
	data := compiler.Context.Literal.Values
	if len(data) != len(expecteds) {
		t.Errorf("wrong data length. want=%d, got=%d", len(expecteds), len(data))
	}

	for i, d := range data {
		if d.EqualTo(expecteds[i]) == false {
			t.Errorf("wrong data at %d. want=%q, got=%q", i, expecteds[i], d)
		}
	}
}

func inst(name int, ops ...int) opcode.Opcode {
	return opcode.Inst(name, ops...)
}

func code(codes ...opcode.Opcode) []opcode.Opcode {
	return codes
}

func data(o ...object.Object) []object.Object {
	return o
}

type testCompilerCase struct {
	code  string
	codes []opcode.Opcode
	data  []object.Object
}

func runCompilerTestCases(t *testing.T, cases []testCompilerCase) {
	for _, c := range cases {
		compiler, err := testCompileCode(t, c.code)
		if err != nil {
			t.Fatalf("compiler error:\n%s", err)
		}

		checkInstructions(t, compiler, c.codes)
		checkData(t, compiler, c.data)
	}
}

type testCompilerErrorCase struct {
	codes   []string
	message []string
}

func runCompilerErrorTestCases(t *testing.T, cases []testCompilerErrorCase) {
	for _, c := range cases {
		code := strings.Join(c.codes, "\n")
		_, err := testCompileCode(t, code)
		if err == nil {
			t.Fatalf("compilation should fail:\n")
		}

		expected := strings.Join(c.message, "\n")
		if err.Error() != expected {
			t.Fatalf("incorrect error message:\n%s\nexpect:\n%s", err, expected)
		}
	}
}
