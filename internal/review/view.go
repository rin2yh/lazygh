package review

import "github.com/rin2yh/lazygh/internal/model"

type view struct {
	rs       *ReviewState
	host     AppState
	setFocus func(FocusTarget)
	comment  *comment
	summary  *summary
}

func newView(rs *ReviewState, host AppState, setFocus func(FocusTarget), comment *comment, summary *summary) *view {
	return &view{
		rs:       rs,
		host:     host,
		setFocus: setFocus,
		comment:  comment,
		summary:  summary,
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

func (f *view) StopInput() {
	f.comment.editor.Blur()
	f.summary.editor.Blur()
	if f.rs.InputMode == model.ReviewInputComment {
		f.rs.ClearRangeStart()
		f.rs.ClearEditingComment()
		f.comment.editor.SetValue("")
	}
	f.rs.StopInput()
	if f.ShouldShowDrawer() {
		f.setFocus(FocusReviewDrawer)
	}
}

func (f *view) HandleEsc() bool {
	f.StopInput()
	f.setFocus(FocusDiffContent)
	return true
}

func (f *view) HandleSummarySave() bool {
	f.summary.Save()
	f.StopInput()
	f.rs.SetNotice("Review summary updated.")
	return true
}
