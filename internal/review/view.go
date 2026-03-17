package review

import "github.com/rin2yh/lazygh/internal/model"

type view struct {
	rs      *ReviewState
	host    AppState
	comment *comment
	summary *summary
}

func newView(rs *ReviewState, host AppState, comment *comment, summary *summary) *view {
	return &view{
		rs:      rs,
		host:    host,
		comment: comment,
		summary: summary,
	}
}

func (f *view) InputMode() model.ReviewInputMode {
	return f.rs.InputMode
}

func (f *view) ShouldShowDrawer() bool {
	if !f.host.IsDiffMode() {
		return false
	}
	rs := f.rs
	return rs.DrawerOpen || rs.InputMode != model.ReviewInputNone || rs.HasPendingReview() || len(rs.Comments) > 0 || rs.Summary != "" || rs.RangeStart != nil
}

// StopInput stops any active input and returns the FocusTarget to move to,
// or nil if focus should not change.
func (f *view) StopInput() *FocusTarget {
	f.comment.editor.Blur()
	f.summary.editor.Blur()
	if f.rs.InputMode == model.ReviewInputComment {
		f.rs.ClearRangeStart()
		f.rs.ClearEditingComment()
		f.comment.editor.SetValue("")
	}
	f.rs.StopInput()
	if f.ShouldShowDrawer() {
		t := FocusReviewDrawer
		return &t
	}
	return nil
}

func (f *view) HandleEsc() bool {
	f.StopInput()
	return true
}

// HandleSummarySave saves the summary, stops input, and returns the FocusTarget
// to move to (or nil if focus should not change).
func (f *view) HandleSummarySave() (bool, *FocusTarget) {
	f.summary.Save()
	target := f.StopInput()
	f.rs.SetNotice("Review summary updated.")
	return true, target
}
