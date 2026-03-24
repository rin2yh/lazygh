package review

import (
	"errors"
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	"github.com/rin2yh/lazygh/internal/pr"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
	reviewstub "github.com/rin2yh/lazygh/pkg/test/stub/pr/review"
)

func TestApplyCommentResult_PersistsPendingReviewContextOnError(t *testing.T) {
	host := &fakeHost{
		repo: "owner/repo",
		pr:   &pr.Item{Number: 1, Title: "Fix bug"},
	}
	controller := NewController(config.Default(), host, &testmock.GHClient{}, reviewstub.Selection{}, func(FocusTarget) {})

	controller.CommentResult(CommentSavedMsg{
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

	if controller.rs.ReviewID != "PRR_kwDO456" {
		t.Fatalf("got %q, want %q", controller.rs.ReviewID, "PRR_kwDO456")
	}
	if controller.rs.PullRequestID != "PR_kwDO123" {
		t.Fatalf("got %q, want %q", controller.rs.PullRequestID, "PR_kwDO123")
	}
	if controller.rs.CommitOID != "deadbeef" {
		t.Fatalf("got %q, want %q", controller.rs.CommitOID, "deadbeef")
	}
	if len(controller.rs.Comments) != 0 {
		t.Fatalf("got %d comments, want 0", len(controller.rs.Comments))
	}
	if controller.rs.Notice != "add failed" {
		t.Fatalf("got %q, want %q", controller.rs.Notice, "add failed")
	}
}

func TestHandleDeleteComment_WithCommentID(t *testing.T) {
	mc := &testmock.GHClient{}
	c, _, _ := setupControllerWithPR(mc, reviewstub.Selection{})
	c.rs.SetContext(1, "PR_1", "abc", "PRR_1")
	c.rs.AddComment(Comment{CommentID: "IC_1", Path: "a.go", Body: "hi", Line: 1})
	c.rs.SelectedCommentIdx = 0

	cmd := c.pending.HandleDeleteComment()
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}
	if len(c.rs.Comments) != 1 {
		t.Errorf("comment should still be in state before API succeeds, got %d", len(c.rs.Comments))
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
	c, _, _ := setupControllerWithPR(mc, reviewstub.Selection{})
	c.rs.SetContext(1, "PR_1", "abc", "PRR_1")
	c.rs.AddComment(Comment{CommentID: "IC_1", Path: "a.go", Body: "hi", Line: 1})

	c.DeleteCommentResult(CommentDeletedMsg{CommentID: "IC_1"})

	if len(c.rs.Comments) != 0 {
		t.Errorf("expected 0 comments after delete, got %d", len(c.rs.Comments))
	}
	if c.rs.Notice != "Comment deleted." {
		t.Errorf("got notice %q, want %q", c.rs.Notice, "Comment deleted.")
	}
}
