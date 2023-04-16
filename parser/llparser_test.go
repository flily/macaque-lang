package parser

import (
	"strings"
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

type parserErrorTestCase struct {
	lines []string
}

func (c parserErrorTestCase) code() string {
	return c.lines[0]
}

func (c parserErrorTestCase) expected() string {
	return strings.Join(c.lines, "\n")
}

func (c parserErrorTestCase) expect(got string) bool {
	return c.expected() == got
}

func runParserErrorTestCase(t *testing.T, cases []parserErrorTestCase) {
	for _, c := range cases {
		program, err := testLLParseCode(c.code())
		if err == nil {
			s := "no PROGRAM returned"
			if program != nil {
				s = program.CanonicalCode()
			}

			t.Fatalf("Parse() should fail.\n%s\ngot:\n%s", c.code(), s)
		}

		if !c.expect(err.Error()) {
			t.Errorf("wrong error message, got:\n%s\nexpected:\n%s",
				err.Error(),
				c.expected(),
			)
		}
	}
}

func TestParseLetStatement(t *testing.T) {
	tests := []parserTestCase{
		{
			`let answer = 42;`,
			program(
				let(
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
			program(
				let(
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
			program(
				let(
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

func TestParseLetStatementError(t *testing.T) {
	tests := []parserErrorTestCase{
		{
			[]string{
				`let 42 = answer;`,
				"    ^^",
				"    expect token IDENTIFIER IN identifier, but got INTEGER",
				"  at testcase:1:5",
			},
		},
		{
			[]string{
				`let answer;`,
				"          ^",
				"          expect token <=> IN let statement, but got <;>",
				"  at testcase:1:11",
			},
		},
		{
			[]string{
				`let answer = 3 + return`,
				"                 ^^^^^^",
				"                 unexpected token in EXPRESSION: RETURN",
				"  at testcase:1:18",
			},
		},
		{
			[]string{
				`let answer = return`,
				"             ^^^^^^",
				"             expect token IDENTIFIER IN expression list, but got RETURN",
				"  at testcase:1:14",
			},
		},
		{
			[]string{
				`let = 42,`,
				"    ^",
				"    expect token IDENTIFIER IN identifier, but got <=>",
				"  at testcase:1:5",
			},
		},
	}

	runParserErrorTestCase(t, tests)
}

func TestParseReturnStatement(t *testing.T) {
	tests := []parserTestCase{
		{
			`return`,
			program(
				ret(),
			),
		},
		{
			`return;`,
			program(
				ret(),
			),
		},
		{
			`return 42`,
			program(
				ret(
					l(42),
				),
			),
		},
		{
			`return 42;`,
			program(
				ret(
					l(42),
				),
			),
		},
		{
			`return 42, answer, "hello, world";`,
			program(
				ret(
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
			program(
				expr(
					l(42),
				),
			),
		},
		{
			`42,`,
			program(
				expr(
					l(42),
				),
			),
		},
		{
			`42, 3.1415926, "hello, world", true, false, null`,
			program(
				expr(
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
			program(
				expr(
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
			program(
				expr(
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
			program(
				expr(
					infix("+",
						l(1),
						infix("*", l(2), l(3)),
					),
				),
			),
		},
		{
			`1 + 2 + 3`,
			program(
				expr(
					infix("+",
						infix("+", l(1), l(2)),
						l(3),
					),
				),
			),
		},
		{
			`-1 + --2 + ---3`,
			program(
				expr(
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
			program(
				expr(
					array(),
				),
			),
		},
		{
			`[1, 2, 3]`,
			program(
				expr(
					array(
						l(1),
						l(2),
						l(3),
					),
				),
			),
		},
		{
			`[1, 2, 3,]`,
			program(
				expr(
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
			program(
				expr(
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

func TestParseArrayLiteralError(t *testing.T) {
	tests := []parserErrorTestCase{
		{
			[]string{
				`[1, 2, 3`,
				"        ^",
				"        expect token <]> IN array literal, but got EOF",
				"  at testcase:1:9",
			},
		},
		{
			[]string{
				`[1, return, 3,]`,
				"    ^^^^^^",
				"    expect token <]> IN array literal, but got RETURN",
				"  at testcase:1:5",
			},
		},
		{
			[]string{
				`[1 + return, 2, 3]`,
				"     ^^^^^^",
				"     unexpected token in EXPRESSION: RETURN",
				"  at testcase:1:6",
			},
		},
	}

	runParserErrorTestCase(t, tests)
}

func TestParseHashLiteral(t *testing.T) {
	tests := []parserTestCase{
		{
			`{}`,
			program(
				expr(
					hash(),
				),
			),
		},
		{
			`{"one": 1, "two": 2, "three": 3}`,
			program(
				expr(
					hash(
						pair(l(`"one"`), l(1)),
						pair(l(`"two"`), l(2)),
						pair(l(`"three"`), l(3)),
					),
				),
			),
		},
		{
			`{"one": 1, "two": 2, "three": 3,}`,
			program(
				expr(
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
			program(
				expr(
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

func TestParseHashLiteralError(t *testing.T) {
	tests := []parserErrorTestCase{
		{
			[]string{
				`{1: 1, 2: 2, 3: 3`,
				"                 ^",
				"                 expect token <}> IN hash literal, but got EOF",
				"  at testcase:1:18",
			},
		},
		{
			[]string{
				`{1 + return: 1},`,
				"     ^^^^^^",
				"     unexpected token in EXPRESSION: RETURN",
				"  at testcase:1:6",
			},
		},
		{
			[]string{
				`{1: 1 + return}`,
				"        ^^^^^^",
				"        unexpected token in EXPRESSION: RETURN",
				"  at testcase:1:9",
			},
		},
		{
			[]string{
				`{1, 2, 3, 4}`,
				"  ^",
				"  expect token <:> IN hash literal, but got <,>",
				"  at testcase:1:3",
			},
		},
	}

	runParserErrorTestCase(t, tests)
}

func TestFunctionLiteral(t *testing.T) {
	tests := []parserTestCase{
		{
			`fn() {}`,
			program(
				expr(
					fn(
						idList(),
						block(),
					),
				),
			),
		},
		{
			`fn(x) { x }`,
			program(
				expr(
					fn(
						idList("x"),
						block(
							expr(id("x")),
						),
					),
				),
			),
		},
		{
			`let add = fn(x) { fn(y) { x + y } }`,
			program(
				let(
					idList("add"),
					exprList(
						fn(
							idList("x"),
							block(
								expr(
									fn(
										idList("y"),
										block(
											expr(
												infix("+", id("x"), id("y")),
											),
										),
									),
								),
							),
						),
					),
				),
			),
		},
	}

	runParserTestCase(t, tests)
}

func TestParseIfExpression(t *testing.T) {
	tests := []parserTestCase{
		{
			`if (10 > 1) { 10 }`,
			program(
				expr(
					ifexp(
						infix(">", l(10), l(1)),
						block(
							expr(l(10)),
						),
						nil,
					),
				),
			),
		},
		{
			`if (10 > 1) { 10 } else { 20 }`,
			program(
				expr(
					ifexp(
						infix(">", l(10), l(1)),
						block(
							expr(l(10)),
						),
						block(
							expr(l(20)),
						),
					),
				),
			),
		},
		{
			`if (10 > 1) {
				10
			} else {
				20
			}`,
			program(
				expr(
					ifexp(
						infix(">", l(10), l(1)),
						block(
							expr(l(10)),
						),
						block(
							expr(l(20)),
						),
					),
				),
			),
		},
		{
			`if (10 > 1) {
				10
			} else if (10 > 20) {
				20
			} else {
				30
			}`,
			program(
				expr(
					ifexp(
						infix(">", l(10), l(1)),
						block(
							expr(l(10)),
						),
						elseif(
							infix(">", l(10), l(20)),
							block(
								expr(l(20)),
							),
							block(
								expr(l(30)),
							),
						),
					),
				),
			),
		},
		{
			`let x = if (10 > 1) { 10 }`,
			program(
				let(
					idList("x"),
					exprList(
						ifexp(
							infix(">", l(10), l(1)),
							block(
								expr(l(10)),
							),
							nil,
						),
					),
				),
			),
		},
		{
			`let x = 3 + if (10 > 1) { 9 } * 5 + 8`,
			program(
				let(
					idList("x"),
					exprList(
						infix("+",
							infix("+",
								l(3),
								infix("*",
									ifexp(
										infix(">", l(10), l(1)),
										block(
											expr(l(9)),
										),
										nil,
									),
									l(5),
								),
							),
							l(8),
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
