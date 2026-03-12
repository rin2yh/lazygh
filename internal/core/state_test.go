package core

import (
	"errors"
	"testing"
)

func TestApplyPRsResult(t *testing.T) {
	s := NewState()
	s.BeginLoadPRs()
	s.ApplyPRsResult("owner/repo", []Item{{Number: 1, Title: "Fix bug"}}, nil)

	if s.PRsLoading {
		t.Fatal("prs loading should be false")
	}
	if s.Repo != "owner/repo" {
		t.Fatalf("got %q, want owner/repo", s.Repo)
	}
	if len(s.PRs) != 1 {
		t.Fatalf("got %d, want 1", len(s.PRs))
	}
	if s.DetailContent != "PR #1 Fix bug" {
		t.Fatalf("got %q, want %q", s.DetailContent, "PR #1 Fix bug")
	}
}

func TestApplyPRsResult_Empty(t *testing.T) {
	s := NewState()
	s.BeginLoadPRs()
	s.ApplyPRsResult("owner/repo", nil, nil)

	if s.DetailContent != "No pull requests" {
		t.Fatalf("got %q, want %q", s.DetailContent, "No pull requests")
	}
}

func TestApplyPRsResult_Error(t *testing.T) {
	s := NewState()
	s.BeginLoadPRs()
	s.ApplyPRsResult("", nil, errors.New("boom"))
	if s.DetailContent == "" {
		t.Fatal("error message should be set")
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
	action := s.PlanEnter(true, "")
	if action.Kind != EnterLoadPRDetail {
		t.Fatalf("got %v, want %v", action.Kind, EnterLoadPRDetail)
	}
	if action.Repo != "owner/repo" || action.Number != 7 {
		t.Fatalf("unexpected action: %+v", action)
	}
}

func TestApplyDetailResult_Error(t *testing.T) {
	s := NewState()
	s.ApplyDetailResult("", errors.New("boom"))
	if s.DetailContent == "" {
		t.Fatal("error message should be set")
	}
}
