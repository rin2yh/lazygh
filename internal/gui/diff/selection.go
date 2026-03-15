package diff

import "github.com/rin2yh/lazygh/internal/gh"

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
