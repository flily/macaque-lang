package vm

import (
	"fmt"
	"testing"

	"github.com/flily/macaque-lang/object"
)

type monkeyTestCase struct {
	input    string
	expected interface{}
}

func monkeyExpectedValueCompare(t *testing.T, expected interface{}, got object.Object) bool {
	t.Helper()

	if expected == nil {
		return got.EqualTo(object.NewNull())
	}

	ok := false

TypeSwitch:
	switch e := expected.(type) {
	case int:
		if got.Type() != object.ObjectTypeInteger {
			t.Errorf("expect integer, got %s", got.Type())
			break TypeSwitch
		}

		ok = int64(e) == got.(*object.IntegerObject).Value

	case bool:
		if got.Type() != object.ObjectTypeBoolean {
			t.Errorf("expect boolean, got %s", got.Type())
			break TypeSwitch
		}

		ok = bool(e) == got.(*object.BooleanObject).Value

	case string:
		if got.Type() != object.ObjectTypeString {
			t.Errorf("expect string, got %s", got.Type())
			break TypeSwitch
		}

		ok = string(e) == got.(*object.StringObject).Value

	case []int:
		if got.Type() != object.ObjectTypeArray {
			t.Errorf("expect array, got %s", got.Type())
			break TypeSwitch
		}

		array := got.(*object.ArrayObject).Elements
		if len(array) != len(e) {
			t.Errorf("expect array length %d, got %d", len(e), len(array))
			break TypeSwitch
		}

		for i, v := range e {
			if !monkeyExpectedValueCompare(t, v, array[i]) {
				t.Errorf("expect array value at [%d] %d, got %s", i, v, array[i])
				break TypeSwitch
			}
		}

		ok = true

	case map[interface{}]int64:
		if got.Type() != object.ObjectTypeHash {
			t.Errorf("expect hash, got %s", got.Type())
			break TypeSwitch
		}

		hash := got.(*object.HashObject).Map
		if len(hash) != len(e) {
			t.Errorf("expect hash length %d, got %d", len(e), len(hash))
			break TypeSwitch
		}

		for k, v := range e {
			hv, has := hash[k]
			if !has {
				t.Errorf("expect hash key %+v, not found", k)
				break TypeSwitch
			}

			if hv.Value.(*object.IntegerObject).Value != v {
				t.Errorf("expect hash value on key [%+v] = %d, got %d", k, v, hv)
				break TypeSwitch
			}
		}

		ok = true
	}

	return ok
}

func runMonkeyCompatibleTest(t *testing.T, tests []monkeyTestCase) {
	t.Helper()

	for i, tt := range tests {
		m, page := testCompileCode(t, tt.input)
		m.LoadCodePage(page)
		main := page.Main().Func(nil)
		err := m.Run(main)
		if err != nil {
			t.Fatalf("run error: %s", err)
		}

		for i := 0; i < len(page.Codes); i++ {
			fmt.Printf("%d: %s\n", i, page.Codes[i].String())
		}

		got := m.Top()
		result := monkeyExpectedValueCompare(t, tt.expected, got)
		if !result {
			t.Errorf("ERROR on %d: expect %v, got %s", i, tt.expected, got)
			t.Errorf("  code: %s", tt.input)
		}
	}
}

func TestMonkeyIntegerArithmetic(t *testing.T) {
	tests := []monkeyTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestMonkeyBooleanExpression(t *testing.T) {
	tests := []monkeyTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) { 5; })", true},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []monkeyTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", nil},
		{"if (false) { 10 }", nil},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestMonkeyGlobalLetStatement(t *testing.T) {
	// NOTES: Macaque does not support global let statement.
	// The top-level variables are local variables of main function.
	tests := []monkeyTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []monkeyTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []monkeyTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []monkeyTestCase{
		{
			"{}", map[interface{}]int64{},
		},
		{
			"{1: 2, 2: 3}", map[interface{}]int64{
				int64(1): 2,
				int64(2): 3,
			},
		},
		{
			"{1 + 1: 2 * 2, 3 + 3: 4 * 4}", map[interface{}]int64{
				int64(2): 4,
				int64(6): 16,
			},
		},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []monkeyTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", nil},
		{"[1, 2, 3][99]", nil},
		{"[1][-1]", nil},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", nil},
		{"{}[0]", nil},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []monkeyTestCase{
		{
			input: `
				let fivePlusTen = fn() { 5 + 10; };
				fivePlusTen();
			`,
			expected: 15,
		},
		{
			input: `
				let one = fn() { 1; };
				let two = fn() { 2; };
				one() + two();
			`,
			expected: 3,
		},
		{
			input: `
				let a = fn() { 1; };
				let b = fn() { a() + 1; };
				let c = fn() { b() + 1; };
				c();
			`,
			expected: 3,
		},
	}

	runMonkeyCompatibleTest(t, tests)
}
