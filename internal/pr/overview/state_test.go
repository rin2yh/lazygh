package overview

import "testing"

func TestState_ZeroValue(t *testing.T) {
	var s State
	if s.Mode() != DetailModeOverview {
		t.Errorf("Mode = %v, want %v", s.Mode(), DetailModeOverview)
	}
	if s.Content() != "" {
		t.Errorf("Content = %q, want empty", s.Content())
	}
	if s.IsFetching() {
		t.Errorf("IsFetching = true, want false")
	}
}

func TestState_StartAndStopFetching(t *testing.T) {
	tests := []struct {
		name string
		kind FetchKind
	}{
		{"PRs", FetchingPRs},
		{"detail", FetchingDetail},
		{"review", FetchingReview},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s State
			s.StartFetching(tt.kind)
			if !s.IsFetching() {
				t.Error("IsFetching = false after StartFetching")
			}
			if s.FetchKind() != tt.kind {
				t.Errorf("FetchKind = %v, want %v", s.FetchKind(), tt.kind)
			}
			s.StopFetching()
			if s.IsFetching() {
				t.Error("IsFetching = true after StopFetching")
			}
		})
	}
}

func TestState_LoadResult(t *testing.T) {
	var s State
	s.StartFetching(FetchingDetail)
	s.LoadResult("body")
	if s.IsFetching() {
		t.Error("IsFetching = true after LoadResult")
	}
	if s.Content() != "body" {
		t.Errorf("Content = %q, want %q", s.Content(), "body")
	}
}

func TestState_ShowContent_DoesNotClearFetching(t *testing.T) {
	var s State
	s.StartFetching(FetchingPRs)
	s.ShowContent("preview")
	if !s.IsFetching() {
		t.Error("IsFetching = false after ShowContent, want true")
	}
	if s.Content() != "preview" {
		t.Errorf("Content = %q, want %q", s.Content(), "preview")
	}
}

func TestState_EnterOverviewMode(t *testing.T) {
	var s State
	s.StartFetching(FetchingDetail)
	s.EnterOverviewMode()
	if s.Mode() != DetailModeOverview {
		t.Errorf("Mode = %v, want %v", s.Mode(), DetailModeOverview)
	}
	if s.IsFetching() {
		t.Error("IsFetching = true after EnterOverviewMode")
	}
}

func TestState_EnterDiffMode(t *testing.T) {
	var s State
	s.StartFetching(FetchingPRs)
	s.EnterDiffMode()
	if s.Mode() != DetailModeDiff {
		t.Errorf("Mode = %v, want %v", s.Mode(), DetailModeDiff)
	}
	if s.IsFetching() {
		t.Error("IsFetching = true after EnterDiffMode")
	}
}
