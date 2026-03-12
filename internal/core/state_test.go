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
				detail:  "PR #1 Fix bug",
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
				t.Fatal("prs loading should be false")
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
		})
	}
}

func TestBeginLoadPRs_OnlySetsLoadingState(t *testing.T) {
	s := NewState()
	s.DetailContent = "keep"

	s.BeginLoadPRs()

	if !s.PRsLoading {
		t.Fatal("prs loading should be true")
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

	s.NavigateDown()
	if s.PRsSelected != 1 {
		t.Fatalf("got %d, want 1", s.PRsSelected)
	}
	if s.DetailContent != "PR #2 two" {
		t.Fatalf("got %q, want %q", s.DetailContent, "PR #2 two")
	}

	s.NavigateUp()
	if s.PRsSelected != 0 {
		t.Fatalf("got %d, want 0", s.PRsSelected)
	}
	if s.DetailContent != "PR #1 one" {
		t.Fatalf("got %q, want %q", s.DetailContent, "PR #1 one")
	}
}

func TestPlanEnter_LoadPRDetail(t *testing.T) {
	s := NewState()
	s.ApplyPRsResult("owner/repo", []Item{{Number: 7, Title: "Fix bug"}}, nil)
	before := s.DetailContent
	action := s.PlanEnter(true, "")
	if action.Kind != EnterLoadPRDetail {
		t.Fatalf("got %v, want %v", action.Kind, EnterLoadPRDetail)
	}
	if action.Repo != "owner/repo" || action.Number != 7 {
		t.Fatalf("unexpected action: %+v", action)
	}
	if s.Loading != LoadingDetail {
		t.Fatalf("got %v, want %v", s.Loading, LoadingDetail)
	}
	if s.DetailContent != before {
		t.Fatalf("got %q, want %q", s.DetailContent, before)
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
