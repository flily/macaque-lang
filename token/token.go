package token

import "fmt"

// Token is set of lexical elements of the language.
type Token int

const (
	Illegal Token = iota
	EOF
	Comment

	literalBegin
	Identifier
	Null
	False
	True
	Integer
	Float
	String
	literalEnd

	keywordBegin
	Let
	Fn
	Return
	If
	Else
	Import
	keywordEnd

	operatorBegin
	Bang     // !
	Plus     // +
	Minus    // -
	Asterisk // *
	Slash    // /
	EQ       // ==
	NE       // !=
	LT       // <
	GT       // >
	Modulo   // %
	LE       // <=
	GE       // >=
	AND      // &&
	OR       // ||
	BITAND   // &
	BITOR    // |
	BITXOR   // ^
	BITNOT   // ~
	operatorEnd

	punctuationBegin
	Assign       // =
	Comma        // ,
	Period       // .
	Colon        // :
	Semicolon    // ;
	LParen       // (
	RParen       // )
	LBrace       // {
	RBrace       // }
	LBracket     // [
	RBracket     // ]
	CommentStart // //
	punctuationEnd

	SLet          = "let"
	SFn           = "fn"
	SReturn       = "return"
	SIf           = "if"
	SElse         = "else"
	SImport       = "import"
	SNull         = "null"
	SFalse        = "false"
	STrue         = "true"
	SBang         = "!"
	SPlus         = "+"
	SMinus        = "-"
	SAsterisk     = "*"
	SSlash        = "/"
	SEQ           = "=="
	SNE           = "!="
	SLT           = "<"
	SGT           = ">"
	SModulo       = "%"
	SLE           = "<="
	SGE           = ">="
	SAND          = "&&"
	SOR           = "||"
	SBITAND       = "&"
	SBITOR        = "|"
	SBITXOR       = "^"
	SBITNOT       = "~"
	SAssign       = "="
	SComma        = ","
	SPeriod       = "."
	SColon        = ":"
	SSemicolon    = ";"
	SLParen       = "("
	SRParen       = ")"
	SLBrace       = "{"
	SRBrace       = "}"
	SLBracket     = "["
	SRBracket     = "]"
	SCommentStart = "//"
)

var displayName = [...]string{
	Illegal:      "ILLEGAL",
	EOF:          "EOF",
	Comment:      "COMMENT",
	Identifier:   "IDENTIFIER",
	Null:         "NULL",
	False:        "FALSE",
	True:         "TRUE",
	Integer:      "INTEGER",
	Float:        "FLOAT",
	String:       "STRING",
	Let:          "LET",
	Fn:           "FN",
	Return:       "RETURN",
	If:           "IF",
	Else:         "ELSE",
	Import:       "IMPORT",
	Bang:         "!",
	Plus:         "+",
	Asterisk:     "*",
	Slash:        "/",
	EQ:           "==",
	NE:           "!=",
	LT:           "<",
	GT:           ">",
	Modulo:       "%",
	LE:           "<=",
	GE:           ">=",
	AND:          "&&",
	OR:           "||",
	BITAND:       "&",
	BITOR:        "|",
	BITXOR:       "^",
	BITNOT:       "~",
	Assign:       "=",
	Comma:        ",",
	Period:       ".",
	Colon:        ":",
	Semicolon:    ";",
	LParen:       "(",
	RParen:       ")",
	LBrace:       "{",
	RBrace:       "}",
	LBracket:     "[",
	RBracket:     "]",
	CommentStart: "//",
}

// String returns a string representation of the token.
func (t Token) String() string {
	if t < 0 || t >= Token(len(displayName)) {
		return "ILLEGAL"
	}

	s := displayName[t]
	if t > operatorBegin {
		return fmt.Sprintf("<%s>", s)
	}

	return s
}

func (t Token) IsLiteral() bool {
	return t > literalBegin && t < literalEnd
}

func (t Token) IsKeyword() bool {
	return t > keywordBegin && t < keywordEnd
}

func (t Token) IsOperator() bool {
	return t > operatorBegin && t < operatorEnd
}

func (t Token) IsPunctuation() bool {
	return t > punctuationBegin && t < punctuationEnd
}
