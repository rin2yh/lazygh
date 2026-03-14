package gui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
)

func (gui *Gui) shouldShowReviewDrawer() bool {
	if !gui.state.IsDiffMode() {
		return false
	}
	review := gui.state.Review
	return review.DrawerOpen || review.InputMode != core.ReviewInputNone || gui.state.HasPendingReview() || len(review.Comments) > 0 || review.Summary != "" || review.RangeStart != nil
}

func (gui *Gui) currentDiffFile() (gh.DiffFile, bool) {
	if len(gui.diffFiles) == 0 || gui.diffFileSelected < 0 || gui.diffFileSelected >= len(gui.diffFiles) {
		return gh.DiffFile{}, false
	}
	return gui.diffFiles[gui.diffFileSelected], true
}

func (gui *Gui) currentDiffLine() (gh.DiffLine, bool) {
	file, ok := gui.currentDiffFile()
	if !ok || len(file.Lines) == 0 || gui.diffLineSelected < 0 || gui.diffLineSelected >= len(file.Lines) {
		return gh.DiffLine{}, false
	}
	return file.Lines[gui.diffLineSelected], true
}

func (gui *Gui) ensureDiffLineSelection() {
	file, ok := gui.currentDiffFile()
	if !ok || len(file.Lines) == 0 {
		gui.diffLineSelected = 0
		return
	}
	if gui.diffLineSelected < 0 {
		gui.diffLineSelected = 0
	}
	if gui.diffLineSelected >= len(file.Lines) {
		gui.diffLineSelected = len(file.Lines) - 1
	}
}

func (gui *Gui) selectNextDiffLine(step int) bool {
	file, ok := gui.currentDiffFile()
	if !ok || len(file.Lines) == 0 {
		return false
	}
	if step < 1 {
		step = 1
	}
	next := gui.diffLineSelected + step
	if next >= len(file.Lines) {
		next = len(file.Lines) - 1
	}
	if next == gui.diffLineSelected {
		return false
	}
	gui.diffLineSelected = next
	return true
}

func (gui *Gui) selectPrevDiffLine(step int) bool {
	file, ok := gui.currentDiffFile()
	if !ok || len(file.Lines) == 0 {
		return false
	}
	if step < 1 {
		step = 1
	}
	prev := gui.diffLineSelected - step
	if prev < 0 {
		prev = 0
	}
	if prev == gui.diffLineSelected {
		return false
	}
	gui.diffLineSelected = prev
	return true
}

func (gui *Gui) gotoFirstDiffLine() bool {
	if gui.diffLineSelected == 0 {
		return false
	}
	gui.diffLineSelected = 0
	return true
}

func (gui *Gui) gotoLastDiffLine() bool {
	file, ok := gui.currentDiffFile()
	if !ok || len(file.Lines) == 0 {
		return false
	}
	last := len(file.Lines) - 1
	if gui.diffLineSelected == last {
		return false
	}
	gui.diffLineSelected = last
	return true
}

func (gui *Gui) beginReviewCommentInput() {
	gui.state.BeginReviewCommentInput()
	gui.state.OpenReviewDrawer()
	gui.focus = panelReviewDrawer
	gui.commentEditor.Focus()
	gui.commentEditor.SetValue("")
}

func (gui *Gui) beginReviewSummaryInput() {
	gui.state.BeginReviewSummaryInput()
	gui.state.OpenReviewDrawer()
	gui.focus = panelReviewDrawer
	gui.summaryEditor.SetValue(gui.state.Review.Summary)
	gui.summaryEditor.Focus()
}

func (gui *Gui) stopReviewInput() {
	gui.commentEditor.Blur()
	gui.summaryEditor.Blur()
	if gui.state.Review.InputMode == core.ReviewInputComment {
		gui.state.ClearReviewRangeStart()
		gui.commentEditor.SetValue("")
	}
	gui.state.Review.InputMode = core.ReviewInputNone
	if gui.shouldShowReviewDrawer() {
		gui.focus = panelReviewDrawer
	}
}

