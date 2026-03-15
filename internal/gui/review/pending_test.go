package review

import (
	"errors"
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
	reviewstub "github.com/rin2yh/lazygh/pkg/test/stub/review"
)

func TestApplyCommentResult_PersistsPendingReviewContextOnError(t *testing.T) {
	state := core.NewState()
	state.ApplyPRsResult("owner/repo", []core.Item{{Number: 1, Title: "Fix bug"}}, nil)
	state.Loading = core.LoadingReview
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
