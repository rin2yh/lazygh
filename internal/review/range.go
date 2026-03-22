package review

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
		f.rs.SetNotice("Current diff line cannot be reviewed.")
		return false
	}
	if f.rs.RangeStart != nil {
		f.rs.ClearRangeStart()
		f.rs.SetNotice("Range selection cleared.")
		return true
	}
	anchor := Range{
		Path:  line.Path,
		Index: f.selection.LineSelected(),
		Side:  string(line.Side),
	}
	if line.NewLine > 0 {
		anchor.Line = line.NewLine
	} else {
		anchor.Line = line.OldLine
	}
	f.rs.MarkRangeStart(anchor)
	f.rs.SetNotice("Range selection started.")
	return true
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
