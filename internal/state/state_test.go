package state

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
				detail:  "Error loading PRs: boom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewState()
			s.BeginFetchPRs()
			s.ApplyPRsResult(tt.repo, tt.prs, tt.err)

			if s.Fetching {
				t.Fatal("prs should not be loading")
			}
			if s.Detail.Fetching != model.FetchNone {
				t.Fatalf("got %v, want %v", s.Detail.Fetching, model.FetchNone)
			}
			if s.Repo != tt.want.repo {
				t.Fatalf("got %q, want %q", s.Repo, tt.want.repo)
			}
			if len(s.Items) != tt.want.prCount {
				t.Fatalf("got %d, want %d", len(s.Items), tt.want.prCount)
			}
			if s.Detail.Content != tt.want.detail {
				t.Fatalf("got %q, want %q", s.Detail.Content, tt.want.detail)
			}
			if s.Detail.Mode != model.DetailModeOverview {
				t.Fatalf("got %v, want %v", s.Detail.Mode, model.DetailModeOverview)
			}
		})
	}
}

func TestBeginFetchPRs_OnlySetsLoadingState(t *testing.T) {
	s := NewState()
	s.Detail.Content = "keep"

	s.BeginFetchPRs()

	if !s.Fetching {
		t.Fatal("expected PRsLoading to be true")
	}
	if s.Detail.Fetching != model.FetchingPRs {
		t.Fatalf("got %v, want %v", s.Detail.Fetching, model.FetchingPRs)
	}
	if s.Detail.Content != "keep" {
		t.Fatalf("got %q, want %q", s.Detail.Content, "keep")
	}
}

func TestNavigatePRs(t *testing.T) {
	s := NewState()
	s.ApplyPRsResult("owner/repo", []model.Item{{Number: 1, Title: "one"}, {Number: 2, Title: "two"}}, nil)

	changed := s.NavigateDown()
	if !changed {
		t.Fatal("expected selection change")
	}
	if s.Selected != 1 {
		t.Fatalf("got %d, want %d", s.Selected, 1)
	}
	if s.Detail.Content != "PR #2 two\nStatus: OPEN\nAssignee: unassigned" {
		t.Fatalf("got %q, want %q", s.Detail.Content, "PR #2 two\nStatus: OPEN\nAssignee: unassigned")
	}

	changed = s.NavigateUp()
	if !changed {
		t.Fatal("expected selection change")
	}
	if s.Selected != 0 {
		t.Fatalf("got %d, want %d", s.Selected, 0)
	}
	if s.Detail.Content != "PR #1 one\nStatus: OPEN\nAssignee: unassigned" {
		t.Fatalf("got %q, want %q", s.Detail.Content, "PR #1 one\nStatus: OPEN\nAssignee: unassigned")
	}
}

func TestNavigatePRs_DiffModeDoesNotOverwriteContent(t *testing.T) {
	s := NewState()
	s.ApplyPRsResult("owner/repo", []model.Item{{Number: 1, Title: "one"}, {Number: 2, Title: "two"}}, nil)
	s.Detail.Content = "diff-body"
	s.SwitchToDiff()

	changed := s.NavigateDown()
	if !changed {
		t.Fatal("expected selection change")
	}
	if s.Selected != 1 {
		t.Fatalf("got %d, want %d", s.Selected, 1)
	}
	if s.Detail.Content != "diff-body" {
		t.Fatalf("got %q, want %q", s.Detail.Content, "diff-body")
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
			s := NewState()
			s.ApplyPRsResult("owner/repo", []model.Item{{Number: 7, Title: "Fix bug"}}, nil)
			if tt.switchDiff {
				s.SwitchToDiff()
			}
			before := s.Detail.Content

			action := s.PlanEnter(true, "")
			if action.Kind != tt.wantKind {
				t.Fatalf("got %v, want %v", action.Kind, tt.wantKind)
			}
			if action.Repo != "owner/repo" {
				t.Fatalf("got %q, want %q", action.Repo, "owner/repo")
			}
			if action.Number != 7 {
				t.Fatalf("got %d, want %d", action.Number, 7)
			}
			if s.Detail.Fetching != model.FetchingDetail {
				t.Fatalf("got %v, want %v", s.Detail.Fetching, model.FetchingDetail)
			}
			if tt.wantKind == model.EnterLoadPRDetail {
				if s.Detail.Content != before {
					t.Fatalf("got %q, want %q", s.Detail.Content, before)
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
				detail: "Error loading detail: boom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewState()
			s.Detail.Fetching = model.FetchingDetail

			s.ApplyDetailResult(tt.content, tt.err)

			if s.Detail.Fetching != model.FetchNone {
				t.Fatalf("got %v, want %v", s.Detail.Fetching, model.FetchNone)
			}
			if s.Detail.Content != tt.want.detail {
				t.Fatalf("got %q, want %q", s.Detail.Content, tt.want.detail)
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
				detail: "Error loading diff: boom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewState()
			s.Detail.Fetching = model.FetchingDetail

			s.ApplyDiffResult(tt.content, tt.err)

			if s.Detail.Fetching != model.FetchNone {
				t.Fatalf("got %v, want %v", s.Detail.Fetching, model.FetchNone)
			}
			if s.Detail.Content != tt.want.detail {
				t.Fatalf("got %q, want %q", s.Detail.Content, tt.want.detail)
			}
		})
	}
}

func TestSwitchMode(t *testing.T) {
	s := NewState()
	s.ApplyPRsResult("owner/repo", []model.Item{{Number: 1, Title: "one"}}, nil)
	s.Detail.Content = "from-overview"

	if !s.SwitchToDiff() {
		t.Fatal("expected switch to diff")
	}
	if s.Detail.Mode != model.DetailModeDiff {
		t.Fatalf("got %v, want %v", s.Detail.Mode, model.DetailModeDiff)
	}
	if !s.SwitchToOverview() {
		t.Fatal("expected switch to overview")
	}
	if s.Detail.Mode != model.DetailModeOverview {
		t.Fatalf("got %v, want %v", s.Detail.Mode, model.DetailModeOverview)
	}
	if s.Detail.Content != "PR #1 one\nStatus: OPEN\nAssignee: unassigned" {
		t.Fatalf("got %q, want %q", s.Detail.Content, "PR #1 one\nStatus: OPEN\nAssignee: unassigned")
	}
}

func TestShouldApplyDetailResult(t *testing.T) {
	s := NewState()
	s.ApplyPRsResult("owner/repo", []model.Item{{Number: 1, Title: "one"}, {Number: 2, Title: "two"}}, nil)

	if !s.ShouldApplyDetailResult(model.DetailModeOverview, 1) {
		t.Fatal("expected overview detail to apply")
	}
	if s.ShouldApplyDetailResult(model.DetailModeDiff, 1) {
		t.Fatal("expected diff detail not to apply in overview mode")
	}

	s.SwitchToDiff()
	if !s.ShouldApplyDetailResult(model.DetailModeDiff, 1) {
		t.Fatal("expected diff detail to apply")
	}
	if s.ShouldApplyDetailResult(model.DetailModeDiff, 2) {
		t.Fatal("expected different PR detail not to apply")
	}
}
