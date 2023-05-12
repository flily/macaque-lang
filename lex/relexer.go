package lex

import (
	"github.com/flily/macaque-lang/token"
)

type RecursiveScanner struct {
	token.FileReader
}

// NewRecursiveScanner creates a new instance of RecursiveScanner.
// Parameter filename is the filename of source file, but the scanner DOES NOT
// read the file directly via filename.
func NewRecursiveScanner(filename string) *RecursiveScanner {
	s := &RecursiveScanner{
		FileReader: *token.NewFileReader(filename),
	}

	return s
}

func (s *RecursiveScanner) Scan() (*token.TokenContext, error) {
	return s.scanStateInit()
}

func (s *RecursiveScanner) skipWhitespace() int {
	i := 0

SkipWhitespace:
	for !s.EOF() {
		switch s.Current() {
		case ' ', '\t', '\r', '\n':
			s.Shift(1)
			i += 1

		default:
			break SkipWhitespace
		}
	}

	return i
}

func (s *RecursiveScanner) RejectError(length int, format string, args ...interface{}) *LexicalError {
	ctx := s.RejectToken(length)
	return NewLexicalError(ctx, format, args...)
}

func (s *RecursiveScanner) scanStateInit() (*token.TokenContext, error) {
	var elem *token.TokenContext
	var err error

	s.skipWhitespace()
	if s.EOF() {
		return s.ScanEOF(), nil
	}

	c := s.Current()
	switch {
	case IsDigit(c):
		elem, err = s.scanElementNumber()

	case IsUpper(c) || IsLower(c) || c == '_':
		elem, err = s.scanElementIdentifierOrKeyword()

	case IsPunct(c) && c != '"':
		elem, err = s.scanStatePunctuation()

	case c == '"':
		elem, err = s.scanStateString()
	}

	return elem, err
}

func (s *RecursiveScanner) scanElementNumber() (*token.TokenContext, error) {
	if s.Peek("0x") {
		return s.scanElementHexdecimal()
	}

	return s.scanElementDecimal()
}

func (s *RecursiveScanner) scanElementHexdecimal() (*token.TokenContext, error) {
	s.StartToken()
	// ensure the first two characters are "0x"
	s.Shift(2)

	for !s.EOF() {
		c := s.Current()
		if IsHexDigit(c) || c == '_' {
			s.Shift(1)

		} else {
			break
		}
	}

	elem := s.FinishToken(token.Integer)
	return elem, nil
}

func (s *RecursiveScanner) scanElementDecimal() (*token.TokenContext, error) {
	s.StartToken()
	for !s.EOF() {
		c := s.Current()
		if IsDigit(c) || c == '_' {
			s.Shift(1)

		} else {
			if c == '.' {
				s.Shift(1)
				return s.scanElementFloat()
			}

			break
		}
	}

	elem := s.FinishToken(token.Integer)
	return elem, nil
}

func (s *RecursiveScanner) scanElementFloat() (*token.TokenContext, error) {
	for !s.EOF() {
		c := s.Current()
		if IsDigit(c) || c == '_' {
			s.Shift(1)

		} else {
			break
		}
	}

	elem := s.FinishToken(token.Float)
	return elem, nil
}

func (s *RecursiveScanner) scanElementIdentifierOrKeyword() (*token.TokenContext, error) {
	s.StartToken()
	for !s.EOF() {
		c := s.Current()
		if IsUpper(c) || IsLower(c) || IsDigit(c) || c == '_' {
			s.Shift(1)

		} else {
			break
		}
	}

	elem := s.FinishToken(token.Identifier)
	elem.Token = token.CheckKeywordToken(elem.Content)

	return elem, nil
}

func (s *RecursiveScanner) scanStateString() (*token.TokenContext, error) {
	s.StartToken()
	s.Shift(1) // include the first '"'

StringLoop:
	for !s.EOF() {
		c := s.Current()
		s.Shift(1)

		switch c {
		case '"':
			break StringLoop

		case '\\':
			if s.EOF() {
				err := s.RejectError(1, "unexpected EOF")
				return nil, err
			}

			n := s.Current()
			switch n {
			case '\\', 'n', 'r', 't', '"':
				s.Shift(1)

			case 'x':
				s.Shift(1)
				if charsLeft := s.CharsLeft(); charsLeft < 2 {
					err := s.RejectError(charsLeft+1, "insufficient characters for escape sequence")
					return nil, err
				}

				n1 := s.Current()
				n2 := s.PeekChar(1)
				if IsHexDigit(n1) && IsHexDigit(n2) {
					s.Shift(2)

				} else {
					err := s.RejectError(2, "invalid escape sequence: \\x%c%c", n1, n2)
					return nil, err
				}

			default:
				err := s.RejectError(1, "invalid escape sequence: \\%c", n)
				return nil, err
			}
		}
	}

	elem := s.FinishToken(token.String)
	return elem, nil
}

func (s *RecursiveScanner) makeForwardLexicalElement(length int) *token.TokenContext {
	s.StartToken()
	s.Shift(length)

	elem := s.FinishToken(token.Illegal)
	tokenType := token.CheckOperatorToken(elem.Content)
	elem.Token = tokenType
	return elem
}

func (s *RecursiveScanner) tryScanPunctuations(ps ...string) *token.TokenContext {
	for _, p := range ps {
		if s.Peek(p) {
			return s.makeForwardLexicalElement(len(p))
		}
	}

	return nil
}

var multiBytesPunctutations = []string{
	token.SEQ, token.SAssign,
	token.SNE, token.SBang,
	token.SGE, token.SGT,
	token.SLE, token.SLT,
	token.SDualColon, token.SColon,
}

func (s *RecursiveScanner) scanStatePunctuation() (*token.TokenContext, error) {
	var elem *token.TokenContext
	var err error

	if s.Peek("//") {
		return s.scanStateComment()
	}

	c := s.Current()
	elem = s.tryScanPunctuations(multiBytesPunctutations...)
	if elem != nil {
		return elem, nil
	}

	tokenType := token.CheckOperatorToken(string(c))
	if tokenType != token.Illegal {
		elem = s.makeForwardLexicalElement(1)

	} else {
		err = s.RejectError(1, "unknown operator '%c'", c)
	}

	return elem, err
}

func (s *RecursiveScanner) scanStateComment() (*token.TokenContext, error) {
	s.StartToken()
	s.Shift(2) // shift the first '//'

	for !s.EOF() {
		c := s.Current()
		if c == '\n' {
			break
		}

		s.Shift(1)
	}

	elem := s.FinishToken(token.Comment)
	return elem, nil
}

func (s *RecursiveScanner) ScanEOF() *token.TokenContext {
	s.StartToken()
	elem := s.FinishToken(token.EOF)
	// FIXME: should not hack token length
	elem.Position.Length = 1
	return elem
}
