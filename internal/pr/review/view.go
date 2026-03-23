package review

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

func (f *view) InputMode() InputMode {
	return f.rs.InputMode
}

func (f *view) ShouldShowDrawer() bool {
	if !f.host.IsDiffMode() {
		return false
	}
	rs := f.rs
	return rs.DrawerOpen || rs.InputMode != InputNone || rs.HasPendingReview() || len(rs.Comments) > 0 || rs.Summary != "" || rs.RangeStart != nil
}

// StopInput stops any active input and returns the FocusTarget to move to
// and whether focus should change.
func (f *view) StopInput() (FocusTarget, bool) {
	f.comment.Blur()
	f.summary.Blur()
	if f.rs.InputMode == InputComment {
		f.rs.ClearRangeStart()
		f.rs.ClearEditingComment()
		f.comment.State.Clear()
	}
	f.rs.StopInput()
	if f.ShouldShowDrawer() {
		return FocusReviewDrawer, true
	}
	return 0, false
}

func (f *view) HandleEsc() bool {
	f.StopInput()
	return true
}

// HandleSummarySave saves the summary, stops input, and returns the FocusTarget
// to move to and whether focus should change.
func (f *view) HandleSummarySave() (FocusTarget, bool) {
	f.summary.Save()
	target, ok := f.StopInput()
	f.rs.Notify("Review summary updated.")
	return target, ok
}
