package token

import "fmt"

// Token is set of lexical elements of the language.
type Token int

const (
	Illegal Token = iota
	EOF
	Comment

	literal_begin
	Identifier
	Null
	False
	True
	Integer
	Float
	String
	literal_end

	keyword_begin
	Let
	Fn
	Return
	If
	Else
	Import
	keyword_end

	operator_begin
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
	operator_end

	punctuation_begin
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
	punctuation_end

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

func (t Token) String() string {
	if t < 0 || t >= Token(len(displayName)) {
		return "ILLEGAL"
	}

	s := displayName[t]
	if t > operator_begin {
		return fmt.Sprintf("<%s>", s)
	}

	return s
}

func (t Token) IsLiteral() bool {
	return t > literal_begin && t < literal_end
}

func (t Token) IsKeyword() bool {
	return t > keyword_begin && t < keyword_end
}

func (t Token) IsOperator() bool {
	return t > operator_begin && t < operator_end
}

func (t Token) IsPunctuation() bool {
	return t > punctuation_begin && t < punctuation_end
}
