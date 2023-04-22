package token

import (
	"fmt"
	"strings"

	"github.com/flily/macaque-lang/errors"
)

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

func (t *TokenInfo) MakeContext() *errors.CodeContext {
	ctx := &errors.CodeContext{
		Filename:  t.Line.File.Filename,
		NumLine:   t.GetLineNumber(),
		NumColumn: t.GetColumnNumber(),
		Line:      t.Line.Content,
		Length:    t.Length,
	}

	return ctx
}

func (t *TokenInfo) GetLeadingSpaces(index int) string {
	line := t.Line.Content
	spaces := make([]byte, 0, index)

	for i := 0; i < index; i++ {
		switch line[i] {
		case '\t', '\v', '\f':
			spaces = append(spaces, line[i])
		default:
			spaces = append(spaces, ' ')
		}
	}

	return string(spaces)
}

func (t *TokenInfo) MakeLineHighlight() string {
	leadingSpaces := t.GetLeadingSpaces(t.Column - 1)
	highlight := leadingSpaces + strings.Repeat("^", t.Length)

	return t.Line.Content + "\n" + highlight
}

func (t *TokenInfo) MakeMessage(format string, args ...interface{}) string {
	leadingSpaces := t.GetLeadingSpaces(t.Column - 1)
	highlight := leadingSpaces + strings.Repeat("^", t.Length)
	lines := []string{
		t.Line.Content,
		highlight,
		leadingSpaces + fmt.Sprintf(format, args...),
		fmt.Sprintf("  at %s:%d:%d", t.Line.File.Filename, t.Line.Line, t.Column),
	}

	return strings.Join(lines, "\n")
}
