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

func (t *TokenInfo) String() string {
	return fmt.Sprintf("Token{%s, %s:%d:%d}",
		t.Content, t.Line.File.Filename, t.Line.Line, t.Column,
	)
}
