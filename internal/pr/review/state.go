package review

import (
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/pkg/sanitize"
)

const noEditingComment = -1

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

	Threads           []gh.ReviewThread
	SelectedThreadIdx int
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
	rs.ClearNotice()
}

func (rs *ReviewState) BeginCommentInput() {
	rs.DrawerOpen = true
	rs.InputMode = InputComment
	rs.ClearNotice()
}

func (rs *ReviewState) BeginSummaryInput() {
	rs.DrawerOpen = true
	rs.InputMode = InputSummary
	rs.ClearNotice()
}

func (rs *ReviewState) SetSummary(summary string) {
	rs.Summary = sanitize.Multiline(summary)
}

func (rs *ReviewState) SetContext(prNumber int, pullRequestID string, commitOID string, reviewID string) {
	rs.PRNumber = prNumber
	rs.PullRequestID = sanitize.SingleLine(pullRequestID)
	rs.CommitOID = sanitize.SingleLine(commitOID)
	rs.ReviewID = sanitize.SingleLine(reviewID)
}

func (rs *ReviewState) AddComment(comment Comment) {
	rs.Comments = append(rs.Comments, Comment{
		CommentID: comment.CommentID,
		Path:      sanitize.SingleLine(comment.Path),
		Body:      sanitize.Multiline(comment.Body),
		Side:      sanitize.SingleLine(comment.Side),
		Line:      comment.Line,
		StartSide: sanitize.SingleLine(comment.StartSide),
		StartLine: comment.StartLine,
	})
	rs.SelectedCommentIdx = len(rs.Comments) - 1
	rs.Notify("Review comment added.")
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
	rs.ClearNotice()
}

func (rs *ReviewState) ApplyEditComment(newBody string) {
	idx := rs.EditingCommentIdx
	if idx < 0 || idx >= len(rs.Comments) {
		return
	}
	rs.Comments[idx].Body = sanitize.Multiline(newBody)
	rs.EditingCommentIdx = noEditingComment
	rs.InputMode = InputNone
	rs.Notify("Comment updated.")
}

func (rs *ReviewState) ClearEditingComment() {
	rs.EditingCommentIdx = noEditingComment
}

func (rs *ReviewState) Notify(msg string) {
	rs.Notice = sanitize.Multiline(msg)
}

func (rs *ReviewState) ClearNotice() {
	rs.Notice = ""
}

func (rs *ReviewState) MarkRangeStart(anchor Range) {
	copied := anchor
	rs.RangeStart = &copied
	rs.DrawerOpen = true
	rs.Notify("Range start selected.")
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

// LoadThreads stores fetched review threads.
func (rs *ReviewState) LoadThreads(threads []gh.ReviewThread) {
	rs.Threads = threads
	rs.SelectedThreadIdx = 0
}

// SelectNextThread moves the thread selection down.
func (rs *ReviewState) SelectNextThread() {
	if len(rs.Threads) == 0 {
		return
	}
	if rs.SelectedThreadIdx < len(rs.Threads)-1 {
		rs.SelectedThreadIdx++
	}
}

// SelectPrevThread moves the thread selection up.
func (rs *ReviewState) SelectPrevThread() {
	if rs.SelectedThreadIdx > 0 {
		rs.SelectedThreadIdx--
	}
}

// SelectedThread returns the currently selected review thread.
func (rs *ReviewState) SelectedThread() (gh.ReviewThread, bool) {
	idx := rs.SelectedThreadIdx
	if idx < 0 || idx >= len(rs.Threads) {
		return gh.ReviewThread{}, false
	}
	return rs.Threads[idx], true
}

// BeginThreadReplyInput switches to thread-reply input mode.
func (rs *ReviewState) BeginThreadReplyInput() {
	rs.DrawerOpen = true
	rs.InputMode = InputThreadReply
	rs.ClearNotice()
}

// Reset clears the review state entirely (e.g. when PR list reloads).
func (rs *ReviewState) Reset() {
	*rs = ReviewState{
		EditingCommentIdx: noEditingComment,
		Comments:          []Comment{},
	}
}
