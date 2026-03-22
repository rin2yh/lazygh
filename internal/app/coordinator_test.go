package app

import (
	"errors"
	"testing"

	"github.com/rin2yh/lazygh/internal/model"
)

func TestApplyPRsResult(t *testing.T) {
	type want struct {
		repo    string
		prCount int
		detail  string
	}

	tests := []struct {
		name string
		repo string
		prs  []model.Item
		err  error
		want want
	}{
		{
			name: "success",
			repo: "owner/repo",
			prs:  []model.Item{{Number: 1, Title: "Fix bug"}},
			want: want{
				repo:    "owner/repo",
				prCount: 1,
				detail:  "PR #1 Fix bug\nStatus: OPEN\nAssignee: unassigned",
			},
		},
		{
			name: "empty",
			repo: "owner/repo",
			want: want{
				repo:    "owner/repo",
				prCount: 0,
				detail:  "No pull requests",
			},
		},
		{
			name: "error",
			err:  errors.New("boom"),
			want: want{
				repo:    "",
				prCount: 0,
				detail:  "Error fetching PRs: boom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCoordinator()
			c.BeginFetchPRs()
			c.ApplyPRsResult(tt.repo, tt.prs, tt.err)

			if c.Fetching {
				t.Fatal("prs should not be loading")
			}
			if c.Overview.Fetching != model.FetchNone {
				t.Fatalf("got %v, want %v", c.Overview.Fetching, model.FetchNone)
			}
			if c.Repo != tt.want.repo {
				t.Fatalf("got %q, want %q", c.Repo, tt.want.repo)
			}
			if len(c.Items) != tt.want.prCount {
				t.Fatalf("got %d, want %d", len(c.Items), tt.want.prCount)
			}
			if c.Overview.Content != tt.want.detail {
				t.Fatalf("got %q, want %q", c.Overview.Content, tt.want.detail)
			}
			if c.Overview.Mode != model.DetailModeOverview {
				t.Fatalf("got %v, want %v", c.Overview.Mode, model.DetailModeOverview)
			}
		})
	}
}

func TestBeginFetchPRs_OnlySetsLoadingState(t *testing.T) {
	c := NewCoordinator()
	c.Overview.Content = "keep"

	c.BeginFetchPRs()

	if !c.Fetching {
		t.Fatal("expected PRsLoading to be true")
	}
	if c.Overview.Fetching != model.FetchingPRs {
		t.Fatalf("got %v, want %v", c.Overview.Fetching, model.FetchingPRs)
	}
	if c.Overview.Content != "keep" {
		t.Fatalf("got %q, want %q", c.Overview.Content, "keep")
	}
}

func TestNavigatePRs(t *testing.T) {
	c := NewCoordinator()
	c.ApplyPRsResult("owner/repo", []model.Item{{Number: 1, Title: "one"}, {Number: 2, Title: "two"}}, nil)

	changed := c.NavigateDown()
	if !changed {
		t.Fatal("expected selection change")
	}
	if c.Selected != 1 {
		t.Fatalf("got %d, want %d", c.Selected, 1)
	}
	if c.Overview.Content != "PR #2 two\nStatus: OPEN\nAssignee: unassigned" {
		t.Fatalf("got %q, want %q", c.Overview.Content, "PR #2 two\nStatus: OPEN\nAssignee: unassigned")
	}

	changed = c.NavigateUp()
	if !changed {
		t.Fatal("expected selection change")
	}
	if c.Selected != 0 {
		t.Fatalf("got %d, want %d", c.Selected, 0)
	}
	if c.Overview.Content != "PR #1 one\nStatus: OPEN\nAssignee: unassigned" {
		t.Fatalf("got %q, want %q", c.Overview.Content, "PR #1 one\nStatus: OPEN\nAssignee: unassigned")
	}
}

func TestNavigatePRs_DiffModeDoesNotOverwriteContent(t *testing.T) {
	c := NewCoordinator()
	c.ApplyPRsResult("owner/repo", []model.Item{{Number: 1, Title: "one"}, {Number: 2, Title: "two"}}, nil)
	c.Overview.Content = "diff-body"
	c.SwitchToDiff()

	changed := c.NavigateDown()
	if !changed {
		t.Fatal("expected selection change")
	}
	if c.Selected != 1 {
		t.Fatalf("got %d, want %d", c.Selected, 1)
	}
	if c.Overview.Content != "diff-body" {
		t.Fatalf("got %q, want %q", c.Overview.Content, "diff-body")
	}
}

func TestPlanEnter_LoadPR(t *testing.T) {
	tests := []struct {
		name       string
		switchDiff bool
		wantKind   model.EnterActionKind
	}{
		{
			name:       "overview",
			switchDiff: false,
			wantKind:   model.EnterLoadPRDetail,
		},
		{
			name:       "diff",
			switchDiff: true,
			wantKind:   model.EnterLoadPRDiff,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCoordinator()
			c.ApplyPRsResult("owner/repo", []model.Item{{Number: 7, Title: "Fix bug"}}, nil)
			if tt.switchDiff {
				c.SwitchToDiff()
			}
			before := c.Overview.Content

			action := c.PlanEnter(true)
			if action.Kind != tt.wantKind {
				t.Fatalf("got %v, want %v", action.Kind, tt.wantKind)
			}
			if action.Repo != "owner/repo" {
				t.Fatalf("got %q, want %q", action.Repo, "owner/repo")
			}
			if action.Number != 7 {
				t.Fatalf("got %d, want %d", action.Number, 7)
			}
			if c.Overview.Fetching != model.FetchingDetail {
				t.Fatalf("got %v, want %v", c.Overview.Fetching, model.FetchingDetail)
			}
			if tt.wantKind == model.EnterLoadPRDetail {
				if c.Overview.Content != before {
					t.Fatalf("got %q, want %q", c.Overview.Content, before)
				}
			}
		})
	}
}

