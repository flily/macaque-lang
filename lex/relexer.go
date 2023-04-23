package lex

import (
	"github.com/flily/macaque-lang/errors"
	"github.com/flily/macaque-lang/token"
)

type RecursiveScanner struct {
	FileInfo *token.FileInfo

	source []byte
	line   int
	column int
	index  int
}

// NewRecursiveScanner creates a new instance of RecursiveScanner.
// Parameter filename is the filename of source file, but the scanner DOES NOT
// read the file directly via filename.
func NewRecursiveScanner(filename string) *RecursiveScanner {
	s := &RecursiveScanner{
		FileInfo: token.NewFileInfo(filename),
		line:     1,
		column:   1,
		index:    0,
	}

	return s
}

func (s *RecursiveScanner) splitToLines(data []byte) []string {
	lines := make([]string, 0)
	start := 0
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			line := string(data[start:i])
			start = i + 1
			lines = append(lines, line)
		}
	}

	if start <= len(data) {
		line := string(data[start:])
		lines = append(lines, line)
	}

	return lines
}

func (s *RecursiveScanner) SetContent(data []byte) error {
	if s.source != nil {
		return ErrScannerHasContentAlready
	}

	s.source = data
	lines := s.splitToLines(data)
	for _, line := range lines {
		s.FileInfo.NewLine(line)
	}

	return nil
}

func (s *RecursiveScanner) ReadContentSlice(start int) string {
	return string(s.source[start:s.index])
}

func (s *RecursiveScanner) Append(data []byte) {
	// FIXME: lines count is wrong
	lines := s.splitToLines(data)
	for _, line := range lines {
		s.FileInfo.NewLine(line)
	}

	s.source = append(s.source, data...)
	s.source = append(s.source, '\n')
}

func (s *RecursiveScanner) makeCurrentPosition(content string) *token.TokenInfo {
	length := len(content)
	column := s.column - length

	return s.FileInfo.Lines[s.line-1].NewToken(column, length, content)
}

func (s *RecursiveScanner) makeEOFPosition() *token.TokenInfo {
	length := 1
	column := s.column

	return s.FileInfo.Lines[s.line-1].NewToken(column, length, "")
}

func (s *RecursiveScanner) Scan() (*LexicalElement, error) {
	return s.scanStateInit()
}

func (s *RecursiveScanner) charsLeft() int {
	return len(s.source) - s.index
}

func (s *RecursiveScanner) EOF() bool {
	return s.index >= len(s.source)
}

func (s *RecursiveScanner) makeCurrentCodeContext(length int) *errors.CodeContext {
	line := s.FileInfo.Lines[s.line-1]
	ctx := &errors.CodeContext{
		Filename:  s.FileInfo.Filename,
		Line:      line.Content,
		NumLine:   s.line,
		NumColumn: s.column,
		Length:    length,
	}

	return ctx
}

func (s *RecursiveScanner) currentChar() byte {
	return s.source[s.index]
}

func (s *RecursiveScanner) peekChar(offset int) byte {
	return s.source[s.index+offset]
}

func (s *RecursiveScanner) shift_one() {
	// Check EOF every time before calling shift() and shift_one()
	// unnecessary to check EOF here

	c := s.currentChar()
	s.index++
	s.column++
	if c == '\n' {
		s.line++
		s.column = 1
	}
}

func (s *RecursiveScanner) shift(length int) {
	// length is checked before calling shift()
	i := 0
	for i = 0; i < length; i++ {
		s.shift_one()
	}
}

func (s *RecursiveScanner) skipWhitespace() int {
	i := 0

SkipWhitespace:
	for !s.EOF() {
		switch s.source[s.index] {
		case ' ', '\t', '\r', '\n':
			s.shift(1)
			i += 1

		default:
			break SkipWhitespace
		}
	}

	return i
}

func (s *RecursiveScanner) peek(content string) bool {
	length := len(content)
	if s.index+length > len(s.source) {
		return false
	}

	for i := 0; i < length; i++ {
		if s.source[s.index+i] != content[i] {
			return false
		}
	}

	return true
}

