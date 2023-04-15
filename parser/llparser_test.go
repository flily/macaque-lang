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
			t.Errorf("Parse() failed.\n%v\n%s", err, c.code)
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
					idList(
						"answer"),
					exprList(
						l(42),
					),
				),
			),
		},
		{
			`let answer, pi, hello,  yes, no, nil = 
				42, 3.1415926, "hello, world", true, false, null;`,
			makeProgram(
				makeLetStatement(
					idList(
						"answer", "pi", "hello", "yes", "no", "nil"),
					exprList(
						l(42),
						float("3.1415926"),
						l(`"hello, world"`),
						l(true),
						l(false),
						l(nil),
					),
				),
			),
		},
		{
			`let answer, pi, hello,  yes, no, nil, = 
				42, 3.1415926, "hello, world", true, false, null,;`,
			makeProgram(
				makeLetStatement(
					idList(
						"answer", "pi", "hello", "yes", "no", "nil"),
					exprList(
						l(42),
						float("3.1415926"),
						l(`"hello, world"`),
						l(true),
						l(false),
						l(nil),
					),
				),
			),
		},
	}

	runParserTestCase(t, tests)
}

func TestParseReturnStatement(t *testing.T) {
	tests := []parserTestCase{
		{
			`return`,
			makeProgram(
				makeReturnStatement(),
			),
		},
		{
			`return;`,
			makeProgram(
				makeReturnStatement(),
			),
		},
		{
			`return 42`,
			makeProgram(
				makeReturnStatement(
					l(42),
				),
			),
		},
		{
			`return 42;`,
			makeProgram(
				makeReturnStatement(
					l(42),
				),
			),
		},
		{
			`return 42, answer, "hello, world";`,
			makeProgram(
				makeReturnStatement(
					l(42),
					id("answer"),
					l(`"hello, world"`),
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
					l(42),
				),
			),
		},
		{
			`42,`,
			makeProgram(
				makeExpressionStatement(
					l(42),
				),
			),
		},
		{
			`42, 3.1415926, "hello, world", true, false, null`,
			makeProgram(
				makeExpressionStatement(
					l(42),
					float("3.1415926"),
					l(`"hello, world"`),
					l(true),
					l(false),
					l(nil),
				),
			),
		},
		{
			`42, 3.1415926, "hello, world", true, false, null,`,
			makeProgram(
				makeExpressionStatement(
					l(42),
					float("3.1415926"),
					l(`"hello, world"`),
					l(true),
					l(false),
					l(nil),
				),
			),
		},
		{
			`42, 3.1415926, "hello, world", true, false, null,;`,
			makeProgram(
				makeExpressionStatement(
					l(42),
					float("3.1415926"),
					l(`"hello, world"`),
					l(true),
					l(false),
					l(nil),
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
					infix("+",
						l(1),
						infix("*", l(2), l(3)),
					),
				),
			),
		},
		{
			`1 + 2 + 3`,
			makeProgram(
				makeExpressionStatement(
					infix("+",
						infix("+", l(1), l(2)),
						l(3),
					),
				),
			),
		},
		{
			`-1 + --2 + ---3`,
			makeProgram(
				makeExpressionStatement(
					infix("+",
						infix("+",
							prefix("-", l(1)),
							prefix("-",
								prefix("-", l(2)),
							),
						),
						prefix("-",
							prefix("-",
								prefix("-", l(3)),
							),
						),
					),
				),
			),
		},
	}

	runParserTestCase(t, tests)
}

func TestParseArrayLiteral(t *testing.T) {
	tests := []parserTestCase{
		{
			`[]`,
			makeProgram(
				makeExpressionStatement(
					array(),
				),
			),
		},
		{
			`[1, 2, 3]`,
			makeProgram(
				makeExpressionStatement(
					array(
						l(1),
						l(2),
						l(3),
					),
				),
			),
		},
		{
			`[1, 2 * 2, 3 + 3]`,
			makeProgram(
				makeExpressionStatement(
					array(
						l(1),
						infix("*", l(2), l(2)),
						infix("+", l(3), l(3)),
					),
				),
			),
		},
	}

	runParserTestCase(t, tests)
}

func TestParseHashLiteral(t *testing.T) {
	tests := []parserTestCase{
		{
			`{}`,
			makeProgram(
				makeExpressionStatement(
					hash(),
				),
			),
		},
		{
			`{"one": 1, "two": 2, "three": 3}`,
			makeProgram(
				makeExpressionStatement(
					hash(
						pair(l(`"one"`), l(1)),
						pair(l(`"two"`), l(2)),
						pair(l(`"three"`), l(3)),
					),
				),
			),
		},
		{
			`{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`,
			makeProgram(
				makeExpressionStatement(
					hash(
						pair(l(`"one"`),
							infix("+", l(0), l(1))),
						pair(l(`"two"`),
							infix("-", l(10), l(8))),
						pair(l(`"three"`),
							infix("/", l(15), l(5))),
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
			"((a + add(b, c)) + (std[sum])(c, d));",
		},
		{
			"a + add(b, c) + std.sum(c, d)",
			"((a + add(b, c)) + (std[sum])(c, d));",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d);",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])));",
		},
		{
			"a + b::c(d + e) * f",
			"(a + (b::c((d + e)) * f));",
		},
		{
			"a * b::c(d + e) + f",
			"((a * b::c((d + e))) + f);",
		},
	}

	for _, tt := range tests {
		program, err := testLLParseCode(tt.input)
		if err != nil {
			t.Errorf("Parse code error:\n%s", err)
		}

		if program.CanonicalCode() != tt.expected {
			t.Errorf("Wrong precedence, got:\n%s\nexpected:\n%s",
				program.CanonicalCode(), tt.expected)
		}
	}
}
