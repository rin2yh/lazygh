package review

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/gh"
)

// Message types returned by async commands.

type CommentSavedMsg struct {
	PRNumber  int
	Context   gh.ReviewContext
	ReviewID  string
	Comment   gh.ReviewComment
	CommentID string
	Err       error
}

type CommentDeletedMsg struct {
	CommentID string
	Err       error
}

type CommentUpdatedMsg struct {
	Body string
	Err  error
}

type SubmittedMsg struct {
	ReviewID string
	Err      error
}

type DiscardedMsg struct {
	Err error
}

type pending struct {
	rs        *ReviewState
	host      AppState
	client    PendingReviewClient
	selection Selection
	comment   *comment
	summary   *summary
}

func newPending(rs *ReviewState, host AppState, client PendingReviewClient, selection Selection, comment *comment, summary *summary) *pending {
	comment.bindSelection(selection)
	return &pending{
		rs:        rs,
		host:      host,
		client:    client,
		selection: selection,
		comment:   comment,
		summary:   summary,
	}
}

func (f *pending) HandleCommentSave() tea.Cmd {
	item, ok := f.host.SelectedPR()
	if !ok {
		f.rs.SetNotice("No pull request selected.")
		return nil
	}
	comment, err := f.comment.BuildDraft(f.comment.CurrentValue(), f.rs.RangeStart)
	if err != nil {
		f.rs.SetNotice(err.Error())
		return nil
	}
	repo := f.host.ListRepo()
	reviewID := f.rs.ReviewID
	ctx := gh.ReviewContext{
		PullRequestID: f.rs.PullRequestID,
		CommitOID:     f.rs.CommitOID,
	}

	f.host.BeginFetchReview()
	return func() tea.Msg {
		var runErr error
		if reviewID == "" {
			ctx, runErr = f.client.GetReviewContext(repo, item.Number)
			if runErr != nil {
				return CommentSavedMsg{Err: runErr}
			}
			reviewID, runErr = f.client.StartPendingReview(repo, item.Number, ctx)
			if runErr != nil {
				return CommentSavedMsg{Err: runErr}
			}
		}
		commentID, runErr := f.client.AddReviewComment(repo, reviewID, comment)
		return CommentSavedMsg{
			PRNumber:  item.Number,
			Context:   ctx,
			ReviewID:  reviewID,
			Comment:   comment,
			CommentID: commentID,
			Err:       runErr,
		}
	}
}

func (f *pending) BeginEditComment() bool {
	comment, ok := f.rs.SelectedComment()
	if !ok {
		return false
	}
	f.rs.BeginEditComment()
	f.comment.StartEdit(comment.Body)
	return true
}

func (f *pending) IsEditingComment() bool {
	return f.rs.EditingCommentIdx != noEditingComment
}

func (f *pending) SelectNextComment() {
	f.rs.SelectNextComment()
}

func (f *pending) SelectPrevComment() {
	f.rs.SelectPrevComment()
}

func (f *pending) HandleDeleteComment() tea.Cmd {
	comment, ok := f.rs.SelectedComment()
	if !ok {
		f.rs.SetNotice("No comment selected.")
		return nil
	}
	if comment.CommentID == "" {
		f.rs.DeleteSelectedComment()
		f.rs.SetNotice("Comment deleted.")
		return nil
	}
	commentID := comment.CommentID
	f.host.BeginFetchReview()
	return func() tea.Msg {
		err := f.client.DeletePendingReviewComment(commentID)
		return CommentDeletedMsg{CommentID: commentID, Err: err}
	}
}

func (f *pending) HandleEditCommentSave() tea.Cmd {
	idx := f.rs.EditingCommentIdx
	if idx < 0 || idx >= len(f.rs.Comments) {
		f.rs.SetNotice("No comment being edited.")
		return nil
	}
	body := strings.TrimSpace(f.comment.CurrentValue())
	if body == "" {
		f.rs.SetNotice("Comment body is empty.")
		return nil
	}
	comment := f.rs.Comments[idx]
	if comment.CommentID == "" {
		f.rs.ApplyEditComment(body)
		f.comment.StopInput()
		return nil
	}
	commentID := comment.CommentID
	f.host.BeginFetchReview()
	return func() tea.Msg {
		err := f.client.UpdatePendingReviewComment(commentID, body)
		return CommentUpdatedMsg{Body: body, Err: err}
	}
}

