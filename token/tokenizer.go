package token

// Simple tokenizer split string into word by space, only use for test.
type SimpleTokenizer struct {
	FileInfo *FileInfo

	source []byte
	line   int
	column int
	index  int
}

func NewSimpleTokenizer(filename string) *SimpleTokenizer {
	s := &SimpleTokenizer{
		FileInfo: NewFileInfo(filename),
		line:     1,
		column:   1,
		index:    0,
	}

	return s
}

func (s *SimpleTokenizer) SetContent(data []byte) {
	s.source = data
	lines := SplitToLines(data)
	for _, line := range lines {
		s.FileInfo.NewLine(line)
	}
}

func (s *SimpleTokenizer) eof() bool {
	return s.index >= len(s.source)
}

func (s *SimpleTokenizer) current() byte {
	return s.source[s.index]
}

func (s *SimpleTokenizer) shift() {
	c := s.source[s.index]
	s.index++
	s.column++
	if c == '\n' {
		s.line++
		s.column = 1
	}
}

func (s *SimpleTokenizer) skipSpace() {
	for s.index < len(s.source) {
		c := s.current()
		switch c {
		case ' ', '\t', '\n', '\r':
			s.shift()
		default:
			return
		}
	}
}

func (s *SimpleTokenizer) ScanToken() TokenContext {
	s.skipSpace()

	start := s.index
	for !s.eof() {
		c := s.current()
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			break
		}

		s.shift()
	}

	content := string(s.source[start:s.index])
	line := s.line
	info := s.FileInfo.Lines[line-1].NewToken(s.column-len(content), len(content), content)
	context := TokenContext{
		Token:    String,
		Position: info,
		Content:  content,
	}

	return context
}

func (s *SimpleTokenizer) TokenList() []TokenContext {
	tokens := make([]TokenContext, 0)
	for !s.eof() {
		token := s.ScanToken()
		tokens = append(tokens, token)
	}

	return tokens
}

func SplitToLines(data []byte) []string {
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
