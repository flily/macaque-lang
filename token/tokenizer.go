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
			s.Shift(1)
		default:
			return
		}
	}
}

func (s *SimpleTokenizer) ScanToken() *TokenContext {
	s.skipSpace()

	s.StartToken()
	for !s.EOF() {
		c := s.Current()
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			break
		} else if c == '~' {
			return s.RejectToken(1)
		}

		s.Shift(1)
	}

	t := s.FinishToken(String)
	return t
}

func (s *SimpleTokenizer) TokenList() []*TokenContext {
	tokens := make([]*TokenContext, 0)
	for !s.EOF() {
		token := s.ScanToken()
		tokens = append(tokens, token)
		if token.Token == Illegal {
			break
		}
	}

	return tokens
}
