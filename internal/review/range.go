package review

import "github.com/rin2yh/lazygh/internal/model"

type rangeState struct {
	rs        *ReviewState
	selection Selection
	setFocus  func(FocusTarget)
}

func newRange(rs *ReviewState, selection Selection, setFocus func(FocusTarget)) *rangeState {
	return &rangeState{
		rs:        rs,
		selection: selection,
		setFocus:  setFocus,
	}
}

func (f *rangeState) RangeStart() *model.ReviewRange {
	return f.rs.RangeStart
}

func (f *rangeState) ToggleSelection() {
	line, ok := f.selection.CurrentDiffLine()
	if !ok || !line.Commentable {
		f.rs.SetNotice("Current diff line cannot be reviewed.")
		return
	}
	if f.rs.RangeStart != nil {
		f.rs.ClearRangeStart()
		f.rs.SetNotice("Range selection cleared.")
		f.setFocus(FocusDiffContent)
		return
	}
	anchor := model.ReviewRange{
		Path:  line.Path,
		Index: f.selection.CurrentLineIndex(),
		Side:  string(line.Side),
	}
	if line.NewLine > 0 {
		anchor.Line = line.NewLine
	} else {
		anchor.Line = line.OldLine
	}
	f.rs.MarkRangeStart(anchor)
	f.rs.SetNotice("Range selection started.")
	f.setFocus(FocusDiffContent)
}

func (f *rangeState) IsIndexWithinPendingRange(path string, commentable bool, idx int) bool {
	start := f.rs.RangeStart
	if start == nil || start.Path != path || !commentable {
		return false
	}
	minIndex := start.Index
	maxIndex := f.selection.CurrentLineIndex()
	if minIndex > maxIndex {
		minIndex, maxIndex = maxIndex, minIndex
	}
	return idx >= minIndex && idx <= maxIndex
}
