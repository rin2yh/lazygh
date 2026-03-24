package diff

import "github.com/rin2yh/lazygh/internal/gh"

type State struct {
	files        []gh.DiffFile
	fileSelected int
	lineSelected int
}

func (s *State) Reset() {
	s.files = nil
	s.fileSelected = 0
	s.lineSelected = 0
}

func (s *State) Files() []gh.DiffFile {
	return s.files
}

func (s *State) SetFiles(files []gh.DiffFile) {
	s.files = files
}

func (s *State) FileSelected() int {
	return s.fileSelected
}

func (s *State) SetFileSelected(selected int) {
	s.fileSelected = selected
}

func (s *State) LineSelected() int {
	return s.lineSelected
}

func (s *State) SetLineSelected(selected int) {
	s.lineSelected = selected
}

func (s *State) CurrentFile() (gh.DiffFile, bool) {
	return CurrentFile(s.files, s.fileSelected)
}

func (s *State) CurrentLine() (gh.DiffLine, bool) {
	file, ok := s.CurrentFile()
	if !ok {
		return gh.DiffLine{}, false
	}
	return CurrentLine(file, s.lineSelected)
}

func (s *State) EnsureLineSelection() {
	file, ok := s.CurrentFile()
	if !ok {
		s.lineSelected = 0
		return
	}
	s.lineSelected = EnsureLineSelection(file, s.lineSelected)
}

func (s *State) SelectNextFile() bool {
	if len(s.files) == 0 || s.fileSelected >= len(s.files)-1 {
		return false
	}
	s.fileSelected++
	s.lineSelected = 0
	s.EnsureLineSelection()
	return true
}

func (s *State) SelectPrevFile() bool {
	if len(s.files) == 0 || s.fileSelected <= 0 {
		return false
	}
	s.fileSelected--
	s.lineSelected = 0
	s.EnsureLineSelection()
	return true
}

func (s *State) SelectNextLine(step int) bool {
	file, ok := s.CurrentFile()
	if !ok {
		return false
	}
	next, changed := SelectNextLine(file, s.lineSelected, step)
	if !changed {
		return false
	}
	s.lineSelected = next
	return true
}

func (s *State) SelectPrevLine(step int) bool {
	file, ok := s.CurrentFile()
	if !ok {
		return false
	}
	prev, changed := SelectPrevLine(file, s.lineSelected, step)
	if !changed {
		return false
	}
	s.lineSelected = prev
	return true
}

func (s *State) GotoFirstLine() bool {
	next, changed := GotoFirstLine(s.lineSelected)
	if !changed {
		return false
	}
	s.lineSelected = next
	return true
}

func (s *State) GotoLastLine() bool {
	file, ok := s.CurrentFile()
	if !ok {
		return false
	}
	next, changed := GotoLastLine(file, s.lineSelected)
	if !changed {
		return false
	}
	s.lineSelected = next
	return true
}

func CurrentFile(files []gh.DiffFile, selected int) (gh.DiffFile, bool) {
	if len(files) == 0 || selected < 0 || selected >= len(files) {
		return gh.DiffFile{}, false
	}
	return files[selected], true
}

func CurrentLine(file gh.DiffFile, selected int) (gh.DiffLine, bool) {
	if len(file.Lines) == 0 || selected < 0 || selected >= len(file.Lines) {
		return gh.DiffLine{}, false
	}
	return file.Lines[selected], true
}

func EnsureLineSelection(file gh.DiffFile, selected int) int {
	if len(file.Lines) == 0 {
		return 0
	}
	if selected < 0 {
		return 0
	}
	if selected >= len(file.Lines) {
		return len(file.Lines) - 1
	}
	return selected
}

func SelectNextLine(file gh.DiffFile, current int, step int) (int, bool) {
	if len(file.Lines) == 0 {
		return current, false
	}
	if step < 1 {
		step = 1
	}
	next := current + step
	if next >= len(file.Lines) {
		next = len(file.Lines) - 1
	}
	if next == current {
		return current, false
	}
	return next, true
}

func SelectPrevLine(file gh.DiffFile, current int, step int) (int, bool) {
	if len(file.Lines) == 0 {
		return current, false
	}
	if step < 1 {
		step = 1
	}
	prev := current - step
	if prev < 0 {
		prev = 0
	}
	if prev == current {
		return current, false
	}
	return prev, true
}

func GotoFirstLine(current int) (int, bool) {
	if current == 0 {
		return current, false
	}
	return 0, true
}

func GotoLastLine(file gh.DiffFile, current int) (int, bool) {
	if len(file.Lines) == 0 {
		return current, false
	}
	last := len(file.Lines) - 1
	if current == last {
		return current, false
	}
	return last, true
}

func ParseFiles(prevFiles []gh.DiffFile, prevSelected int, content string) ([]gh.DiffFile, int, int) {
	files := gh.ParseUnifiedDiff(content)
	if len(files) == 0 {
		return nil, 0, 0
	}

	prevPath := ""
	if prevSelected >= 0 && prevSelected < len(prevFiles) {
		prevPath = prevFiles[prevSelected].Path
	}

	selected := 0
	if prevPath != "" {
		for idx, file := range files {
			if file.Path == prevPath {
				selected = idx
				break
			}
		}
	}

	return files, selected, 0
}

func LineIndex(file gh.DiffFile, target gh.DiffLine) int {
	for idx, candidate := range file.Lines {
		if candidate == target {
			return idx
		}
	}
	return -1
}
