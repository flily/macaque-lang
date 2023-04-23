package vm

import (
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
		return got.EqualTo(null)
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

		got := m.Top()
		result := monkeyExpectedValueCompare(t, tt.expected, got)
		if !result {
			t.Errorf("ERROR on %d: expect %v, got %s", i, tt.expected, got)
			t.Errorf("  code: %s", tt.input)
		}

		// there are ONLY one element on the stack, ONLY true in original monkey.
		expectedStack := uint64(main.FrameSize + 2)
		if m.sp != expectedStack {
			t.Errorf("ERROR on %d: expect stack size %d, got %d", i, expectedStack, m.sp)
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

func TestFunctionsWithReturnStatement(t *testing.T) {
	tests := []monkeyTestCase{
		{
			input: `
				let earlyExit = fn() { return 99; 100; };
				earlyExit();
			`,
			expected: 99,
		},
		{
			input: `
				let earlyExit = fn() { return 99; return 100; };
				earlyExit();
			`,
			expected: 99,
		},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []monkeyTestCase{
		{
			input: `
				let noReturn = fn() { };
				noReturn();
			`,
			expected: nil,
		},
		{
			input: `
				let noReturn = fn() { };
				let noReturnTwo = fn() { noReturn(); };
				noReturn();
				noReturnTwo();
			`,
			expected: nil,
		},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestFirstClassFunctions(t *testing.T) {
	tests := []monkeyTestCase{
		{
			input: `
				let returnsOne = fn() { 1; };
				let returnsOneReturner = fn() { returnsOne; };
				returnsOneReturner()();
			`,
			expected: 1,
		},
		{
			input: `
				let returnsOneReturner = fn() {
					let returnsOne = fn() { 1; };
					returnsOne;
				};
				returnsOneReturner()();
			`,
			expected: 1,
		},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []monkeyTestCase{
		{
			input: `
				let one = fn() { let one = 1; one; };
				one();
			`,
			expected: 1,
		},
		{
			input: `
				let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
				oneAndTwo();
			`,
			expected: 3,
		},
		{
			input: `
				let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
				let threeAndFour = fn() { let three = 3; let four = 4; three + four; };
				oneAndTwo() + threeAndFour();
			`,
			expected: 10,
		},
		{
			input: `
				let firstFoobar = fn() { let foobar = 50; foobar; };
				let secondFoobar = fn() { let foobar = 100; foobar; };
				firstFoobar() + secondFoobar();
			`,
			expected: 150,
		},
		{
			input: `
				let globalSeed = 50;
				let minusOne = fn() {
					let num = 1;
					globalSeed - num;
				};
				let minusTwo = fn() {
					let num = 2;
					globalSeed - num;
				};
				minusOne() + minusTwo();
			`,
			expected: 97,
		},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []monkeyTestCase{
		{
			input: `
				let identity = fn(a) { a; };
				identity(4);
			`,
			expected: 4,
		},
		{
			input: `
				let sum = fn(a, b) { a + b; };
				sum(1, 2);
			`,
			expected: 3,
		},
		{
			input: `
				let sum = fn(a, b) {
					let c = a + b;
					c;
				};
				sum(1, 2);
			`,
			expected: 3,
		},
		{
			input: `
				let sum = fn(a, b) {
					let c = a + b;
					c;
				};
				sum(1, 2) + sum(3, 4);
			`,
			expected: 10,
		},
		{
			input: `
				let sum = fn(a, b) {
					let c = a + b;
					c;
				};
				let outer = fn() {
					sum(1, 2) + sum(3, 4);
				};
				outer();
			`,
			expected: 10,
		},
		{
			input: `
				let globalNum = 10;

				let sum = fn(a, b) {
					let c = a + b;
					c + globalNum;
				};

				let outer = fn() {
					sum(1, 2) + sum(3, 4) + globalNum;
				};

				outer() + globalNum;
			`,
			expected: 50,
		},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestClosure(t *testing.T) {
	tests := []monkeyTestCase{
		{
			input: `
				let newClosure = fn(a) {
					fn() { a; };
				};
				let closure = newClosure(99);
				closure();
			`,
			expected: 99,
		},
		{
			input: `
				let newAdder = fn(a, b) {
					fn(c) { a + b + c };
				};
				let adder = newAdder(1, 2);
				adder(8);
			`,
			expected: 11,
		},
		{
			input: `
				let newAdder = fn(a, b) {
					let c = a + b;
					fn(d) { c + d };
				};
				let adder = newAdder(1, 2);
				adder(8);
			`,
			expected: 11,
		},
		{
			input: `
				let newAdderOuter = fn(a, b) {
					let c = a + b;
					fn(d) {
						let e = d + c;
						fn(f) { e + f };
					};
				};
				let newAdderInner = newAdderOuter(1, 2);
				let adder = newAdderInner(3);
				adder(8);
			`,
			expected: 14,
		},
		{
			input: `
				let a = 1;
				let newAdderOuter = fn(b) {
					fn(c) {
						fn(d) { a + b + c + d };
					};
				};
				let newAdderInner = newAdderOuter(2);
				let adder = newAdderInner(3);
				adder(8);
			`,
			expected: 14,
		},
		{
			input: `
				let newClosure = fn(a, b) {
					let one = fn() { a; };
					let two = fn() { b; };
					fn() { one() + two(); };
				};
				let closure = newClosure(9, 90);
				closure();
			`,
			expected: 99,
		},
	}

	runMonkeyCompatibleTest(t, tests)
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []monkeyTestCase{
		{
			input: `
				let countDown = fn(x) {
					if (x == 0) {
						return 0;
					} else {
						fn(x - 1);
					}
				};
				countDown(1);
			`,
			expected: 0,
		},
		{
			input: `
				let countDown = fn(x) {
					if (x == 0) {
						return 0;
					} else {
						fn(x - 1);
					}
				};
				let wrapper = fn() {
					countDown(1);
				};
				wrapper();
			`,
			expected: 0,
		},
		{
			input: `
				let wrapper = fn() {
					let countDown = fn(x) {
						if (x == 0) {
							return 0;
						} else {
							fn(x - 1);
						}
					};
					countDown(1);
				};
				wrapper();
			`,
			expected: 0,
		},
	}

	runMonkeyCompatibleTest(t, tests)
}
