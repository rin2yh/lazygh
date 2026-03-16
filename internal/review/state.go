package review

import "github.com/rin2yh/lazygh/internal/model"

// ReviewState holds all mutable state for the pending review workflow.
type ReviewState struct {
	PRNumber           int
	PullRequestID      string
	CommitOID          string
	ReviewID           string
	DrawerOpen         bool
	InputMode          model.ReviewInputMode
	Event              model.ReviewEvent
	Summary            string
	Comments           []model.ReviewComment
	RangeStart         *model.ReviewRange
	Notice             string
	SelectedCommentIdx int
	EditingCommentIdx  int
}

func newReviewState() *ReviewState {
	return &ReviewState{
		Comments:          []model.ReviewComment{},
		EditingCommentIdx: model.NoEditingComment,
	}
}

func (rs *ReviewState) HasPendingReview() bool {
	return rs.ReviewID != ""
}

func (rs *ReviewState) StopInput() {
	rs.InputMode = model.ReviewInputNone
}

func (rs *ReviewState) OpenDrawer() {
	rs.DrawerOpen = true
}

func (rs *ReviewState) CloseDrawer() {
	rs.DrawerOpen = false
	rs.InputMode = model.ReviewInputNone
	rs.Notice = ""
}

func (rs *ReviewState) BeginCommentInput() {
	rs.DrawerOpen = true
	rs.InputMode = model.ReviewInputComment
	rs.Notice = ""
}

func (rs *ReviewState) BeginSummaryInput() {
	rs.DrawerOpen = true
	rs.InputMode = model.ReviewInputSummary
	rs.Notice = ""
}

func (rs *ReviewState) SetSummary(summary string) {
	rs.Summary = model.SanitizeMultiline(summary)
}

func (rs *ReviewState) SetContext(prNumber int, pullRequestID string, commitOID string, reviewID string) {
	rs.PRNumber = prNumber
	rs.PullRequestID = model.SanitizeSingleLine(pullRequestID)
	rs.CommitOID = model.SanitizeSingleLine(commitOID)
	rs.ReviewID = model.SanitizeSingleLine(reviewID)
}

func (rs *ReviewState) AddComment(comment model.ReviewComment) {
	rs.Comments = append(rs.Comments, model.ReviewComment{
		CommentID: comment.CommentID,
		Path:      model.SanitizeSingleLine(comment.Path),
		Body:      model.SanitizeMultiline(comment.Body),
		Side:      model.SanitizeSingleLine(comment.Side),
		Line:      comment.Line,
		StartSide: model.SanitizeSingleLine(comment.StartSide),
		StartLine: comment.StartLine,
	})
	rs.SelectedCommentIdx = len(rs.Comments) - 1
	rs.Notice = "Review comment added."
	rs.DrawerOpen = true
	rs.InputMode = model.ReviewInputNone
	rs.RangeStart = nil
}

func (rs *ReviewState) SelectNextComment() {
	if len(rs.Comments) == 0 {
		return
	}
	if rs.SelectedCommentIdx < len(rs.Comments)-1 {
		rs.SelectedCommentIdx++
	}
}

func (rs *ReviewState) SelectPrevComment() {
	if rs.SelectedCommentIdx > 0 {
		rs.SelectedCommentIdx--
	}
}

func (rs *ReviewState) DeleteSelectedComment() (model.ReviewComment, bool) {
	idx := rs.SelectedCommentIdx
	if idx < 0 || idx >= len(rs.Comments) {
		return model.ReviewComment{}, false
	}
	deleted := rs.Comments[idx]
	rs.Comments = append(rs.Comments[:idx], rs.Comments[idx+1:]...)
	if len(rs.Comments) == 0 {
		rs.SelectedCommentIdx = 0
	} else if rs.SelectedCommentIdx >= len(rs.Comments) {
		rs.SelectedCommentIdx = len(rs.Comments) - 1
	}
	return deleted, true
}

func (rs *ReviewState) SelectedComment() (model.ReviewComment, bool) {
	idx := rs.SelectedCommentIdx
	if idx < 0 || idx >= len(rs.Comments) {
		return model.ReviewComment{}, false
	}
	return rs.Comments[idx], true
}

func (rs *ReviewState) BeginEditComment() {
	rs.EditingCommentIdx = rs.SelectedCommentIdx
	rs.InputMode = model.ReviewInputComment
	rs.DrawerOpen = true
	rs.Notice = ""
}

func (rs *ReviewState) ApplyEditComment(newBody string) {
	idx := rs.EditingCommentIdx
	if idx < 0 || idx >= len(rs.Comments) {
		return
	}
	rs.Comments[idx].Body = model.SanitizeMultiline(newBody)
	rs.EditingCommentIdx = model.NoEditingComment
	rs.InputMode = model.ReviewInputNone
	rs.Notice = "Comment updated."
}

func (rs *ReviewState) ClearEditingComment() {
	rs.EditingCommentIdx = model.NoEditingComment
}

func (rs *ReviewState) SetNotice(msg string) {
	rs.Notice = model.SanitizeMultiline(msg)
}

func (rs *ReviewState) ClearNotice() {
	rs.Notice = ""
}

func (rs *ReviewState) MarkRangeStart(anchor model.ReviewRange) {
	copied := anchor
	rs.RangeStart = &copied
	rs.DrawerOpen = true
	rs.Notice = "Range start selected."
}

func (rs *ReviewState) CycleEvent() {
	switch rs.Event {
	case model.ReviewEventComment:
		rs.Event = model.ReviewEventApprove
	case model.ReviewEventApprove:
		rs.Event = model.ReviewEventRequestChanges
	default:
		rs.Event = model.ReviewEventComment
	}
}

func (rs *ReviewState) ClearRangeStart() {
	rs.RangeStart = nil
}

func (rs *ReviewState) ResetAfterSubmit(notice string) {
	rs.reset()
	rs.Notice = model.SanitizeMultiline(notice)
}

func (rs *ReviewState) ResetAfterDiscard(notice string) {
	rs.reset()
	rs.Notice = model.SanitizeMultiline(notice)
}

// Reset clears the review state entirely (e.g. when PR list reloads).
func (rs *ReviewState) Reset() {
	rs.reset()
}

func (rs *ReviewState) reset() {
	notice := rs.Notice
	*rs = ReviewState{
		Notice:            notice,
		EditingCommentIdx: model.NoEditingComment,
	}
	rs.Comments = []model.ReviewComment{}
}
