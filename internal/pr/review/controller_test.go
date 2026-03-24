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

func defaultTestConfig() *config.Config { return config.Default() }

// fakeHost implements AppState for tests.
type fakeHost struct {
	repo string
	pr   *pr.Item
	// fetching call counts
	beginFetchReviewCalls int
	clearFetchingCalls    int
	diffMode              bool
}

func (h *fakeHost) SelectedPR() (pr.Item, bool) {
	if h.pr == nil {
		return pr.Item{}, false
	}
	return *h.pr, true
}
func (h *fakeHost) ListRepo() string  { return h.repo }
func (h *fakeHost) BeginFetchReview() { h.beginFetchReviewCalls++ }
func (h *fakeHost) ClearFetching()    { h.clearFetchingCalls++ }
func (h *fakeHost) IsDiffMode() bool  { return h.diffMode }

func setupControllerWithPR(client *testmock.GHClient, sel reviewstub.Selection) (*Controller, *fakeHost, *FocusTarget) {
	host := &fakeHost{
		repo: "owner/repo",
		pr:   &pr.Item{Number: 1, Title: "PR"},
	}
	focus := FocusDiffContent
	c := NewController(defaultTestConfig(), host, client, sel, func(f FocusTarget) { focus = f })
	return c, host, &focus
}

func TestHandleCommentSave_NoPRSelected(t *testing.T) {
	host := &fakeHost{}
	c := NewController(defaultTestConfig(), host, &testmock.GHClient{}, reviewstub.Selection{}, func(FocusTarget) {})

	cmd := c.saveComment()
	if cmd != nil {
		t.Fatal("expected nil cmd when no PR selected")
	}
	if c.rs.Notice != "No pull request selected." {
		t.Fatalf("got %q, want %q", c.rs.Notice, "No pull request selected.")
	}
}

func TestHandleCommentSave_BuildDraftError(t *testing.T) {
	sel := reviewstub.Selection{
		Line: gh.DiffLine{Path: "a.go", Commentable: true, Side: gh.DiffSideRight},
	}
	c, _, _ := setupControllerWithPR(&testmock.GHClient{}, sel)
	c.SetCommentValue("")

	cmd := c.saveComment()
	if cmd != nil {
		t.Fatal("expected nil cmd on draft error")
	}
	if c.rs.Notice != "comment body is empty" {
		t.Fatalf("got notice %q, want %q", c.rs.Notice, "comment body is empty")
	}
}

func TestHandleCommentSave_InvalidLine(t *testing.T) {
	sel := reviewstub.Selection{
		Line: gh.DiffLine{Path: "a.go", Commentable: false},
	}
	c, _, _ := setupControllerWithPR(&testmock.GHClient{}, sel)
	c.SetCommentValue("hello")

	cmd := c.saveComment()
	if cmd != nil {
		t.Fatal("expected nil cmd on non-commentable line")
	}
	if c.rs.Notice != "current line is not commentable" {
		t.Fatalf("got notice %q, want %q", c.rs.Notice, "current line is not commentable")
	}
}

func TestHandleSubmit_NoPendingReview(t *testing.T) {
	host := &fakeHost{}
	c := NewController(defaultTestConfig(), host, &testmock.GHClient{}, reviewstub.Selection{}, func(FocusTarget) {})

	cmd := c.submit()
	if cmd != nil {
		t.Fatal("expected nil cmd when no pending review")
	}
	if c.rs.Notice != "No pending review to submit." {
		t.Fatalf("got %q, want %q", c.rs.Notice, "No pending review to submit.")
	}
}

func TestApplySubmitResult_ErrorPreservesState(t *testing.T) {
	c, host, _ := setupControllerWithPR(&testmock.GHClient{}, reviewstub.Selection{})
	host.beginFetchReviewCalls = 0
	c.rs.SetContext(1, "PR_1", "abc123", "PRR_1")
	c.rs.AddComment(Comment{Path: "a.go", Body: "hi", Line: 10})

	c.SubmitResult(SubmittedMsg{ReviewID: "PRR_1", Err: errors.New("network error")})

	if c.rs.Notice != "network error" {
		t.Fatalf("got %q, want %q", c.rs.Notice, "network error")
	}
	if c.rs.ReviewID != "PRR_1" {
		t.Fatalf("review ID cleared on error, got %q", c.rs.ReviewID)
	}
	if len(c.rs.Comments) != 1 {
		t.Fatalf("comments cleared on error, got %d", len(c.rs.Comments))
	}
}

func TestApplySubmitResult_SuccessClearsReview(t *testing.T) {
	c, _, focus := setupControllerWithPR(&testmock.GHClient{}, reviewstub.Selection{})
	c.rs.SetContext(1, "PR_1", "abc123", "PRR_1")
	c.rs.AddComment(Comment{Path: "a.go", Body: "hi", Line: 10})

	c.SubmitResult(SubmittedMsg{ReviewID: "PRR_1"})

	if c.rs.ReviewID != "" {
		t.Fatalf("review ID not cleared, got %q", c.rs.ReviewID)
	}
	if len(c.rs.Comments) != 0 {
		t.Fatalf("comments not cleared, got %d", len(c.rs.Comments))
	}
	if *focus != FocusDiffContent {
		t.Fatalf("focus not restored to diff content, got %v", *focus)
	}
}

