package token

import (
	"fmt"
	"strings"
)

type TokenContext struct {
	Token    Token
	Position *TokenInfo
	Content  string
}

func (c *TokenContext) Length() int {
	return c.Position.Length
}

func (c *TokenContext) Filename() string {
	return c.Position.Line.File.Filename
}

func (c *TokenContext) LineNo() int {
	return c.Position.Line.Line
}

func (c *TokenContext) ColumnStart() int {
	return c.Position.Column
}

func (c *TokenContext) ColumnEnd() int {
	return c.Position.Column + c.Position.Length
}

func (c *TokenContext) NewSyntaxError(format string, args ...interface{}) *SyntaxNError {
	return MakeSyntaxError(c, format, args...)
}

type Context struct {
	Tokens []*TokenContext
}

func (c *Context) makeLineTokenMap() [][]*TokenContext {
	lines := make([]int, 0)
	tokens := make(map[int][]*TokenContext)

	for _, token := range c.Tokens {
		lineInfo := token.Position.Line
		lineNo := lineInfo.Line
		if _, hasLine := tokens[lineNo]; !hasLine {
			tokens[lineNo] = make([]*TokenContext, 0)
			lines = append(lines, lineNo)
		}

		tokens[lineNo] = append(tokens[lineNo], token)
	}

	result := make([][]*TokenContext, 0, len(lines))
	for _, lineNo := range lines {
		result = append(result, tokens[lineNo])
	}

	return result
}

func (c *Context) makeLineHighlight(content string, tokens []*TokenContext) (string, string) {
	parts := make([]string, 0)
	last := 1
	lead := ""

	for i, token := range tokens {
		spaces := readSpaces(content, last-1, token.ColumnStart()-1)
		if i == 0 {
			lead = spaces
		}

		highlight := strings.Repeat("^", token.Length())
		parts = append(parts, spaces, highlight)
		last = token.ColumnEnd()
	}

	return strings.Join(parts, ""), lead
}

func (c *Context) HighLight() string {
	lineTokens := c.makeLineTokenMap()
	results := make([]string, 2*len(lineTokens))

	for i, tokens := range lineTokens {
		first := tokens[0]
		content := first.Position.Line.Content
		highlight, _ := c.makeLineHighlight(content, tokens)
		results[2*i] = content
		results[2*i+1] = highlight
	}

	return strings.Join(results, "\n")
}

func (c *Context) Message(format string, args ...interface{}) string {
	lineTokens := c.makeLineTokenMap()
	results := make([]string, 2*len(lineTokens)+2)
	lead := ""

	for i, tokens := range lineTokens {
		first := tokens[0]
		content := first.Position.Line.Content
		highlight, l := c.makeLineHighlight(content, tokens)
		lead = l
		results[2*i] = content
		results[2*i+1] = highlight
	}

	results[2*len(lineTokens)] = lead + fmt.Sprintf(format, args...)
	results[2*len(lineTokens)+1] = fmt.Sprintf("  at %s:%d:%d",
		c.Tokens[0].Filename(),
		c.Tokens[0].LineNo(),
		c.Tokens[0].ColumnStart(),
	)

	return strings.Join(results, "\n")
}

func readSpaces(s string, start int, end int) string {
	chars := make([]byte, end-start)

	for i := start; i < end; i++ {
		c := s[i]
		switch c {
		case '\t', '\v', '\f':
			chars[i-start] = c
		default:
			chars[i-start] = ' '
		}
	}

	return string(chars)
}
