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
			makeProgram(
				makeLetStatement(
					makeIdentifierList(
						"answer"),
					makeExpressionList(
						makeInteger("42", nil),
					),
				),
			),
		},
		{
			`let answer, pi, hello,  yes, no, nil = 
				42, 3.1415926, "hello, world", true, false, null;`,
			makeProgram(
				makeLetStatement(
					makeIdentifierList(
						"answer", "pi", "hello", "yes", "no", "nil"),
					makeExpressionList(
						makeInteger("42", nil),
						makeFloat("3.1415926", nil),
						makeString(`"hello, world"`, nil),
						makeBoolean(true, nil),
						makeBoolean(false, nil),
						makeNull(nil),
					),
				),
			),
		},
		{
			`let answer, pi, hello,  yes, no, nil, = 
				42, 3.1415926, "hello, world", true, false, null,;`,
			makeProgram(
				makeLetStatement(
					makeIdentifierList(
						"answer", "pi", "hello", "yes", "no", "nil"),
					makeExpressionList(
						makeInteger("42", nil),
						makeFloat("3.1415926", nil),
						makeString(`"hello, world"`, nil),
						makeBoolean(true, nil),
						makeBoolean(false, nil),
						makeNull(nil),
					),
				),
			),
		},
	}

	runParserTestCase(t, tests)
}

func TestParseExpressionList(t *testing.T) {
	tests := []parserTestCase{
		{
			`42`,
			makeProgram(
				makeExpressionStatement(
					makeInteger("42", nil),
				),
			),
		},
		{
			`42,`,
			makeProgram(
				makeExpressionStatement(
					makeInteger("42", nil),
				),
			),
		},
		{
			`42, 3.1415926, "hello, world", true, false, null`,
			makeProgram(
				makeExpressionStatement(
					makeInteger("42", nil),
					makeFloat("3.1415926", nil),
					makeString(`"hello, world"`, nil),
					makeBoolean(true, nil),
					makeBoolean(false, nil),
					makeNull(nil),
				),
			),
		},
		{
			`42, 3.1415926, "hello, world", true, false, null,`,
			makeProgram(
				makeExpressionStatement(
					makeInteger("42", nil),
					makeFloat("3.1415926", nil),
					makeString(`"hello, world"`, nil),
					makeBoolean(true, nil),
					makeBoolean(false, nil),
					makeNull(nil),
				),
			),
		},
	}

	runParserTestCase(t, tests)
}

func TestExpressionParsing(t *testing.T) {
	tests := []parserTestCase{
		{
			`1 + 2 * 3`,
			makeProgram(
				makeExpressionStatement(
					makeInfixExpression("+",
						makeInteger("1", nil),
						makeInfixExpression("*",
							makeInteger("2", nil),
							makeInteger("3", nil),
						),
					),
				),
			),
		},
		{
			`1 + 2 + 3`,
			makeProgram(
				makeExpressionStatement(
					makeInfixExpression("+",
						makeInfixExpression("*",
							makeInteger("1", nil),
							makeInteger("2", nil),
						),
						makeInteger("3", nil),
					),
				),
			),
		},
		{
			`-1 + --2 + ---3`,
			makeProgram(
				makeExpressionStatement(
					makeInfixExpression("+",
						makeInfixExpression("+",
							makePrefixExpression("-", makeInteger("1", nil)),
							makePrefixExpression("-",
								makePrefixExpression("-", makeInteger("2", nil)),
							),
						),
						makePrefixExpression("-",
							makePrefixExpression("-",
								makePrefixExpression("-", makeInteger("3", nil)),
							),
						),
					),
				),
			),
		},
	}

	runParserTestCase(t, tests)
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((- a) * b);",
		},
		{
			"!-a",
			"(! (- a));",
		},
		{
			"a + b + c",
			"((a + b) + c);",
		},
		{
			"a + b - c",
			"((a + b) - c);",
		},
		{
			"a * b * c",
			"((a * b) * c);",
		},
		{
			"a * b / c",
			"((a * b) / c);",
		},
		{
			"a + b / c",
			"(a + (b / c));",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f);",
		},
		{
			"3 + 4; -5 * 5",
			makeMultilines(
				"(3 + 4);",
				"((- 5) * 5);",
			),
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4));",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4));",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));",
		},
		{
			"true",
			"true;",
		},
		{
			"false",
			"false;",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false);",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true);",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4);",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2);",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5));",
		},
		{
			"-(5 + 5)",
			"(- (5 + 5));",
		},
		{
			"!(true == true)",
			"(! (true == true));",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d);",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)));",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g));",
		},
		{
			"a + add(b, c) + std[sum](c, d)",
			"((a + add(b, c)) + std[sum](c, d));",
		},
		{
			"a + add(b, c) + std.sum(c, d)",
			"((a + add(b, c)) + std[sum](c, d));",
		},
		// {
		// 	"a * [1, 2, 3, 4][b * c] * d",
		// 	"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		// },
		// {
		// 	"add(a * b[2], b[1], 2 * [1, 2][1])",
		// 	"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		// },
	}

	for _, tt := range tests {
		program, err := testLLParseCode(tt.input)
		if err != nil {
			t.Fatalf("Parse code error:\n%s", err)
		}

		if program.CanonicalCode() != tt.expected {
			t.Errorf("Wrong precedence, got:\n%s\nexpected:\n%s",
				program.CanonicalCode(), tt.expected)
		}
	}
}
