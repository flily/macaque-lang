package token

type FileReader struct {
	FileInfo *FileInfo
	Source   []byte
	Line     int
	Column   int
	Index    int

	start int
}

func NewFileReader(filename string) *FileReader {
	r := &FileReader{
		FileInfo: NewFileInfo(filename),
		Line:     1,
		Column:   1,
		Index:    0,
	}

	return r
}

func (r *FileReader) SetContent(data []byte) {
	r.Source = data
	lines := SplitToLines(data)
	for _, line := range lines {
		r.FileInfo.NewLine(line)
	}
}

func (r *FileReader) EOF() bool {
	return r.Index >= len(r.Source)
}

func (r *FileReader) Current() byte {
	return r.Source[r.Index]
}

func (r *FileReader) PeekChar(offset int) byte {
	return r.Source[r.Index+offset]
}

func (r *FileReader) Peek(content string) bool {
	length := len(content)
	if r.Index+length > len(r.Source) {
		return false
	}

	for i := 0; i < length; i++ {
		if r.Source[r.Index+i] != content[i] {
			return false
		}
	}

	return true
}

func (r *FileReader) CharsLeft() int {
	return len(r.Source) - r.Index
}

func (r *FileReader) shiftChar() {
	c := r.Current()
	r.Index++
	r.Column++
	if c == '\n' {
		r.Line++
		r.Column = 1
	}
}

func (r *FileReader) Shift(count int) {
	for i := 0; i < count; i++ {
		r.shiftChar()
	}
}

func (r *FileReader) StartToken() int {
	r.start = r.Index
	return r.start
}

func (r *FileReader) FinishToken(token Token) *TokenContext {
	start := r.start
	length := r.Index - r.start
	content := string(r.Source[start:r.Index])
	lineInfo := r.FileInfo.Lines[r.Line-1]
	tokenInfo := lineInfo.NewToken(r.Column-length, length, content)

	context := &TokenContext{
		Token:    token,
		Position: tokenInfo,
		Content:  content,
	}

	return context
}

// Reject token starts from current position
func (r *FileReader) RejectToken(length int) *TokenContext {
	start := r.Index
	content := string(r.Source[start : start+length])
	lineInfo := r.FileInfo.Lines[r.Line-1]
	tokenInfo := lineInfo.NewToken(r.Column, length, content)

	context := &TokenContext{
		Token:    Illegal,
		Position: tokenInfo,
		Content:  content,
	}

	return context
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
