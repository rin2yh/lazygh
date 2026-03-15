package review

import "github.com/rin2yh/lazygh/internal/gh"

type Selection struct {
	File      gh.DiffFile
	Line      gh.DiffLine
	LineIndex int
}

func (s Selection) CurrentDiffFile() (gh.DiffFile, bool) {
	if s.File.Path == "" {
		return gh.DiffFile{}, false
	}
	return s.File, true
}

func (s Selection) CurrentDiffLine() (gh.DiffLine, bool) {
	if s.Line.Path == "" {
		return gh.DiffLine{}, false
	}
	return s.Line, true
}

func (s Selection) CurrentLineIndex() int {
	return s.LineIndex
}
