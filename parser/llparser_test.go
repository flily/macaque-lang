package parser

import (
	"testing"

	"github.com/flily/macaque-lang/ast"
	"github.com/flily/macaque-lang/lex"
)

func testLLParseCode(code string) (*ast.Program, error) {
	scanner := lex.NewRecursiveScanner("testcase")
	_ = scanner.SetContent([]byte(code))

	parser := NewLLParser(scanner)
	if err := parser.ReadTokens(); err != nil {
		return nil, err
	}

	return parser.Parse()
}

type parserTestCase struct {
	code     string
	expected *ast.Program
}

func runParserTestCase(t *testing.T, cases []parserTestCase) {
	for _, c := range cases {
		program, err := testLLParseCode(c.code)
		if err != nil {
			t.Fatalf("Parse() failed.\n%v\n%s", err, c.code)
		}

		if !program.EqualTo(c.expected) {
			t.Errorf("wrong parse result, got:\n%s\nexpected:\n%s",
				program.CanonicalCode(),
				c.expected.CanonicalCode(),
			)
		}
	}
}

func TestParseSimpleLetStatement(t *testing.T) {
	tests := []parserTestCase{
		{
			`let answer = 42;`,
			&ast.Program{
				Statements: []ast.Statement{
					&ast.LetStatement{
						Identifiers: &ast.IdentifierList{
							Identifiers: []*ast.Identifier{
								{Value: "answer"},
							},
						},
						Expressions: &ast.ExpressionList{
							Expressions: []ast.Expression{
								makeInteger("42", nil),
							},
						},
					},
				},
			},
		},
		{
			`let answer, pi, yes, no, nil = 42, 3.1415926, true, false, null;`,
			&ast.Program{
				Statements: []ast.Statement{
					&ast.LetStatement{
						Identifiers: &ast.IdentifierList{
							Identifiers: []*ast.Identifier{
								{Value: "answer"},
								{Value: "pi"},
								{Value: "yes"},
								{Value: "no"},
								{Value: "nil"},
							},
						},
						Expressions: &ast.ExpressionList{
							Expressions: []ast.Expression{
								makeInteger("42", nil),
								makeFloat("3.1415926", nil),
								makeBoolean(true, nil),
								makeBoolean(false, nil),
								makeNull(nil),
							},
						},
					},
				},
			},
		},
	}

	runParserTestCase(t, tests)
}
