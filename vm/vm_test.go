package vm

import (
	"strings"
	"testing"

	"github.com/flily/macaque-lang/compiler"
	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/opcode"
)

func testCompileCode(t *testing.T, code string) *opcode.CodePage {
	t.Helper()

	c := compiler.NewCompiler()
	block, err := c.CompileCode("testcase", []byte(code))
	if err != nil {
		t.Fatalf("compiler error:\n%s", err)
	}

	page := c.Link(block)
	return page
}

func checkVMStackTop(t *testing.T, name string, m VM, expecteds []object.Object) {
	t.Helper()

	if len(expecteds) == 0 {
		return
	}

	for i, expected := range expecteds {
		sp := m.GetSP()
		index := int(sp) - i - 1
		if index >= int(m.GetSP()) || index < 0 {
			t.Fatalf("[%s] ERROR on %d: stack do not have enough elements", name, i)
		}

		got := m.GetStackObject(index)
		if got.EqualTo(expected) == false {
			t.Fatalf("[%s] ERROR on stack %d: expect %s, got %s", name, i, expected, got)
		}
	}
}

func checkVMResult(t *testing.T, name string, m VM, result []object.Object, expected []object.Object) {
	t.Helper()

	// if len(result) != len(expected) {
	// 	t.Fatalf("[%s] ERROR: result length not matched, expect %d, got %d", name, len(expected), len(result))
	// }

	// for i, r := range result {
	// 	if r.EqualTo(expected[i]) == false {
	// 		t.Fatalf("[%s] ERROR on %d result: expect %s, got %s", name, i, expected[i], r)
	// 	}
	// }
}

type registerAssertion struct {
	register string
	value    uint64
}

type vmTest struct {
	code      string
	stack     []object.Object
	registers []registerAssertion
	// result    []object.Object
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

// func result(o ...object.Object) []object.Object {
// 	return o
// }

func sp(v uint64) registerAssertion {
	return registerAssertion{"sp", v}
}

func bp(v uint64) registerAssertion {
	return registerAssertion{"bp", v}
}

func checkVMRegisters(t *testing.T, name string, m VM, cases []registerAssertion) {
	t.Helper()

	for _, c := range cases {
		regValue := m.GetRegister(c.register)
		if regValue != c.value {
			t.Errorf("[%s] register %s error: expect %d, got %d", name, c.register, c.value, regValue)

			switch c.register {
			case "sp":
				vs, vv := m.InspectStack()
				t.Errorf("STACK: %s\n", vs)
				t.Errorf("       %s\n", vv)
			}
		}
	}
}

func runVMTestOnInstance(t *testing.T, name string, vm VM, c vmTest) {
	t.Helper()

	page := testCompileCode(t, c.code)
	main := page.Main().Func(nil)

	vm.LoadCodePage(page)
	_, err := vm.Run(main)
	if err != nil {
		t.Fatalf("%s error: %s", name, err)
	}

	checkVMStackTop(t, name, vm, c.stack)
	checkVMRegisters(t, name, vm, c.registers)
	// checkVMResult(t, name, vm, result, c.result)
	checkVMResult(t, name, vm, nil, nil)
}

func runVMTest(t *testing.T, cases []vmTest) {
	t.Helper()

	for _, c := range cases {
		runVMTestOnInstance(t, "vme", NewNaiveVM(), c)
		runVMTestOnInstance(t, "vmi", NewNaiveVMInterpreter(), c)
	}
}
