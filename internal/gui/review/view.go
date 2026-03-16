package review

import (
	"github.com/rin2yh/lazygh/internal/core"
	appstate "github.com/rin2yh/lazygh/internal/state"
)

type view struct {
	state    *appstate.State
	setFocus func(FocusTarget)
	comment  *comment
	summary  *summary
}

func newView(state *appstate.State, setFocus func(FocusTarget), comment *comment, summary *summary) *view {
	return &view{
		state:    state,
		setFocus: setFocus,
		comment:  comment,
		summary:  summary,
	}
}

func (f *view) InputMode() core.ReviewInputMode {
	return f.state.Review.InputMode
}

func (f *view) ShouldShowDrawer() bool {
	if !f.state.IsDiffMode() {
		return false
	}
	review := f.state.Review
	return review.DrawerOpen || review.InputMode != core.ReviewInputNone || f.state.HasPendingReview() || len(review.Comments) > 0 || review.Summary != "" || review.RangeStart != nil
}

func (f *view) StopInput() {
	f.comment.editor.Blur()
	f.summary.editor.Blur()
	if f.state.Review.InputMode == core.ReviewInputComment {
		f.state.ClearReviewRangeStart()
		f.state.ClearEditingComment()
		f.comment.editor.SetValue("")
	}
	f.state.StopReviewInput()
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
	f.state.SetReviewNotice("Review summary updated.")
	return true
}
