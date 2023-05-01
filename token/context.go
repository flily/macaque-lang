package token

type TokenContext struct {
	Token    Token
	Position *TokenInfo
	Content  string
}

type Context struct {
	Tokens []TokenContext
}
