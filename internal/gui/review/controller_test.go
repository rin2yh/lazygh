package review

import (
	"errors"
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/model"
	appstate "github.com/rin2yh/lazygh/internal/state"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
	reviewstub "github.com/rin2yh/lazygh/pkg/test/stub/review"
)

func defaultTestConfig() *config.Config { return config.Default() }

func setupControllerWithPR(client *testmock.GHClient, sel reviewstub.Selection) (*Controller, *appstate.State, *FocusTarget) {
	state := appstate.NewState()
	state.ApplyPRsResult("owner/repo", []model.Item{{Number: 1, Title: "PR"}}, nil)
	focus := FocusDiffContent
	c := NewController(defaultTestConfig(), state, client, sel, func(f FocusTarget) { focus = f })
	return c, state, &focus
}

func TestHandleCommentSave_NoPRSelected(t *testing.T) {
	state := appstate.NewState()
	c := NewController(defaultTestConfig(), state, &testmock.GHClient{}, reviewstub.Selection{}, func(FocusTarget) {})

	cmd := c.HandleCommentSave()
	if cmd != nil {
		t.Fatal("expected nil cmd when no PR selected")
	}
	if state.Review.Notice != "No pull request selected." {
		t.Fatalf("got %q, want %q", state.Review.Notice, "No pull request selected.")
	}
}

func TestHandleCommentSave_BuildDraftError(t *testing.T) {
	sel := reviewstub.Selection{
		Line: gh.DiffLine{Path: "a.go", Commentable: true, Side: gh.DiffSideRight},
	}
	c, state, _ := setupControllerWithPR(&testmock.GHClient{}, sel)
	// Set empty comment value to trigger "comment body is empty"
	c.SetCommentValue("")

	cmd := c.HandleCommentSave()
	if cmd != nil {
		t.Fatal("expected nil cmd on draft error")
	}
	if state.Review.Notice != "comment body is empty" {
		t.Fatalf("got notice %q, want %q", state.Review.Notice, "comment body is empty")
	}
}

func TestHandleCommentSave_InvalidLine(t *testing.T) {
	sel := reviewstub.Selection{
		Line: gh.DiffLine{Path: "a.go", Commentable: false},
	}
	c, state, _ := setupControllerWithPR(&testmock.GHClient{}, sel)
	c.SetCommentValue("hello")

	cmd := c.HandleCommentSave()
	if cmd != nil {
		t.Fatal("expected nil cmd on non-commentable line")
	}
	if state.Review.Notice != "current line is not commentable" {
		t.Fatalf("got notice %q, want %q", state.Review.Notice, "current line is not commentable")
	}
}

func TestHandleSubmit_NoPendingReview(t *testing.T) {
	state := appstate.NewState()
	c := NewController(defaultTestConfig(), state, &testmock.GHClient{}, reviewstub.Selection{}, func(FocusTarget) {})

	cmd := c.HandleSubmit()
	if cmd != nil {
		t.Fatal("expected nil cmd when no pending review")
	}
	if state.Review.Notice != "No pending review to submit." {
		t.Fatalf("got %q, want %q", state.Review.Notice, "No pending review to submit.")
	}
}

func TestApplySubmitResult_ErrorPreservesState(t *testing.T) {
	c, state, _ := setupControllerWithPR(&testmock.GHClient{}, reviewstub.Selection{})
	state.BeginReviewLoad()
	state.SetReviewContext(1, "PR_1", "abc123", "PRR_1")
	state.AddReviewComment(model.ReviewComment{Path: "a.go", Body: "hi", Line: 10})

	c.ApplySubmitResult(SubmittedMsg{ReviewID: "PRR_1", Err: errors.New("network error")})

	if state.Review.Notice != "network error" {
		t.Fatalf("got %q, want %q", state.Review.Notice, "network error")
	}
	if state.Review.ReviewID != "PRR_1" {
		t.Fatalf("review ID cleared on error, got %q", state.Review.ReviewID)
	}
	if len(state.Review.Comments) != 1 {
		t.Fatalf("comments cleared on error, got %d", len(state.Review.Comments))
	}
}

