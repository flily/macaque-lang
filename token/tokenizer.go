package token

// Simple tokenizer split string into word by space, only use for test.
type SimpleTokenizer struct {
	FileReader
}

func NewSimpleTokenizer(filename string) *SimpleTokenizer {
	s := &SimpleTokenizer{
		FileReader: *NewFileReader(filename),
	}

	return s
}

func (s *SimpleTokenizer) skipSpace() {
	for !s.EOF() {
		c := s.Current()
		switch c {
		case ' ', '\t', '\n', '\r':
			s.Shift()
		default:
			return
		}
	}
}

func (s *SimpleTokenizer) ScanToken() *TokenContext {
	s.skipSpace()

	start := s.Index
	for !s.EOF() {
		c := s.Current()
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			break
		}

		s.Shift()
	}

	t := s.FinishToken(String, start)
	return t
}

func (s *SimpleTokenizer) TokenList() []*TokenContext {
	tokens := make([]*TokenContext, 0)
	for !s.EOF() {
		token := s.ScanToken()
		tokens = append(tokens, token)
	}

	return tokens
}
