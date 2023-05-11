package token

import "fmt"

// Token is set of lexical elements of the language.
type Token int

const (
	Illegal Token = iota
	Nil           // Go language level nil
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

	punctuationBegin //
	Assign           // =
	Comma            // ,
	Period           // .
	Colon            // :
	Semicolon        // ;
	DualColon        // ::
	bracketBegin     //
	LParen           // (
	RParen           // )
	LBrace           // {
	RBrace           // }
	LBracket         // [
	RBracket         // ]
	CommentStart     // //
	bracketEnd       //
	punctuationEnd   //
	operatorEnd      //
	LastToken        // This is not a token, but a marker for the last token.

	SIllegal      = "ILLEGAL"
	SEOF          = "EOF"
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
	SDualColon    = "::"
	SLParen       = "("
	SRParen       = ")"
	SLBrace       = "{"
	SRBrace       = "}"
	SLBracket     = "["
	SRBracket     = "]"
	SCommentStart = "//"
	SLastToken    = "<LAST_TOKEN>"
)

var displayToken = [...]string{
	Illegal:   SIllegal,
	EOF:       SEOF,
	Let:       SLet,
	Fn:        SFn,
	Return:    SReturn,
	If:        SIf,
	Else:      SElse,
	Import:    SImport,
	Null:      SNull,
	False:     SFalse,
	True:      STrue,
	Bang:      SBang,
	Plus:      SPlus,
	Minus:     SMinus,
	Asterisk:  SAsterisk,
	Slash:     SSlash,
	EQ:        SEQ,
	NE:        SNE,
	LT:        SLT,
	GT:        SGT,
	Modulo:    SModulo,
	LE:        SLE,
	GE:        SGE,
	AND:       SAND,
	OR:        SOR,
	BITAND:    SBITAND,
	BITOR:     SBITOR,
	BITXOR:    SBITXOR,
	BITNOT:    SBITNOT,
	Assign:    SAssign,
	Comma:     SComma,
	Period:    SPeriod,
	Colon:     SColon,
	Semicolon: SSemicolon,
	DualColon: SDualColon,
	LParen:    SLParen,
	RParen:    SRParen,
	LBrace:    SLBrace,
	RBrace:    SRBrace,
	LBracket:  SLBracket,
	RBracket:  SRBracket,
	LastToken: SLastToken,
}

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
	Bang:         "BANG",
	Plus:         "PLUS",
	Minus:        "MINUS",
	Asterisk:     "ASTERISK",
	Slash:        "SLASH",
	EQ:           "EQ",
	NE:           "NE",
	LT:           "LT",
	GT:           "GT",
	Modulo:       "MODULO",
	LE:           "LE",
	GE:           "GE",
	AND:          "AND",
	OR:           "OR",
	BITAND:       "BITAND",
	BITOR:        "BITOR",
	BITXOR:       "BITXOR",
	BITNOT:       "BITNOT",
	Assign:       "ASSIGN",
	Comma:        "COMMA",
	Period:       "PERIOD",
	Colon:        "COLON",
	Semicolon:    "SEMICOLON",
	DualColon:    "DUALCOLON",
	LParen:       "LPAREN",
	RParen:       "RPAREN",
	LBrace:       "LBRACE",
	RBrace:       "RBRACE",
	LBracket:     "LBRACKET",
	RBracket:     "RBRACKET",
	CommentStart: "COMMENT",
	LastToken:    "<LAST>",
}

// String returns a string representation of the token.
func (t Token) String() string {
	if t < 0 || t >= Token(len(displayName)) {
		return "ILLEGAL"
	}

	name := displayName[t]
	if t > operatorBegin {
		token := displayToken[t]
		if bracketBegin < t && t < bracketEnd {
			return fmt.Sprintf("%s('%s')", name, token)
		} else {
			return fmt.Sprintf("%s(%s)", name, token)
		}
	}

	return name
}

func (t Token) CodeName() string {
	if t < 0 || t >= Token(len(displayName)) {
		return "ILLEGAL"
	}

	if t > operatorBegin {
		return displayToken[t]
	} else {
		return displayName[t]
	}
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

var keywordMap = map[string]Token{
	SLet:    Let,
	SFn:     Fn,
	SReturn: Return,
	SIf:     If,
	SElse:   Else,
	SImport: Import,
	SNull:   Null,
	SFalse:  False,
	STrue:   True,
}

// CheckKeywordToken returns keyword token when the given string is keyword,
// otherwise returns Identifier.
func CheckKeywordToken(s string) Token {
	if v, ok := keywordMap[s]; ok {
		return v
	}

	return Identifier
}

var operatorMap = map[string]Token{
	SBang:         Bang,
	SPlus:         Plus,
	SMinus:        Minus,
	SAsterisk:     Asterisk,
	SSlash:        Slash,
	SEQ:           EQ,
	SNE:           NE,
	SLT:           LT,
	SGT:           GT,
	SModulo:       Modulo,
	SLE:           LE,
	SGE:           GE,
	SAND:          AND,
	SOR:           OR,
	SBITAND:       BITAND,
	SBITOR:        BITOR,
	SBITXOR:       BITXOR,
	SBITNOT:       BITNOT,
	SAssign:       Assign,
	SComma:        Comma,
	SPeriod:       Period,
	SColon:        Colon,
	SSemicolon:    Semicolon,
	SDualColon:    DualColon,
	SLParen:       LParen,
	SRParen:       RParen,
	SLBrace:       LBrace,
	SRBrace:       RBrace,
	SLBracket:     LBracket,
	SRBracket:     RBracket,
	SCommentStart: CommentStart,
}

// CheckOperatorToken returns operator token when the given string is operator,
// otherwise returns Illegal.
func CheckOperatorToken(s string) Token {
	if v, ok := operatorMap[s]; ok {
		return v
	}

	return Illegal
}
