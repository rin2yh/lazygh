package core

import (
	"errors"
	"testing"
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
		prs  []Item
		err  error
		want want
	}{
		{
			name: "success",
			repo: "owner/repo",
			prs:  []Item{{Number: 1, Title: "Fix bug"}},
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
			s.BeginLoadPRs()
			s.ApplyPRsResult(tt.repo, tt.prs, tt.err)

			if s.PRsLoading {
				t.Fatal("prs should not be loading")
			}
			if s.Loading != LoadingNone {
				t.Fatalf("got %v, want %v", s.Loading, LoadingNone)
			}
			if s.Repo != tt.want.repo {
				t.Fatalf("got %q, want %q", s.Repo, tt.want.repo)
			}
			if len(s.PRs) != tt.want.prCount {
				t.Fatalf("got %d, want %d", len(s.PRs), tt.want.prCount)
			}
			if s.DetailContent != tt.want.detail {
				t.Fatalf("got %q, want %q", s.DetailContent, tt.want.detail)
			}
			if s.DetailMode != DetailModeOverview {
				t.Fatalf("got %v, want %v", s.DetailMode, DetailModeOverview)
			}
		})
	}
}

func TestBeginLoadPRs_OnlySetsLoadingState(t *testing.T) {
	s := NewState()
	s.DetailContent = "keep"

	s.BeginLoadPRs()

	if !s.PRsLoading {
		t.Fatal("expected PRsLoading to be true")
	}
	if s.Loading != LoadingPRs {
		t.Fatalf("got %v, want %v", s.Loading, LoadingPRs)
	}
	if s.DetailContent != "keep" {
		t.Fatalf("got %q, want %q", s.DetailContent, "keep")
	}
}

func TestNavigatePRs(t *testing.T) {
	s := NewState()
	s.ApplyPRsResult("owner/repo", []Item{{Number: 1, Title: "one"}, {Number: 2, Title: "two"}}, nil)

	changed := s.NavigateDown()
	if !changed {
		t.Fatal("expected selection change")
	}
	if s.PRsSelected != 1 {
		t.Fatalf("got %d, want %d", s.PRsSelected, 1)
	}
	if s.DetailContent != "PR #2 two\nStatus: OPEN\nAssignee: unassigned" {
		t.Fatalf("got %q, want %q", s.DetailContent, "PR #2 two\nStatus: OPEN\nAssignee: unassigned")
	}

	changed = s.NavigateUp()
	if !changed {
		t.Fatal("expected selection change")
	}
	if s.PRsSelected != 0 {
		t.Fatalf("got %d, want %d", s.PRsSelected, 0)
	}
	if s.DetailContent != "PR #1 one\nStatus: OPEN\nAssignee: unassigned" {
		t.Fatalf("got %q, want %q", s.DetailContent, "PR #1 one\nStatus: OPEN\nAssignee: unassigned")
	}
}

func TestNavigatePRs_DiffModeDoesNotOverwriteContent(t *testing.T) {
	s := NewState()
	s.ApplyPRsResult("owner/repo", []Item{{Number: 1, Title: "one"}, {Number: 2, Title: "two"}}, nil)
	s.DetailContent = "diff-body"
	s.SwitchToDiff()

	changed := s.NavigateDown()
	if !changed {
		t.Fatal("expected selection change")
	}
	if s.PRsSelected != 1 {
		t.Fatalf("got %d, want %d", s.PRsSelected, 1)
	}
	if s.DetailContent != "diff-body" {
		t.Fatalf("got %q, want %q", s.DetailContent, "diff-body")
	}
}

func TestPlanEnter_LoadPR(t *testing.T) {
	tests := []struct {
		name       string
		switchDiff bool
		wantKind   EnterActionKind
	}{
		{
			name:       "overview",
			switchDiff: false,
			wantKind:   EnterLoadPRDetail,
		},
		{
			name:       "diff",
			switchDiff: true,
			wantKind:   EnterLoadPRDiff,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewState()
			s.ApplyPRsResult("owner/repo", []Item{{Number: 7, Title: "Fix bug"}}, nil)
			if tt.switchDiff {
				s.SwitchToDiff()
			}
			before := s.DetailContent

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
			if s.Loading != LoadingDetail {
				t.Fatalf("got %v, want %v", s.Loading, LoadingDetail)
			}
			if tt.wantKind == EnterLoadPRDetail {
				if s.DetailContent != before {
					t.Fatalf("got %q, want %q", s.DetailContent, before)
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
			s.Loading = LoadingDetail

			s.ApplyDetailResult(tt.content, tt.err)

			if s.Loading != LoadingNone {
				t.Fatalf("got %v, want %v", s.Loading, LoadingNone)
			}
			if s.DetailContent != tt.want.detail {
				t.Fatalf("got %q, want %q", s.DetailContent, tt.want.detail)
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
			s.Loading = LoadingDetail

			s.ApplyDiffResult(tt.content, tt.err)

			if s.Loading != LoadingNone {
				t.Fatalf("got %v, want %v", s.Loading, LoadingNone)
			}
			if s.DetailContent != tt.want.detail {
				t.Fatalf("got %q, want %q", s.DetailContent, tt.want.detail)
			}
		})
	}
}

func TestSwitchMode(t *testing.T) {
	s := NewState()
	s.ApplyPRsResult("owner/repo", []Item{{Number: 1, Title: "one"}}, nil)
	s.DetailContent = "from-overview"

	if !s.SwitchToDiff() {
		t.Fatal("expected switch to diff")
	}
	if s.DetailMode != DetailModeDiff {
		t.Fatalf("got %v, want %v", s.DetailMode, DetailModeDiff)
	}
	if !s.SwitchToOverview() {
		t.Fatal("expected switch to overview")
	}
	if s.DetailMode != DetailModeOverview {
		t.Fatalf("got %v, want %v", s.DetailMode, DetailModeOverview)
	}
	if s.DetailContent != "PR #1 one\nStatus: OPEN\nAssignee: unassigned" {
		t.Fatalf("got %q, want %q", s.DetailContent, "PR #1 one\nStatus: OPEN\nAssignee: unassigned")
	}
}

func TestShouldApplyDetailResult(t *testing.T) {
	s := NewState()
	s.ApplyPRsResult("owner/repo", []Item{{Number: 1, Title: "one"}, {Number: 2, Title: "two"}}, nil)

	if !s.ShouldApplyDetailResult(DetailModeOverview, 1) {
		t.Fatal("expected overview detail to apply")
	}
	if s.ShouldApplyDetailResult(DetailModeDiff, 1) {
		t.Fatal("expected diff detail not to apply in overview mode")
	}

	s.SwitchToDiff()
	if !s.ShouldApplyDetailResult(DetailModeDiff, 1) {
		t.Fatal("expected diff detail to apply")
	}
	if s.ShouldApplyDetailResult(DetailModeDiff, 2) {
		t.Fatal("expected different PR detail not to apply")
	}
}
