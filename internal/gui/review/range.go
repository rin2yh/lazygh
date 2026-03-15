package review

import (
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/gui/diff"
)

type rangeState struct {
	state     *core.State
	selection Selection
	setFocus  func(FocusTarget)
}

func newRange(state *core.State, selection Selection, setFocus func(FocusTarget)) *rangeState {
	return &rangeState{
		state:     state,
		selection: selection,
		setFocus:  setFocus,
	}
}

func (f *rangeState) RangeStart() *core.ReviewRange {
	return f.state.Review.RangeStart
}

func (f *rangeState) Clear() {
	f.state.ClearReviewRangeStart()
}

func (f *rangeState) ToggleSelection() {
	line, ok := f.selection.CurrentDiffLine()
	if !ok || !line.Commentable {
		f.state.SetReviewNotice("Current diff line cannot be reviewed.")
		return
	}
	if f.state.Review.RangeStart != nil {
		f.state.ClearReviewRangeStart()
		f.state.SetReviewNotice("Range selection cleared.")
		f.setFocus(FocusDiffContent)
		return
	}
	anchor := core.ReviewRange{
		Path:  line.Path,
		Index: f.selection.CurrentLineIndex(),
		Side:  string(line.Side),
	}
	if line.NewLine > 0 {
		anchor.Line = line.NewLine
	} else {
		anchor.Line = line.OldLine
	}
	f.state.MarkReviewRangeStart(anchor)
	f.state.SetReviewNotice("Range selection started.")
	f.setFocus(FocusDiffContent)
}

func (f *rangeState) IsLineWithinPendingRange(line gh.DiffLine) bool {
	start := f.state.Review.RangeStart
	if start == nil || start.Path != line.Path || !line.Commentable {
		return false
	}
	file, ok := f.selection.CurrentDiffFile()
	if !ok {
		return false
	}
	lineIndex := diff.LineIndex(file, line)
	if lineIndex < 0 {
		return false
	}
	minIndex := start.Index
	maxIndex := f.selection.CurrentLineIndex()
	if minIndex > maxIndex {
		minIndex, maxIndex = maxIndex, minIndex
	}
	return lineIndex >= minIndex && lineIndex <= maxIndex
}
