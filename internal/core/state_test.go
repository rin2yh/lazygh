package core

import (
	"errors"
	"testing"
)

func TestPanelCycle(t *testing.T) {
	s := NewState()
	s.NextPanel()
	if s.ActivePanel != PanelIssues {
		t.Fatalf("got %v, want %v", s.ActivePanel, PanelIssues)
	}
	s.PrevPanel()
	if s.ActivePanel != PanelRepos {
		t.Fatalf("got %v, want %v", s.ActivePanel, PanelRepos)
	}
}

func TestApplyReposResult(t *testing.T) {
	s := NewState()
	s.ReposLoading = true
	s.ApplyReposResult([]string{"owner/repo1"}, nil)
	if s.ReposLoading {
		t.Fatal("repos loading should be false")
	}
	if !s.ReposLoaded {
		t.Fatal("repos loaded should be true")
	}
	if len(s.Repos) != 1 {
		t.Fatalf("got %d, want 1", len(s.Repos))
	}
}

func TestApplyItemsResult(t *testing.T) {
	s := NewState()
	s.Repos = []Item{{Title: "owner/repo"}}
	s.ApplyItemsResult(
		"owner/repo",
		[]Item{{Number: 10, Title: "Issue one"}},
		[]Item{{Number: 1, Title: "Fix bug"}},
		nil,
	)
	if len(s.Issues) != 1 || len(s.PRs) != 1 {
		t.Fatal("issues/prs should be loaded")
	}
}

func TestApplyDetailResult_Error(t *testing.T) {
	s := NewState()
	s.ApplyDetailResult("", errors.New("boom"))
	if s.DetailContent == "" {
		t.Fatal("error message should be set")
	}
}

func TestPlanEnter_LoadItems(t *testing.T) {
	s := NewState()
	s.Repos = []Item{{Title: "owner/repo"}}
	action := s.PlanEnter(true, "")
	if action.Kind != EnterLoadItems {
		t.Fatalf("got %v, want %v", action.Kind, EnterLoadItems)
	}
}
