package review

import "github.com/rin2yh/lazygh/internal/model"

// ReviewState holds all mutable state for the pending review workflow.
type ReviewState struct {
	PRNumber           int
	PullRequestID      string
	CommitOID          string
	ReviewID           string
	DrawerOpen         bool
	InputMode          InputMode
	Event              Event
	Summary            string
	Comments           []Comment
	RangeStart         *Range
	Notice             string
	SelectedCommentIdx int
	EditingCommentIdx  int
}

func newReviewState() *ReviewState {
	return &ReviewState{
		Comments:          []Comment{},
		EditingCommentIdx: noEditingComment,
	}
}

func (rs *ReviewState) HasPendingReview() bool {
	return rs.ReviewID != ""
}

func (rs *ReviewState) StopInput() {
	rs.InputMode = InputNone
}

func (rs *ReviewState) OpenDrawer() {
	rs.DrawerOpen = true
}

func (rs *ReviewState) CloseDrawer() {
	rs.DrawerOpen = false
	rs.InputMode = InputNone
	rs.Notice = ""
}

func (rs *ReviewState) BeginCommentInput() {
	rs.DrawerOpen = true
	rs.InputMode = InputComment
	rs.Notice = ""
}

func (rs *ReviewState) BeginSummaryInput() {
	rs.DrawerOpen = true
	rs.InputMode = InputSummary
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

func (rs *ReviewState) AddComment(comment Comment) {
	rs.Comments = append(rs.Comments, Comment{
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
	rs.InputMode = InputNone
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

func (rs *ReviewState) DeleteSelectedComment() (Comment, bool) {
	idx := rs.SelectedCommentIdx
	if idx < 0 || idx >= len(rs.Comments) {
		return Comment{}, false
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

func (rs *ReviewState) SelectedComment() (Comment, bool) {
	idx := rs.SelectedCommentIdx
	if idx < 0 || idx >= len(rs.Comments) {
		return Comment{}, false
	}
	return rs.Comments[idx], true
}

func (rs *ReviewState) BeginEditComment() {
	rs.EditingCommentIdx = rs.SelectedCommentIdx
	rs.InputMode = InputComment
	rs.DrawerOpen = true
	rs.Notice = ""
}

func (rs *ReviewState) ApplyEditComment(newBody string) {
	idx := rs.EditingCommentIdx
	if idx < 0 || idx >= len(rs.Comments) {
		return
	}
	rs.Comments[idx].Body = model.SanitizeMultiline(newBody)
	rs.EditingCommentIdx = noEditingComment
	rs.InputMode = InputNone
	rs.Notice = "Comment updated."
}

func (rs *ReviewState) ClearEditingComment() {
	rs.EditingCommentIdx = noEditingComment
}

func (rs *ReviewState) SetNotice(msg string) {
	rs.Notice = model.SanitizeMultiline(msg)
}

func (rs *ReviewState) ClearNotice() {
	rs.Notice = ""
}

func (rs *ReviewState) MarkRangeStart(anchor Range) {
	copied := anchor
	rs.RangeStart = &copied
	rs.DrawerOpen = true
	rs.Notice = "Range start selected."
}

func (rs *ReviewState) CycleEvent() {
	switch rs.Event {
	case EventComment:
		rs.Event = EventApprove
	case EventApprove:
		rs.Event = EventRequestChanges
	default:
		rs.Event = EventComment
	}
}

func (rs *ReviewState) ClearRangeStart() {
	rs.RangeStart = nil
}

// Reset clears the review state entirely (e.g. when PR list reloads).
func (rs *ReviewState) Reset() {
	rs.reset()
}

func (rs *ReviewState) reset() {
	*rs = ReviewState{
		EditingCommentIdx: noEditingComment,
		Comments:          []Comment{},
	}
}
