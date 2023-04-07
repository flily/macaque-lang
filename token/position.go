package token

import "fmt"

type FileInfo struct {
	Filename string
	Lines    []*LineInfo
}

type LineInfo struct {
	Line    int
	Content string
	File    *FileInfo
	Tokens  []*TokenInfo
}

type TokenInfo struct {
	Column  int
	Length  int
	Content string
	Line    *LineInfo
}

func NewFileInfo(filename string) *FileInfo {
	info := &FileInfo{
		Filename: filename,
		Lines:    make([]*LineInfo, 0),
	}

	return info
}

func (f *FileInfo) NewLine(content string) *LineInfo {
	lineno := len(f.Lines) + 1
	info := &LineInfo{
		Content: content,
		Line:    lineno,
		File:    f,
	}

	f.Lines = append(f.Lines, info)
	return info
}

func (l *LineInfo) NewToken(start int, length int, content string) *TokenInfo {
	info := &TokenInfo{
		Column:  start,
		Length:  length,
		Content: content,
		Line:    l,
	}

	l.Tokens = append(l.Tokens, info)
	return info
}

// GetLineNumber returns the line number of the token in source file.
// Line number starts from 1.
func (t *TokenInfo) GetLineNumber() int {
	return t.Line.Line
}

// GetColumnNumber returns the column number of the token in source file.
// Column number starts from 1.
func (t *TokenInfo) GetColumnNumber() int {
	return t.Column
}

// GetPosition returns the line number and column number of the token in source file.
func (t *TokenInfo) GetPosition() (int, int) {
	line := t.GetLineNumber()
	column := t.GetColumnNumber()

	return line, column
}

func (t *TokenInfo) String() string {
	line, column := t.GetPosition()
	return fmt.Sprintf("Token{%s, %s:%d:%d}",
		t.Content, t.Line.File.Filename, line, column,
	)
}