func (gui *Gui) handleReviewEditorKey(msg tea.KeyMsg) bool {
	switch msg.Type {
	case tea.KeyEsc:
		gui.stopReviewInput()
		gui.focus = panelDiffContent
		return true
	}
	switch msg.String() {
	case "ctrl+s":
		switch gui.state.Review.InputMode {
		case core.ReviewInputComment:
			return true
		case core.ReviewInputSummary:
			gui.state.SetReviewSummary(gui.summaryEditor.Value())
			gui.stopReviewInput()
			gui.state.SetReviewNotice("Review summary updated.")
			return true
		}
	}

	switch gui.state.Review.InputMode {
	case core.ReviewInputComment:
		updated, _ := gui.commentEditor.Update(msg)
		gui.commentEditor = updated
		return true
	case core.ReviewInputSummary:
		updated, _ := gui.summaryEditor.Update(msg)
		gui.summaryEditor = updated
		return true
	default:
		return false
	}
}

func (gui *Gui) toggleReviewRangeSelection() {
	line, ok := gui.currentDiffLine()
	if !ok || !line.Commentable {
		gui.state.SetReviewNotice("Current diff line cannot be reviewed.")
		return
	}
	if gui.state.Review.RangeStart != nil {
		gui.state.ClearReviewRangeStart()
		gui.state.SetReviewNotice("Range selection cleared.")
		gui.focus = panelDiffContent
		return
	}
	anchor := core.ReviewRange{
		Path:  line.Path,
		Index: gui.diffLineSelected,
		Side:  string(line.Side),
	}
	if line.NewLine > 0 {
		anchor.Line = line.NewLine
	} else {
		anchor.Line = line.OldLine
	}
	gui.state.MarkReviewRangeStart(anchor)
	gui.state.SetReviewNotice("Range selection started.")
	gui.focus = panelDiffContent
}

func (gui *Gui) beginReviewCommentFlow() {
	if gui.state.Review.RangeStart == nil {
		gui.beginReviewCommentInput()
		return
	}
	gui.beginReviewCommentInput()
}

func (gui *Gui) buildReviewCommentDraft(body string) (gh.ReviewComment, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return gh.ReviewComment{}, fmt.Errorf("comment body is empty")
	}
	line, ok := gui.currentDiffLine()
	if !ok || !line.Commentable {
		return gh.ReviewComment{}, fmt.Errorf("current line is not commentable")
	}
	comment := gh.ReviewComment{
		Path: line.Path,
		Body: body,
		Side: line.Side,
	}
	if line.NewLine > 0 && line.Side != gh.DiffSideLeft {
		comment.Line = line.NewLine
	} else {
		comment.Line = line.OldLine
	}
	if comment.Line <= 0 {
		return gh.ReviewComment{}, fmt.Errorf("comment line is invalid")
	}

	start := gui.state.Review.RangeStart
	if start == nil {
		return comment, nil
	}
	if start.Path != comment.Path {
		return gh.ReviewComment{}, fmt.Errorf("range must stay within one file")
	}
	if start.Index != gui.diffLineSelected {
		comment.StartLine = start.Line
		comment.StartSide = gh.DiffSide(start.Side)
		if start.Index > gui.diffLineSelected {
			comment.StartLine, comment.Line = comment.Line, comment.StartLine
			comment.StartSide, comment.Side = comment.Side, comment.StartSide
		}
	}
	return comment, nil
}

func (gui *Gui) isDiffLineWithinPendingRange(line gh.DiffLine) bool {
	start := gui.state.Review.RangeStart
	if start == nil {
		return false
	}
	if start.Path != line.Path || !line.Commentable {
		return false
	}
	file, ok := gui.currentDiffFile()
	if !ok {
		return false
	}
	lineIndex := -1
	for idx, candidate := range file.Lines {
		if candidate == line {
			lineIndex = idx
			break
		}
	}
	if lineIndex < 0 {
		return false
	}
	minIndex := start.Index
	maxIndex := gui.diffLineSelected
	if minIndex > maxIndex {
		minIndex, maxIndex = maxIndex, minIndex
	}
	return lineIndex >= minIndex && lineIndex <= maxIndex
}
