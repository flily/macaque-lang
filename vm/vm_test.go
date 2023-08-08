package vm

import (
	"strings"
	"testing"

	"github.com/flily/macaque-lang/compiler"
	"github.com/flily/macaque-lang/lex"
	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/opcode"
	"github.com/flily/macaque-lang/parser"
)

func testCompileCode(t *testing.T, code string) *opcode.CodePage {
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
	block, err := compiler.CompileCode(program)
	if err != nil {
		t.Fatalf("compiler error:\n%s", err)
	}

	page := compiler.Link(block)
	return page
}

func checkVMStackTop(t *testing.T, m VM, expecteds []object.Object) {
	t.Helper()

	if len(expecteds) == 0 {
		return
	}

	for i, expected := range expecteds {
		sp := m.GetSP()
		index := int(sp) - i - 1
		if index >= int(m.GetSP()) || index < 0 {
			t.Fatalf("ERROR on %d: stack do not have enough elements", i)
		}

		got := m.GetStackObject(index)
		if got.EqualTo(expected) == false {
			t.Fatalf("ERROR on stack %d: expect %s, got %s", i, expected, got)
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

func runVMRegisterCheck(t *testing.T, m VM, cases []registerAssertion) {
	t.Helper()

	for _, c := range cases {
		regValue := m.GetRegister(c.register)
		if regValue != c.value {
			t.Errorf("register %s error: expect %d, got %d", c.register, c.value, regValue)
		}
	}
}

func runVMTestOnInstance(t *testing.T, name string, vm VM, c vmTest) {
	t.Helper()

	page := testCompileCode(t, c.code)
	main := page.Main().Func(nil)

	vm.LoadCodePage(page)
	err := vm.Run(main)
	if err != nil {
		t.Fatalf("%s error: %s", name, err)
	}

	checkVMStackTop(t, vm, c.stack)
	runVMRegisterCheck(t, vm, c.registers)
}

func runVMTest(t *testing.T, cases []vmTest) {
	t.Helper()

	for _, c := range cases {
		runVMTestOnInstance(t, "vme", NewNaiveVM(), c)
		runVMTestOnInstance(t, "vmi", NewNaiveVMInterpreter(), c)
	}
}
