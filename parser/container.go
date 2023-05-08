package parser

import (
	"github.com/flily/macaque-lang/lex"
	"github.com/flily/macaque-lang/token"
)

type CodeContainer struct {
	Elements []*token.TokenContext
	Index    int

	scanner lex.Scanner
}

func NewContainer(scanner lex.Scanner) *CodeContainer {
	c := &CodeContainer{
		Elements: make([]*token.TokenContext, 0),
		scanner:  scanner,
	}

	return c
}

func (c *CodeContainer) ReadTokens() error {
	var err error

	for {
		var elem *token.TokenContext
		elem, err = c.scanner.Scan()
		if err != nil {
			break
		}

		c.Elements = append(c.Elements, elem)
		if elem.Token == token.EOF {
			break
		}
	}

	return err
}

func (c *CodeContainer) current() *token.TokenContext {
	return c.Elements[c.Index]
}

func (c *CodeContainer) Current() *token.TokenContext {
	if c.Index < len(c.Elements) {
		return c.current()
	}

	return nil
}

func (c *CodeContainer) Peek(offset int) *token.TokenContext {
	index := c.Index + offset
	if 0 < index && index < len(c.Elements) {
		return c.Elements[index]
	}

	return nil
}

func (c *CodeContainer) Next() *token.TokenContext {
	if c.Index < len(c.Elements) {
		c.Index++
	}

	return c.current()
}
