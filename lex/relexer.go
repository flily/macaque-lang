package lex

import (
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

	if start < len(data)-1 {
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

	lines := s.splitToLines(data)
	for _, line := range lines {
		s.FileInfo.NewLine(line)
	}

	s.source = append(s.source, data...)
}

func (s *RecursiveScanner) makeCurrentPosition(content string) *token.TokenInfo {
	length := len(content)
	column := s.column - length

	return s.FileInfo.Lines[s.line-1].NewToken(column, length, content)
}

func (s *RecursiveScanner) Scan() (*LexicalElement, error) {
	return s.scan_state_init()
}

func (s *RecursiveScanner) EOF() bool {
	return s.index >= len(s.source)
}

func (s *RecursiveScanner) currentChar() byte {
	return s.source[s.index]
}

func (s *RecursiveScanner) peekChar() byte {
	return s.source[s.index+1]
}

func (s *RecursiveScanner) shift_one() int {
	if s.EOF() {
		return 0
	}

	c := s.currentChar()
	s.index++
	s.column++
	if c == '\n' {
		s.line++
		s.column = 1
	}

	return 1
}

func (s *RecursiveScanner) shift(length ...int) int {
	if len(length) <= 0 {
		return s.shift_one()
	}

	l := length[0]
	for i := 0; i < l; i++ {
		if ok := s.shift_one(); ok <= 0 {
			break
		}
	}

	return l
}

func (s *RecursiveScanner) skip_whitespace() {
	length := len(s.source)
	for s.index < length {
		switch s.source[s.index] {
		case ' ', '\t', '\r', '\n':
			s.shift()
		default:
			return
		}
	}
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

func (s *RecursiveScanner) scan_state_init() (*LexicalElement, error) {
	var elem *LexicalElement
	var err error

	s.skip_whitespace()
	c := s.currentChar()
	switch {
	case c >= '0' && c <= '9':
		elem, err = s.scan_element_number()
	}

	return elem, err
}

func (s *RecursiveScanner) scan_element_number() (*LexicalElement, error) {
	if s.peek("0x") {
		return s.scan_element_hexdecimal()
	}

	return s.scan_element_decimal()
}

func (s *RecursiveScanner) scan_element_hexdecimal() (*LexicalElement, error) {
	start := s.index
	// ensure the first two characters are "0x"
	s.shift(2)

	for !s.EOF() {
		c := s.currentChar()
		if IsHexDigit(c) || c == '_' {
			s.shift()
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

func (s *RecursiveScanner) scan_element_decimal() (*LexicalElement, error) {
	start := s.index
	for !s.EOF() {
		c := s.currentChar()
		if IsDigit(c) || c == '_' {
			s.shift()
		} else {
			if c == '.' {
				s.shift()
				return s.scan_element_float(start)
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

func (s *RecursiveScanner) scan_element_float(start int) (*LexicalElement, error) {
	for !s.EOF() {
		c := s.currentChar()
		if IsDigit(c) || c == '_' {
			s.shift()
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