func TestApplyDetailResult(t *testing.T) {
	type want struct {
		detail string
	}

	tests := []struct {
		name    string
		content string
		err     error
		want    want
	}{
		{
			name:    "success",
			content: "hello",
			want: want{
				detail: "hello",
			},
		},
		{
			name: "error",
			err:  errors.New("boom"),
			want: want{
				detail: "Error fetching detail: boom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCoordinator()
			c.Overview.Fetching = model.FetchingDetail

			c.ApplyDetailResult(tt.content, tt.err)

			if c.Overview.Fetching != model.FetchNone {
				t.Fatalf("got %v, want %v", c.Overview.Fetching, model.FetchNone)
			}
			if c.Overview.Content != tt.want.detail {
				t.Fatalf("got %q, want %q", c.Overview.Content, tt.want.detail)
			}
		})
	}
}

func TestApplyDiffResult(t *testing.T) {
	type want struct {
		detail string
	}

	tests := []struct {
		name    string
		content string
		err     error
		want    want
	}{
		{
			name:    "success",
			content: "diff body",
			want: want{
				detail: "diff body",
			},
		},
		{
			name: "error",
			err:  errors.New("boom"),
			want: want{
				detail: "Error fetching diff: boom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCoordinator()
			c.Overview.Fetching = model.FetchingDetail

			c.ApplyDiffResult(tt.content, tt.err)

			if c.Overview.Fetching != model.FetchNone {
				t.Fatalf("got %v, want %v", c.Overview.Fetching, model.FetchNone)
			}
			if c.Overview.Content != tt.want.detail {
				t.Fatalf("got %q, want %q", c.Overview.Content, tt.want.detail)
			}
		})
	}
}

func TestSwitchMode(t *testing.T) {
	c := NewCoordinator()
	c.ApplyPRsResult("owner/repo", []model.Item{{Number: 1, Title: "one"}}, nil)
	c.Overview.Content = "from-overview"

	if !c.SwitchToDiff() {
		t.Fatal("expected switch to diff")
	}
	if c.Overview.Mode != model.DetailModeDiff {
		t.Fatalf("got %v, want %v", c.Overview.Mode, model.DetailModeDiff)
	}
	if !c.SwitchToOverview() {
		t.Fatal("expected switch to overview")
	}
	if c.Overview.Mode != model.DetailModeOverview {
		t.Fatalf("got %v, want %v", c.Overview.Mode, model.DetailModeOverview)
	}
	if c.Overview.Content != "PR #1 one\nStatus: OPEN\nAssignee: unassigned" {
		t.Fatalf("got %q, want %q", c.Overview.Content, "PR #1 one\nStatus: OPEN\nAssignee: unassigned")
	}
}

func TestShouldApplyDetailResult(t *testing.T) {
	c := NewCoordinator()
	c.ApplyPRsResult("owner/repo", []model.Item{{Number: 1, Title: "one"}, {Number: 2, Title: "two"}}, nil)

	if !c.ShouldApplyDetailResult(model.DetailModeOverview, 1) {
		t.Fatal("expected overview detail to apply")
	}
	if c.ShouldApplyDetailResult(model.DetailModeDiff, 1) {
		t.Fatal("expected diff detail not to apply in overview mode")
	}

	c.SwitchToDiff()
	if !c.ShouldApplyDetailResult(model.DetailModeDiff, 1) {
		t.Fatal("expected diff detail to apply")
	}
	if c.ShouldApplyDetailResult(model.DetailModeDiff, 2) {
		t.Fatal("expected different PR detail not to apply")
	}
}

func TestBlocksPRSelectionChange(t *testing.T) {
	c := NewCoordinator()
	c.ApplyPRsResult("owner/repo", []model.Item{{Number: 5, Title: "x"}}, nil)

	// review hook が nil の場合はブロックしない
	if c.BlocksPRSelectionChange() {
		t.Fatal("expected no block without review hook")
	}

	// pending review なし
	hook := &fakeReviewHook{hasPending: false, prNumber: 5}
	c.SetReviewHook(hook)
	if c.BlocksPRSelectionChange() {
		t.Fatal("expected no block without pending review")
	}

	// pending review あり、同じ PR
	hook.hasPending = true
	if !c.BlocksPRSelectionChange() {
		t.Fatal("expected block with pending review for current PR")
	}

	// pending review あり、別の PR
	hook.prNumber = 99
	if c.BlocksPRSelectionChange() {
		t.Fatal("expected no block when pending review is for a different PR")
	}
}

func TestApplyPRsResult_ResetsReview(t *testing.T) {
	c := NewCoordinator()
	hook := &fakeReviewHook{}
	c.SetReviewHook(hook)

	c.ApplyPRsResult("owner/repo", []model.Item{{Number: 1}}, nil)

	if !hook.resetCalled {
		t.Fatal("expected review.Reset() to be called")
	}
}

type fakeReviewHook struct {
	hasPending  bool
	prNumber    int
	resetCalled bool
}

func (f *fakeReviewHook) HasPendingReview() bool { return f.hasPending }
func (f *fakeReviewHook) PRNumber() int          { return f.prNumber }
func (f *fakeReviewHook) Reset()                 { f.resetCalled = true }
