package token

type TokenContext struct {
	Token    Token
	Position *TokenInfo
	Content  string
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

type Context struct {
	Tokens []*TokenContext
}
