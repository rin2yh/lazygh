package review

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
)

type CommentSavedMsg struct {
	PRNumber int
	Context  gh.ReviewContext
	ReviewID string
	Comment  gh.ReviewComment
	Err      error
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
	client    gh.ClientInterface
	selection Selection
	setFocus  func(FocusTarget)
	comment   *comment
	summary   *summary
}

func newPending(state *core.State, client gh.ClientInterface, selection Selection, setFocus func(FocusTarget), comment *comment, summary *summary) *pending {
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
	repo := f.state.Repo
	reviewID := f.state.Review.ReviewID
	ctx := gh.ReviewContext{
		PullRequestID: f.state.Review.PullRequestID,
		CommitOID:     f.state.Review.CommitOID,
	}

	f.state.Loading = core.LoadingReview
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
		runErr = f.client.AddReviewComment(repo, reviewID, comment)
		return CommentSavedMsg{
			PRNumber: item.Number,
			Context:  ctx,
			ReviewID: reviewID,
			Comment:  comment,
			Err:      runErr,
		}
	}
}

func (f *pending) HandleSubmit() tea.Cmd {
	if f.state.Review.InputMode == core.ReviewInputSummary {
		f.summary.Save()
		f.summary.StopInput()
		f.state.Review.InputMode = core.ReviewInputNone
	}
	if !f.state.HasPendingReview() {
		f.state.SetReviewNotice("No pending review to submit.")
		return nil
	}
	f.state.Loading = core.LoadingReview
	reviewID := f.state.Review.ReviewID
	body := f.state.Review.Summary
	repo := f.state.Repo
	return func() tea.Msg {
		err := f.client.SubmitReview(repo, reviewID, body)
		return SubmittedMsg{ReviewID: reviewID, Err: err}
	}
}

func (f *pending) HandleDiscard() tea.Cmd {
	if f.state.Review.InputMode == core.ReviewInputSummary {
		f.summary.StopInput()
		f.state.Review.InputMode = core.ReviewInputNone
	}
	reviewID := f.state.Review.ReviewID
	if reviewID == "" {
		f.state.ResetReviewAfterDiscard("Review draft discarded.")
		return nil
	}
	f.state.Loading = core.LoadingReview
	repo := f.state.Repo
	return func() tea.Msg {
		err := f.client.DeletePendingReview(repo, reviewID)
		return DiscardedMsg{Err: err}
	}
}

func (f *pending) ApplyCommentResult(msg CommentSavedMsg) {
	f.state.Loading = core.LoadingNone
	if msg.ReviewID != "" || msg.Context.PullRequestID != "" || msg.Context.CommitOID != "" {
		f.state.SetReviewContext(msg.PRNumber, msg.Context.PullRequestID, msg.Context.CommitOID, msg.ReviewID)
	}
	if msg.Err != nil {
		f.state.SetReviewNotice(msg.Err.Error())
		return
	}
	f.state.AddReviewComment(core.ReviewComment{
		Path:      msg.Comment.Path,
		Body:      msg.Comment.Body,
		Side:      string(msg.Comment.Side),
		Line:      msg.Comment.Line,
		StartSide: string(msg.Comment.StartSide),
		StartLine: msg.Comment.StartLine,
	})
	f.comment.ApplySaved()
}

func (f *pending) ApplySubmitResult(msg SubmittedMsg) {
	f.state.Loading = core.LoadingNone
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
	f.state.Loading = core.LoadingNone
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