func (f *pending) HandleSubmit() tea.Cmd {
	if f.rs.InputMode == InputSummary {
		f.summary.Save()
		f.summary.StopInput()
		f.rs.StopInput()
	}
	if !f.rs.HasPendingReview() {
		f.rs.SetNotice("No pending review to submit.")
		return nil
	}
	f.host.BeginFetchReview()
	reviewID := f.rs.ReviewID
	body := f.rs.Summary
	repo := f.host.ListRepo()
	event := coreEventToGH(f.rs.Event)
	return func() tea.Msg {
		err := f.client.SubmitReview(repo, reviewID, event, body)
		return SubmittedMsg{ReviewID: reviewID, Err: err}
	}
}

func coreEventToGH(e Event) gh.ReviewEvent {
	switch e {
	case EventApprove:
		return gh.ReviewEventApprove
	case EventRequestChanges:
		return gh.ReviewEventRequestChanges
	default:
		return gh.ReviewEventComment
	}
}

func (f *pending) HandleDiscard() tea.Cmd {
	if f.rs.InputMode == InputSummary {
		f.summary.StopInput()
		f.rs.StopInput()
	}
	reviewID := f.rs.ReviewID
	if reviewID == "" {
		f.rs.ResetAfterDiscard("Review draft discarded.")
		return nil
	}
	f.host.BeginFetchReview()
	repo := f.host.ListRepo()
	return func() tea.Msg {
		err := f.client.DeletePendingReview(repo, reviewID)
		return DiscardedMsg{Err: err}
	}
}

// ApplyCommentResult applies the result of a comment save operation and returns
// true if the comment was saved successfully (caller should set focus to FocusReviewDrawer).
func (f *pending) ApplyCommentResult(msg CommentSavedMsg) bool {
	f.host.ClearFetching()
	if msg.ReviewID != "" || msg.Context.PullRequestID != "" || msg.Context.CommitOID != "" {
		f.rs.SetContext(msg.PRNumber, msg.Context.PullRequestID, msg.Context.CommitOID, msg.ReviewID)
	}
	if msg.Err != nil {
		f.rs.SetNotice(msg.Err.Error())
		return false
	}
	f.rs.AddComment(Comment{
		CommentID: msg.CommentID,
		Path:      msg.Comment.Path,
		Body:      msg.Comment.Body,
		Side:      string(msg.Comment.Side),
		Line:      msg.Comment.Line,
		StartSide: string(msg.Comment.StartSide),
		StartLine: msg.Comment.StartLine,
	})
	f.comment.ApplySaved()
	return true
}

func (f *pending) ApplyDeleteCommentResult(msg CommentDeletedMsg) {
	f.host.ClearFetching()
	if msg.Err != nil {
		f.rs.SetNotice(msg.Err.Error())
		return
	}
	f.rs.DeleteSelectedComment()
	f.rs.SetNotice("Comment deleted.")
}

func (f *pending) ApplyEditCommentResult(msg CommentUpdatedMsg) {
	f.host.ClearFetching()
	if msg.Err != nil {
		f.rs.SetNotice(msg.Err.Error())
		return
	}
	f.rs.ApplyEditComment(msg.Body)
	f.comment.StopInput()
}

func (f *pending) ApplySubmitResult(msg SubmittedMsg) {
	f.host.ClearFetching()
	if msg.Err != nil {
		f.rs.SetNotice(msg.Err.Error())
		return
	}
	f.comment.StopInput()
	f.summary.StopInput()
	f.rs.ResetAfterSubmit("Review submitted.")
}

func (f *pending) ApplyDiscardResult(msg DiscardedMsg) {
	f.host.ClearFetching()
	if msg.Err != nil {
		f.rs.SetNotice(msg.Err.Error())
		return
	}
	f.comment.StopInput()
	f.summary.StopInput()
	f.summary.Clear()
	f.rs.ResetAfterDiscard("Review draft discarded.")
}
