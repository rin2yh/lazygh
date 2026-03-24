package review

import "github.com/rin2yh/lazygh/internal/gh"

// Range identifies a diff line position for range-based comments.
type Range struct {
	Path      string
	Index     int
	Side      gh.DiffSide
	Line      int
	StartSide gh.DiffSide
	StartLine int
}

type rangeState struct {
	rs        *ReviewState
	selection Selection
}

func newRange(rs *ReviewState, selection Selection) *rangeState {
	return &rangeState{
		rs:        rs,
		selection: selection,
	}
}

func (f *rangeState) RangeStart() *Range {
	return f.rs.RangeStart
}

// ToggleSelection toggles range selection and returns true if focus should
// move to FocusDiffContent.
func (f *rangeState) ToggleSelection() bool {
	line, ok := f.selection.CurrentLine()
	if !ok || !line.Commentable {
		f.rs.Notify("Current diff line cannot be reviewed.")
		return false
	}
	if f.rs.RangeStart != nil {
		f.rs.ClearRangeStart()
		f.rs.Notify("Range selection cleared.")
		return true
	}
	anchor := Range{
		Path:  line.Path,
		Index: f.selection.LineSelected(),
		Side:  line.Side,
	}
	if line.NewLine > 0 {
		anchor.Line = line.NewLine
	} else {
		anchor.Line = line.OldLine
	}
	f.rs.MarkRangeStart(anchor)
	f.rs.Notify("Range selection started.")
	return true
}

// HasConflict reports whether the range start is in a different file than the
// current cursor, making range-based commenting impossible without clearing it.
func (f *rangeState) HasConflict() bool {
	start := f.rs.RangeStart
	if start == nil {
		return false
	}
	file, ok := f.selection.CurrentFile()
	if !ok {
		return false
	}
	return file.Path != start.Path
}

func (f *rangeState) IsIndexWithinPendingRange(path string, commentable bool, idx int) bool {
	start := f.rs.RangeStart
	if start == nil || start.Path != path || !commentable {
		return false
	}
	minIndex := start.Index
	maxIndex := f.selection.LineSelected()
	if minIndex > maxIndex {
		minIndex, maxIndex = maxIndex, minIndex
	}
	return idx >= minIndex && idx <= maxIndex
}