func TestApplySubmitResult_SuccessClearsReview(t *testing.T) {
	c, state, focus := setupControllerWithPR(&testmock.GHClient{}, reviewstub.Selection{})
	state.BeginReviewLoad()
	state.SetReviewContext(1, "PR_1", "abc123", "PRR_1")
	state.AddReviewComment(model.ReviewComment{Path: "a.go", Body: "hi", Line: 10})

	c.ApplySubmitResult(SubmittedMsg{ReviewID: "PRR_1"})

	if state.Review.ReviewID != "" {
		t.Fatalf("review ID not cleared, got %q", state.Review.ReviewID)
	}
	if len(state.Review.Comments) != 0 {
		t.Fatalf("comments not cleared, got %d", len(state.Review.Comments))
	}
	if *focus != FocusDiffContent {
		t.Fatalf("focus not restored to diff content, got %v", *focus)
	}
}

func TestApplyDiscardResult_ErrorPreservesState(t *testing.T) {
	c, state, _ := setupControllerWithPR(&testmock.GHClient{}, reviewstub.Selection{})
	state.BeginReviewLoad()
	state.SetReviewContext(1, "PR_1", "abc123", "PRR_1")
	state.AddReviewComment(model.ReviewComment{Path: "a.go", Body: "hi", Line: 10})

	c.ApplyDiscardResult(DiscardedMsg{Err: errors.New("discard failed")})

	if state.Review.Notice != "discard failed" {
		t.Fatalf("got %q, want %q", state.Review.Notice, "discard failed")
	}
	if state.Review.ReviewID != "PRR_1" {
		t.Fatalf("review ID cleared on error, got %q", state.Review.ReviewID)
	}
}

func TestApplyDiscardResult_SuccessClearsReview(t *testing.T) {
	c, state, focus := setupControllerWithPR(&testmock.GHClient{}, reviewstub.Selection{})
	state.BeginReviewLoad()
	state.SetReviewContext(1, "PR_1", "abc123", "PRR_1")
	state.AddReviewComment(model.ReviewComment{Path: "a.go", Body: "hi", Line: 10})
	state.SetReviewSummary("my summary")

	c.ApplyDiscardResult(DiscardedMsg{})

	if state.Review.ReviewID != "" {
		t.Fatalf("review ID not cleared, got %q", state.Review.ReviewID)
	}
	if len(state.Review.Comments) != 0 {
		t.Fatalf("comments not cleared, got %d", len(state.Review.Comments))
	}
	if c.CurrentSummaryValue() != "" {
		t.Fatalf("summary editor not cleared, got %q", c.CurrentSummaryValue())
	}
	if *focus != FocusDiffContent {
		t.Fatalf("focus not restored, got %v", *focus)
	}
}

func TestApplyCommentResult_SuccessAddsComment(t *testing.T) {
	c, state, _ := setupControllerWithPR(&testmock.GHClient{}, reviewstub.Selection{})
	state.BeginReviewLoad()

	c.ApplyCommentResult(CommentSavedMsg{
		PRNumber: 1,
		Context:  gh.ReviewContext{PullRequestID: "PR_1", CommitOID: "abc"},
		ReviewID: "PRR_1",
		Comment:  gh.ReviewComment{Path: "a.go", Body: "nice", Line: 5, Side: gh.DiffSideRight},
	})

	if len(state.Review.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(state.Review.Comments))
	}
	if state.Review.Comments[0].Body != "nice" {
		t.Fatalf("got body %q, want %q", state.Review.Comments[0].Body, "nice")
	}
	if state.Review.ReviewID != "PRR_1" {
		t.Fatalf("got review ID %q, want %q", state.Review.ReviewID, "PRR_1")
	}
}

func TestHandleSubmit_SavesSummaryIfInSummaryMode(t *testing.T) {
	c, state, _ := setupControllerWithPR(&testmock.GHClient{}, reviewstub.Selection{})
	state.SetReviewContext(1, "PR_1", "abc", "PRR_1")
	state.AddReviewComment(model.ReviewComment{Path: "a.go", Body: "hi", Line: 1})
	state.BeginReviewSummaryInput()

	// Set summary text in the editor via controller
	c.summary.editor.SetValue("my summary text")

	cmd := c.HandleSubmit()
	// Should produce a command (async submit)
	if cmd == nil {
		t.Fatal("expected non-nil cmd for submit")
	}
	if state.Review.Summary != "my summary text" {
		t.Fatalf("summary not saved, got %q", state.Review.Summary)
	}
	if state.Review.InputMode != model.ReviewInputNone {
		t.Fatalf("input mode not cleared, got %v", state.Review.InputMode)
	}
}