func (s *RecursiveScanner) scanStateInit() (*LexicalElement, error) {
	var elem *LexicalElement
	var err error

	s.skipWhitespace()
	if s.EOF() {
		return s.ScanEOF(), nil
	}

	c := s.currentChar()
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

func (s *RecursiveScanner) scanElementNumber() (*LexicalElement, error) {
	if s.peek("0x") {
		return s.scanElementHexdecimal()
	}

	return s.scanElementDecimal()
}

func (s *RecursiveScanner) scanElementHexdecimal() (*LexicalElement, error) {
	start := s.index
	// ensure the first two characters are "0x"
	s.shift(2)

	for !s.EOF() {
		c := s.currentChar()
		if IsHexDigit(c) || c == '_' {
			s.shift(1)

		} else {
			break
		}
	}

	content := s.ReadContentSlice(start)
	elem := &LexicalElement{
		Token:    token.Integer,
		Content:  content,
		Position: s.makeCurrentPosition(content),
	}

	return elem, nil
}

func (s *RecursiveScanner) scanElementDecimal() (*LexicalElement, error) {
	start := s.index
	for !s.EOF() {
		c := s.currentChar()
		if IsDigit(c) || c == '_' {
			s.shift(1)

		} else {
			if c == '.' {
				s.shift(1)
				return s.scanElementFloat(start)
			}

			break
		}
	}

	content := s.ReadContentSlice(start)
	elem := &LexicalElement{
		Token:    token.Integer,
		Content:  content,
		Position: s.makeCurrentPosition(content),
	}
	return elem, nil
}

func (s *RecursiveScanner) scanElementFloat(start int) (*LexicalElement, error) {
	for !s.EOF() {
		c := s.currentChar()
		if IsDigit(c) || c == '_' {
			s.shift(1)

		} else {
			break
		}
	}

	content := s.ReadContentSlice(start)
	elem := &LexicalElement{
		Token:    token.Float,
		Content:  content,
		Position: s.makeCurrentPosition(content),
	}
	return elem, nil
}

func (s *RecursiveScanner) scanElementIdentifierOrKeyword() (*LexicalElement, error) {
	start := s.index
	for !s.EOF() {
		c := s.currentChar()
		if IsUpper(c) || IsLower(c) || IsDigit(c) || c == '_' {
			s.shift(1)

		} else {
			break
		}
	}

	content := s.ReadContentSlice(start)
	tokenType := token.CheckKeywordToken(content)

	elem := &LexicalElement{
		Token:    tokenType,
		Content:  content,
		Position: s.makeCurrentPosition(content),
	}

	return elem, nil
}

func (s *RecursiveScanner) scanStateString() (*LexicalElement, error) {
	start := s.index
	s.shift(1) // include the first '"'

StringLoop:
	for !s.EOF() {
		c := s.currentChar()
		s.shift(1)

		switch c {
		case '"':
			break StringLoop

		case '\\':
			if s.EOF() {
				ctx := s.makeCurrentCodeContext(1)
				return nil, ctx.NewSyntaxError("unexpected EOF")
			}

			n := s.currentChar()
			switch n {
			case '\\', 'n', 'r', 't', '"':
				s.shift(1)

			case 'x':
				s.shift(1)
				if charsLeft := s.charsLeft(); charsLeft < 2 {
					ctx := s.makeCurrentCodeContext(charsLeft + 1)
					return nil, ctx.NewSyntaxError("insufficient characters for escape sequence")
				}

				n1 := s.currentChar()
				n2 := s.peekChar(1)
				if IsHexDigit(n1) && IsHexDigit(n2) {
					s.shift(2)

				} else {
					ctx := s.makeCurrentCodeContext(2)
					return nil, ctx.NewSyntaxError("invalid escape sequence: \\x%c%c", n1, n2)
				}

			default:
				ctx := s.makeCurrentCodeContext(1)
				return nil, ctx.NewSyntaxError("invalid escape sequence: \\%c", n)
			}
		}
	}

	content := s.ReadContentSlice(start)
	elem := &LexicalElement{
		Token:    token.String,
		Content:  content,
		Position: s.makeCurrentPosition(content),
	}

	return elem, nil
}

func (s *RecursiveScanner) makeForwardLexicalElement(length int) *LexicalElement {
	start := s.index
	s.shift(length)

	content := s.ReadContentSlice(start)
	tokenType := token.CheckOperatorToken(content)
	elem := &LexicalElement{
		Token:    tokenType,
		Content:  content,
		Position: s.makeCurrentPosition(content),
	}

	return elem
}

func (s *RecursiveScanner) tryScanPunctuations(ps ...string) *LexicalElement {
	for _, p := range ps {
		if s.peek(p) {
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

func (s *RecursiveScanner) scanStatePunctuation() (*LexicalElement, error) {
	var elem *LexicalElement
	var err error

	if s.peek("//") {
		return s.scanStateComment()
	}

	c := s.currentChar()
	elem = s.tryScanPunctuations(multiBytesPunctutations...)
	if elem != nil {
		return elem, nil
	}

	tokenType := token.CheckOperatorToken(string(c))
	if tokenType != token.Illegal {
		elem = s.makeForwardLexicalElement(1)

	} else {
		ctx := s.makeCurrentCodeContext(1)
		err = ctx.NewSyntaxError("unknown operator '%c'", c)
	}

	return elem, err
}

func (s *RecursiveScanner) scanStateComment() (*LexicalElement, error) {
	start := s.index
	s.shift(2) // shift the first '//'

	for !s.EOF() {
		c := s.currentChar()
		if c == '\n' {
			break
		}

		s.shift(1)
	}

	content := s.ReadContentSlice(start)
	elem := &LexicalElement{
		Token:    token.Comment,
		Content:  content,
		Position: s.makeCurrentPosition(content),
	}

	return elem, nil
}

func (s *RecursiveScanner) ScanEOF() *LexicalElement {
	elem := &LexicalElement{
		Token:    token.EOF,
		Position: s.makeEOFPosition(),
	}

	return elem
}
