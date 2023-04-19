package vm

import (
	"testing"

	"github.com/flily/macaque-lang/compiler"
	"github.com/flily/macaque-lang/lex"
	"github.com/flily/macaque-lang/object"
	"github.com/flily/macaque-lang/parser"
)

func testCompileCode(t *testing.T, code string) *NaiveVM {
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

	compiler := compiler.NewCompiler()
	_, err = compiler.Compile(program)
	if err != nil {
		t.Fatalf("compiler error: %s", err)
	}

	m := NewNaiveVM()
	m.LoadCode(compiler)
	m.LoadData(compiler)
	return m
}

func TestVMSimpleAction(t *testing.T) {
	code := `1 + 2`
	m := testCompileCode(t, code)

	err := m.Run(0)
	if err != nil {
		t.Fatalf("vm error: %s", err)
	}

	result := m.Top()
	if result.EqualTo(object.NewInteger(3)) == false {
		t.Fatalf("vm error: expect 3, got %s", result)
	}
}

func TestVMSimpleAction2(t *testing.T) {
	code := `"hello" + " " + "world"`
	m := testCompileCode(t, code)

	err := m.Run(0)
	if err != nil {
		t.Fatalf("vm error: %s", err)
	}

	expected := object.NewString("hello world")
	result := m.Top()
	if result.EqualTo(expected) == false {
		t.Fatalf("vm error: expect %s, got %s", expected, result)
	}
}
