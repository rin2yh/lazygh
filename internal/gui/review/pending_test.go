package review

import (
	"errors"
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	appstate "github.com/rin2yh/lazygh/internal/state"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
	reviewstub "github.com/rin2yh/lazygh/pkg/test/stub/review"
)

func TestApplyCommentResult_PersistsPendingReviewContextOnError(t *testing.T) {
	state := appstate.NewState()
	state.ApplyPRsResult("owner/repo", []core.Item{{Number: 1, Title: "Fix bug"}}, nil)
	state.BeginReviewLoad()
	controller := NewController(config.Default(), state, &testmock.GHClient{}, reviewstub.Selection{}, func(FocusTarget) {})

	controller.ApplyCommentResult(CommentSavedMsg{
		PRNumber: 1,
		Context: gh.ReviewContext{
			PullRequestID: "PR_kwDO123",
			CommitOID:     "deadbeef",
		},
		ReviewID: "PRR_kwDO456",
		Comment: gh.ReviewComment{
			Path: "a.txt",
			Body: "body",
			Line: 1,
			Side: gh.DiffSideRight,
		},
		Err: errors.New("add failed"),
	})

	if state.Review.ReviewID != "PRR_kwDO456" {
		t.Fatalf("got %q, want %q", state.Review.ReviewID, "PRR_kwDO456")
	}
	if state.Review.PullRequestID != "PR_kwDO123" {
		t.Fatalf("got %q, want %q", state.Review.PullRequestID, "PR_kwDO123")
	}
	if state.Review.CommitOID != "deadbeef" {
		t.Fatalf("got %q, want %q", state.Review.CommitOID, "deadbeef")
	}
	if len(state.Review.Comments) != 0 {
		t.Fatalf("got %d comments, want 0", len(state.Review.Comments))
	}
	if state.Review.Notice != "add failed" {
		t.Fatalf("got %q, want %q", state.Review.Notice, "add failed")
	}
}

func TestHandleDeleteComment_WithCommentID(t *testing.T) {
	mc := &testmock.GHClient{}
	c, state, _ := setupControllerWithPR(mc, reviewstub.Selection{})
	state.SetReviewContext(1, "PR_1", "abc", "PRR_1")
	state.AddReviewComment(core.ReviewComment{CommentID: "IC_1", Path: "a.go", Body: "hi", Line: 1})
	state.Review.SelectedCommentIdx = 0

	cmd := c.HandleDeleteComment()
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}
	// State should NOT be mutated yet (API call is pending)
	if len(state.Review.Comments) != 1 {
		t.Errorf("comment should still be in state before API succeeds, got %d", len(state.Review.Comments))
	}
	msg := cmd().(CommentDeletedMsg)
	if msg.Err != nil {
		t.Fatalf("unexpected error: %v", msg.Err)
	}
	if len(mc.DeletedComments) != 1 || mc.DeletedComments[0] != "IC_1" {
		t.Errorf("got deleted %v, want [IC_1]", mc.DeletedComments)
	}
}

func TestApplyDeleteCommentResult_SuccessDeletesFromState(t *testing.T) {
	mc := &testmock.GHClient{}
	c, state, _ := setupControllerWithPR(mc, reviewstub.Selection{})
	state.SetReviewContext(1, "PR_1", "abc", "PRR_1")
	state.AddReviewComment(core.ReviewComment{CommentID: "IC_1", Path: "a.go", Body: "hi", Line: 1})
	state.Review.SelectedCommentIdx = 0
	state.BeginReviewLoad()

	c.ApplyDeleteCommentResult(CommentDeletedMsg{CommentID: "IC_1"})

	if len(state.Review.Comments) != 0 {
		t.Errorf("expected comment removed from state, got %d", len(state.Review.Comments))
	}
	if state.Review.Notice != "Comment deleted." {
		t.Errorf("got %q, want %q", state.Review.Notice, "Comment deleted.")
	}
}

func TestHandleDeleteComment_WithoutCommentID(t *testing.T) {
	mc := &testmock.GHClient{}
	c, state, _ := setupControllerWithPR(mc, reviewstub.Selection{})
	state.SetReviewContext(1, "PR_1", "abc", "PRR_1")
	state.AddReviewComment(core.ReviewComment{Path: "a.go", Body: "hi", Line: 1})
	state.Review.SelectedCommentIdx = 0

	cmd := c.HandleDeleteComment()
	if cmd != nil {
		t.Fatal("expected nil cmd for local-only comment")
	}
	if len(state.Review.Comments) != 0 {
		t.Errorf("expected comment deleted locally, got %d", len(state.Review.Comments))
	}
}

func TestApplyDeleteCommentResult_Error(t *testing.T) {
	mc := &testmock.GHClient{}
	c, state, _ := setupControllerWithPR(mc, reviewstub.Selection{})
	state.SetReviewContext(1, "PR_1", "abc", "PRR_1")
	state.AddReviewComment(core.ReviewComment{CommentID: "IC_1", Path: "a.go", Body: "hi", Line: 1})
	state.BeginReviewLoad()

	c.ApplyDeleteCommentResult(CommentDeletedMsg{CommentID: "IC_1", Err: errors.New("network error")})

	if state.Review.Notice != "network error" {
		t.Errorf("got %q, want %q", state.Review.Notice, "network error")
	}
	// Comment should NOT be removed on error
	if len(state.Review.Comments) != 1 {
		t.Errorf("comment should remain in state on error, got %d", len(state.Review.Comments))
	}
}

func TestHandleSubmit_PassesReviewEventToClient(t *testing.T) {
	tests := []struct {
		name      string
		event     core.ReviewEvent
		wantEvent string
	}{
		{"comment", core.ReviewEventComment, "COMMENT"},
		{"approve", core.ReviewEventApprove, "APPROVE"},
		{"request_changes", core.ReviewEventRequestChanges, "REQUEST_CHANGES"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &testmock.GHClient{}
			c, state, _ := setupControllerWithPR(mc, reviewstub.Selection{})
			state.SetReviewContext(1, "PR_1", "abc", "PRR_1")
			state.AddReviewComment(core.ReviewComment{Path: "a.go", Body: "hi", Line: 1})
			state.Review.Event = tt.event

			cmd := c.HandleSubmit()
			if cmd == nil {
				t.Fatal("expected non-nil cmd")
			}
			msg := cmd().(SubmittedMsg)
			if msg.Err != nil {
				t.Fatalf("unexpected error: %v", msg.Err)
			}
			if len(mc.SubmittedReviews) != 1 {
				t.Fatalf("got %d submissions, want 1", len(mc.SubmittedReviews))
			}
			want := "PRR_1:" + tt.wantEvent + ":"
			if mc.SubmittedReviews[0] != want {
				t.Errorf("got %q, want %q", mc.SubmittedReviews[0], want)
			}
		})
	}
}
