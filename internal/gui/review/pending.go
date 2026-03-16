package review

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
)

type PendingReviewClient interface {
	GetReviewContext(repo string, number int) (gh.ReviewContext, error)
	StartPendingReview(repo string, number int, ctx gh.ReviewContext) (string, error)
	AddReviewComment(repo string, reviewID string, comment gh.ReviewComment) (string, error)
	SubmitReview(repo string, reviewID string, event gh.ReviewEvent, body string) error
	DeletePendingReview(repo string, reviewID string) error
	DeletePendingReviewComment(commentID string) error
	UpdatePendingReviewComment(commentID string, body string) error
}

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
	Idx  int
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
	state     *core.State
	client    PendingReviewClient
	selection Selection
	setFocus  func(FocusTarget)
	comment   *comment
	summary   *summary
}

func newPending(state *core.State, client PendingReviewClient, selection Selection, setFocus func(FocusTarget), comment *comment, summary *summary) *pending {
	comment.bindSelection(selection)
	return &pending{
		state:     state,
		client:    client,
		selection: selection,
		setFocus:  setFocus,
		comment:   comment,
		summary:   summary,
	}
}

func (f *pending) HandleCommentSave() tea.Cmd {
	item, ok := f.state.SelectedPR()
	if !ok {
		f.state.SetReviewNotice("No pull request selected.")
		return nil
	}
	comment, err := f.comment.BuildDraft(f.comment.CurrentValue(), f.state.Review.RangeStart)
	if err != nil {
		f.state.SetReviewNotice(err.Error())
		return nil
	}
	repo := f.state.List.Repo
	reviewID := f.state.Review.ReviewID
	ctx := gh.ReviewContext{
		PullRequestID: f.state.Review.PullRequestID,
		CommitOID:     f.state.Review.CommitOID,
	}

	f.state.BeginReviewLoad()
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

func (f *pending) HandleDeleteComment() tea.Cmd {
	comment, ok := f.state.DeleteSelectedComment()
	if !ok {
		f.state.SetReviewNotice("No comment selected.")
		return nil
	}
	if comment.CommentID == "" {
		f.state.SetReviewNotice("Comment deleted.")
		return nil
	}
	commentID := comment.CommentID
	f.state.BeginReviewLoad()
	return func() tea.Msg {
		err := f.client.DeletePendingReviewComment(commentID)
		return CommentDeletedMsg{CommentID: commentID, Err: err}
	}
}

func (f *pending) HandleEditCommentSave() tea.Cmd {
	idx := f.state.Review.EditingCommentIdx
	if idx < 0 || idx >= len(f.state.Review.Comments) {
		f.state.SetReviewNotice("No comment being edited.")
		return nil
	}
	body := strings.TrimSpace(f.comment.CurrentValue())
	if body == "" {
		f.state.SetReviewNotice("Comment body is empty.")
		return nil
	}
	comment := f.state.Review.Comments[idx]
	if comment.CommentID == "" {
		f.state.ApplyEditComment(body)
		f.comment.StopInput()
		return nil
	}
	commentID := comment.CommentID
	f.state.BeginReviewLoad()
	return func() tea.Msg {
		err := f.client.UpdatePendingReviewComment(commentID, body)
		return CommentUpdatedMsg{Idx: idx, Body: body, Err: err}
	}
}

func (f *pending) HandleSubmit() tea.Cmd {
	if f.state.Review.InputMode == core.ReviewInputSummary {
		f.summary.Save()
		f.summary.StopInput()
		f.state.StopReviewInput()
	}
	if !f.state.HasPendingReview() {
		f.state.SetReviewNotice("No pending review to submit.")
		return nil
	}
	f.state.BeginReviewLoad()
	reviewID := f.state.Review.ReviewID
	body := f.state.Review.Summary
	repo := f.state.List.Repo
	event := coreEventToGH(f.state.Review.Event)
	return func() tea.Msg {
		err := f.client.SubmitReview(repo, reviewID, event, body)
		return SubmittedMsg{ReviewID: reviewID, Err: err}
	}
}

func coreEventToGH(e core.ReviewEvent) gh.ReviewEvent {
	switch e {
	case core.ReviewEventApprove:
		return gh.ReviewEventApprove
	case core.ReviewEventRequestChanges:
		return gh.ReviewEventRequestChanges
	default:
		return gh.ReviewEventComment
	}
}

func (f *pending) HandleDiscard() tea.Cmd {
	if f.state.Review.InputMode == core.ReviewInputSummary {
		f.summary.StopInput()
		f.state.StopReviewInput()
	}
	reviewID := f.state.Review.ReviewID
	if reviewID == "" {
		f.state.ResetReviewAfterDiscard("Review draft discarded.")
		return nil
	}
	f.state.BeginReviewLoad()
	repo := f.state.List.Repo
	return func() tea.Msg {
		err := f.client.DeletePendingReview(repo, reviewID)
		return DiscardedMsg{Err: err}
	}
}

func (f *pending) ApplyCommentResult(msg CommentSavedMsg) {
	f.state.ClearLoading()
	if msg.ReviewID != "" || msg.Context.PullRequestID != "" || msg.Context.CommitOID != "" {
		f.state.SetReviewContext(msg.PRNumber, msg.Context.PullRequestID, msg.Context.CommitOID, msg.ReviewID)
	}
	if msg.Err != nil {
		f.state.SetReviewNotice(msg.Err.Error())
		return
	}
	f.state.AddReviewComment(core.ReviewComment{
		CommentID: msg.CommentID,
		Path:      msg.Comment.Path,
		Body:      msg.Comment.Body,
		Side:      string(msg.Comment.Side),
		Line:      msg.Comment.Line,
		StartSide: string(msg.Comment.StartSide),
		StartLine: msg.Comment.StartLine,
	})
	f.comment.ApplySaved()
}

func (f *pending) ApplyDeleteCommentResult(msg CommentDeletedMsg) {
	f.state.ClearLoading()
	if msg.Err != nil {
		f.state.SetReviewNotice(msg.Err.Error())
		return
	}
	f.state.SetReviewNotice("Comment deleted.")
}

func (f *pending) ApplyEditCommentResult(msg CommentUpdatedMsg) {
	f.state.ClearLoading()
	if msg.Err != nil {
		f.state.SetReviewNotice(msg.Err.Error())
		return
	}
	f.state.ApplyEditComment(msg.Body)
	f.comment.StopInput()
	f.setFocus(FocusReviewDrawer)
}

func (f *pending) ApplySubmitResult(msg SubmittedMsg) {
	f.state.ClearLoading()
	if msg.Err != nil {
		f.state.SetReviewNotice(msg.Err.Error())
		return
	}
	f.comment.StopInput()
	f.summary.StopInput()
	f.state.ResetReviewAfterSubmit("Review submitted.")
	f.setFocus(FocusDiffContent)
}

func (f *pending) ApplyDiscardResult(msg DiscardedMsg) {
	f.state.ClearLoading()
	if msg.Err != nil {
		f.state.SetReviewNotice(msg.Err.Error())
		return
	}
	f.comment.StopInput()
	f.summary.StopInput()
	f.summary.Clear()
	f.state.ResetReviewAfterDiscard("Review draft discarded.")
	f.setFocus(FocusDiffContent)
}
