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

func testCompileCode(t *testing.T, code string) (*Compiler, *CodePage, error) {
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
	page, err := compiler.Compile(program)
	return compiler, page, err
}

func checkInstructions(t *testing.T, text string, page *CodePage, compiler *Compiler, expecteds []opcode.Opcode) {
	codes := page.Codes
	if len(codes) != len(expecteds) {
		t.Errorf("wrong answer in code: %s", text)
		t.Errorf("wrong instructions length. want=%d, got=%d", len(expecteds), len(codes))
		t.Fatalf("want:\n%s\ngot:\n%s", makeCodePrint(expecteds), makeCodePrint(codes))
	}

	for i, ins := range codes {
		if ins != expecteds[i] {
			t.Errorf("wrong answer in code: %s", text)
			t.Errorf("wrong instruction at %d. want=%q, got=%q", i, expecteds[i], ins)
			t.Fatalf("want:\n%s\ngot:\n%s", makeCodePrint(expecteds), makeCodePrint(codes))
		}
	}
}

func checkData(t *testing.T, text string, compiler *Compiler, expecteds []object.Object) {
	data := compiler.Context.Literal.Values
	if len(data) != len(expecteds) {
		t.Errorf("wrong answer in code: %s", text)
		t.Errorf("wrong data length. want=%d, got=%d", len(expecteds), len(data))
	}

	for i, d := range data {
		if d.EqualTo(expecteds[i]) == false {
			t.Errorf("wrong answer in code: %s", text)
			t.Errorf("wrong data at %d. want=%q, got=%q", i, expecteds[i], d)
		}
	}
}

func inst(name int, ops ...int) opcode.Opcode {
	return opcode.Code(name, ops...)
}

func text(lines ...string) string {
	return strings.Join(lines, "\n")
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
		code := c.code
		compiler, page, err := testCompileCode(t, c.code)
		if err != nil {
			t.Fatalf("compiler error:\n%s", err)
		}

		checkInstructions(t, code, page, compiler, c.codes)
		checkData(t, code, compiler, c.data)
	}
}

type testCompilerErrorCase struct {
	code     string
	expected string
}

func runCompilerErrorTestCases(t *testing.T, cases []testCompilerErrorCase) {
	for _, c := range cases {
		_, _, err := testCompileCode(t, c.code)
		if err == nil {
			t.Fatalf("compilation should fail:\n")
		}

		got := err.Error()
		if got != c.expected {
			t.Fatalf("incorrect error message:\n%s\nexpect:\n%s", err, c.expected)
		}
	}
}