func TestApplyDiscardResult_ErrorPreservesState(t *testing.T) {
	c, _, _ := setupControllerWithPR(&testmock.GHClient{}, reviewstub.Selection{})
	c.rs.SetContext(1, "PR_1", "abc123", "PRR_1")
	c.rs.AddComment(Comment{Path: "a.go", Body: "hi", Line: 10})

	c.DiscardResult(DiscardedMsg{Err: errors.New("discard failed")})

	if c.rs.Notice != "discard failed" {
		t.Fatalf("got %q, want %q", c.rs.Notice, "discard failed")
	}
	if c.rs.ReviewID != "PRR_1" {
		t.Fatalf("review ID cleared on error, got %q", c.rs.ReviewID)
	}
}

func TestApplyDiscardResult_SuccessClearsReview(t *testing.T) {
	c, _, focus := setupControllerWithPR(&testmock.GHClient{}, reviewstub.Selection{})
	c.rs.SetContext(1, "PR_1", "abc123", "PRR_1")
	c.rs.AddComment(Comment{Path: "a.go", Body: "hi", Line: 10})
	c.rs.SetSummary("my summary")

	c.DiscardResult(DiscardedMsg{})

	if c.rs.ReviewID != "" {
		t.Fatalf("review ID not cleared, got %q", c.rs.ReviewID)
	}
	if len(c.rs.Comments) != 0 {
		t.Fatalf("comments not cleared, got %d", len(c.rs.Comments))
	}
	if c.SummaryValue() != "" {
		t.Fatalf("summary editor not cleared, got %q", c.SummaryValue())
	}
	if *focus != FocusDiffContent {
		t.Fatalf("focus not restored, got %v", *focus)
	}
}

func TestApplyCommentResult_SuccessAddsComment(t *testing.T) {
	c, _, _ := setupControllerWithPR(&testmock.GHClient{}, reviewstub.Selection{})

	c.CommentResult(CommentSavedMsg{
		PRNumber: 1,
		Context:  gh.ReviewContext{PullRequestID: "PR_1", CommitOID: "abc"},
		ReviewID: "PRR_1",
		Comment:  gh.ReviewComment{Path: "a.go", Body: "nice", Line: 5, Side: gh.DiffSideRight},
	})

	if len(c.rs.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(c.rs.Comments))
	}
	if c.rs.Comments[0].Body != "nice" {
		t.Fatalf("got body %q, want %q", c.rs.Comments[0].Body, "nice")
	}
	if c.rs.ReviewID != "PRR_1" {
		t.Fatalf("got review ID %q, want %q", c.rs.ReviewID, "PRR_1")
	}
}

func TestHandleSubmit_SavesSummaryIfInSummaryMode(t *testing.T) {
	c, _, _ := setupControllerWithPR(&testmock.GHClient{}, reviewstub.Selection{})
	c.rs.SetContext(1, "PR_1", "abc", "PRR_1")
	c.rs.AddComment(Comment{Path: "a.go", Body: "hi", Line: 1})
	c.rs.BeginSummaryInput()

	c.summary.Load("my summary text")

	cmd := c.submit()
	if cmd == nil {
		t.Fatal("expected non-nil cmd for submit")
	}
	if c.rs.Summary != "my summary text" {
		t.Fatalf("summary not saved, got %q", c.rs.Summary)
	}
	if c.rs.InputMode != InputNone {
		t.Fatalf("input mode not cleared, got %v", c.rs.InputMode)
	}
}

func TestMarkStaleComments_MarksOrphaned(t *testing.T) {
	c, _, _ := setupControllerWithPR(&testmock.GHClient{}, reviewstub.Selection{})
	// Add two comments: one whose position still exists, one that is orphaned.
	c.rs.AddComment(Comment{Path: "a.go", Side: string(gh.DiffSideRight), Line: 10})
	c.rs.AddComment(Comment{Path: "b.go", Side: string(gh.DiffSideRight), Line: 5})

	files := []gh.DiffFile{
		{
			Path: "a.go",
			Lines: []gh.DiffLine{
				{Path: "a.go", Kind: gh.DiffLineKindAdd, Side: gh.DiffSideRight, NewLine: 10, Commentable: true},
			},
		},
		// b.go is absent from the refreshed diff.
	}
	c.MarkStaleComments(files)

	if c.rs.Comments[0].Stale {
		t.Fatal("a.go:10 should not be stale")
	}
	if !c.rs.Comments[1].Stale {
		t.Fatal("b.go:5 should be stale (file not in diff)")
	}
}

func TestHasAnchorConflict_DifferentFile(t *testing.T) {
	sel := reviewstub.Selection{
		File: gh.DiffFile{Path: "b.go"},
		Line: gh.DiffLine{Path: "b.go", Commentable: true, Side: gh.DiffSideRight, NewLine: 1},
	}
	c, _, _ := setupControllerWithPR(&testmock.GHClient{}, sel)
	// Range start is on a.go, cursor is on b.go → conflict.
	c.rs.MarkRangeStart(Range{Path: "a.go", Index: 0, Line: 5})

	if !c.rng.HasConflict() {
		t.Fatal("expected anchor conflict when files differ")
	}
}

func TestHasAnchorConflict_SameFile(t *testing.T) {
	sel := reviewstub.Selection{
		File: gh.DiffFile{Path: "a.go"},
		Line: gh.DiffLine{Path: "a.go", Commentable: true, Side: gh.DiffSideRight, NewLine: 5},
	}
	c, _, _ := setupControllerWithPR(&testmock.GHClient{}, sel)
	c.rs.MarkRangeStart(Range{Path: "a.go", Index: 0, Line: 1})

	if c.rng.HasConflict() {
		t.Fatal("expected no conflict when files are the same")
	}
}
