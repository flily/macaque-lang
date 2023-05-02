package token

import (
	"fmt"
)

// FileInfo is the root node of all positions.
type FileInfo struct {
	Filename string
	Lines    []*LineInfo
}

// LineInfo holds the information of a line, and tokens in this line.
type LineInfo struct {
	Line    int
	Content string
	File    *FileInfo
	Tokens  []*TokenInfo
}

// TokenInfo represents a token in source file, it can be a lexical token parsed
// by lexer or parser, or an invalid token which is not accepted by lexer.
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

func (l *LineInfo) NewToken(startColumn int, length int, content string) *TokenInfo {
	info := &TokenInfo{
		Column:  startColumn,
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

func (t *TokenInfo) MakeContext() *CodeContext {
	ctx := &CodeContext{
		Filename:  t.Line.File.Filename,
		NumLine:   t.GetLineNumber(),
		NumColumn: t.GetColumnNumber(),
		Line:      t.Line.Content,
		Length:    t.Length,
	}

	return ctx
}

func (t *TokenInfo) MakeLineHighlight() string {
	return t.MakeContext().MakeLineHighlight()
}

func (t *TokenInfo) MakeMessage(format string, args ...interface{}) string {
	return t.MakeContext().MakeMessage(format, args...)
}
