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

func makeCodePrint(left []opcode.Opcode, right []opcode.Opcode) string {
	length := len(left)
	if len(right) > length {
		length = len(right)
	}

	parts := make([]string, 0, length)

	for i := 0; i < length; i++ {
		l, r := "", ""
		if i < len(left) {
			l = left[i].String()
		}

		if i < len(right) {
			r = right[i].String()
		}

		flag := "      "
		if l != r {
			flag = ">>>>>>"
		}

		parts = append(parts, fmt.Sprintf("%2d %s   %-40s %s", i, flag, l, r))
		if l != r {
			parts = append(parts, strings.Repeat("-", 80))
		}
	}

	return strings.Join(parts, "\n")
}

func testCompileCode(t *testing.T, code string) (*Compiler, *CodePage, error) {
	t.Helper()

	scanner := lex.NewRecursiveScanner("testcase")
	_ = scanner.SetContent([]byte(code))
	parser := parser.NewLLParser(scanner)
	err := parser.ReadTokens()
	if err != nil {
		t.Fatalf("parser error:\n%s", err)
	}

	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("parser error:\n%s", err)
	}

	compiler := NewCompiler()
	page, err := compiler.Compile(program)
	return compiler, page, err
}

func checkInstructions(t *testing.T, text string, page *CodePage, compiler *Compiler, expecteds []opcode.Opcode) {
	t.Helper()

	codes := page.Codes
	if len(codes) != len(expecteds) {
		t.Errorf("wrong answer in code: %s", text)
		t.Errorf("wrong instructions length. want=%d, got=%d", len(expecteds), len(codes))
		t.Fatalf("want and got:\n%s", makeCodePrint(expecteds, codes))
	}

	for i, ins := range codes {
		if ins != expecteds[i] {
			t.Errorf("wrong answer in code: %s", text)
			t.Errorf("wrong instruction at %d. want=%q, got=%q", i, expecteds[i], ins)
			t.Fatalf("want and got:\n%s", makeCodePrint(expecteds, codes))
		}
	}
}

func checkData(t *testing.T, text string, compiler *Compiler, expecteds []object.Object) {
	t.Helper()

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
	t.Helper()

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
	t.Helper()

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
