package vm

import (
	"strings"
	"testing"

	"github.com/flily/macaque-lang/compiler"
	"github.com/flily/macaque-lang/lex"
	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/parser"
)

func testCompileCode(t *testing.T, code string) *compiler.CodePage {
	t.Helper()

	scanner := lex.NewRecursiveScanner("testcase")
	scanner.SetContent([]byte(code))
	parser := parser.NewLLParser(scanner)
	err := parser.ReadTokens()
	if err != nil {
		t.Fatalf("parser error: %s", err)
	}

	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("parser error:\n%s", err)
	}

	compiler := compiler.NewCompiler()
	page, err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("compiler error:\n%s", err)
	}

	return page
}

func checkVMStackTop(t *testing.T, m *NaiveVM, expecteds []object.Object) {
	t.Helper()

	if len(expecteds) == 0 {
		return
	}

	for i, expected := range expecteds {
		index := int(m.sp) - i - 1
		if index >= int(m.sp) {
			t.Fatalf("ERROR on %d: stack do not have enough elements", i)
		}

		got := m.Stack[index]
		if got.EqualTo(expected) == false {
			t.Fatalf("ERROR on %d: expect %s, got %s", i, expected, got)
		}
	}
}

type registerAssertion struct {
	register string
	value    uint64
}

type vmTest struct {
	code      string
	stack     []object.Object
	registers []registerAssertion
}

func assertRegister(r ...registerAssertion) []registerAssertion {
	return r
}

func text(lines ...string) string {
	return strings.Join(lines, "\n")
}

func stack(o ...object.Object) []object.Object {
	return o
}

func sp(v uint64) registerAssertion {
	return registerAssertion{"sp", v}
}

func bp(v uint64) registerAssertion {
	return registerAssertion{"bp", v}
}

func runVMRegisterCheck(t *testing.T, m *NaiveVM, cases []registerAssertion) {
	t.Helper()

	for _, c := range cases {
		switch c.register {
		case "sp":
			if m.sp != c.value {
				t.Errorf("sp error: expect %d, got %d", c.value, m.sp)
			}

		case "bp":
			if m.bp != c.value {
				t.Errorf("bp error: expect %d, got %d", c.value, m.bp)
			}
		}
	}
}

func runVMTest(t *testing.T, cases []vmTest) {
	t.Helper()

	for _, c := range cases {
		m := NewNaiveVM()
		page := testCompileCode(t, c.code)
		m.LoadCodePage(page)
		main := page.Main().Func(nil)
		err := m.Run(main)
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		checkVMStackTop(t, m, c.stack)
		runVMRegisterCheck(t, m, c.registers)
	}
}
