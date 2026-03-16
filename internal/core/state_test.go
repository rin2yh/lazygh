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

			if s.List.PRsLoading {
				t.Fatal("prs should not be loading")
			}
			if s.Detail.Loading != LoadingNone {
				t.Fatalf("got %v, want %v", s.Detail.Loading, LoadingNone)
			}
			if s.List.Repo != tt.want.repo {
				t.Fatalf("got %q, want %q", s.List.Repo, tt.want.repo)
			}
			if len(s.List.PRs) != tt.want.prCount {
				t.Fatalf("got %d, want %d", len(s.List.PRs), tt.want.prCount)
			}
			if s.Detail.Content != tt.want.detail {
				t.Fatalf("got %q, want %q", s.Detail.Content, tt.want.detail)
			}
			if s.Detail.Mode != DetailModeOverview {
				t.Fatalf("got %v, want %v", s.Detail.Mode, DetailModeOverview)
			}
		})
	}
}

func TestBeginLoadPRs_OnlySetsLoadingState(t *testing.T) {
	s := NewState()
	s.Detail.Content = "keep"

	s.BeginLoadPRs()

	if !s.List.PRsLoading {
		t.Fatal("expected PRsLoading to be true")
	}
	if s.Detail.Loading != LoadingPRs {
		t.Fatalf("got %v, want %v", s.Detail.Loading, LoadingPRs)
	}
	if s.Detail.Content != "keep" {
		t.Fatalf("got %q, want %q", s.Detail.Content, "keep")
	}
}

func TestNavigatePRs(t *testing.T) {
	s := NewState()
	s.ApplyPRsResult("owner/repo", []Item{{Number: 1, Title: "one"}, {Number: 2, Title: "two"}}, nil)

	changed := s.NavigateDown()
	if !changed {
		t.Fatal("expected selection change")
	}
	if s.List.PRsSelected != 1 {
		t.Fatalf("got %d, want %d", s.List.PRsSelected, 1)
	}
	if s.Detail.Content != "PR #2 two\nStatus: OPEN\nAssignee: unassigned" {
		t.Fatalf("got %q, want %q", s.Detail.Content, "PR #2 two\nStatus: OPEN\nAssignee: unassigned")
	}

	changed = s.NavigateUp()
	if !changed {
		t.Fatal("expected selection change")
	}
	if s.List.PRsSelected != 0 {
		t.Fatalf("got %d, want %d", s.List.PRsSelected, 0)
	}
	if s.Detail.Content != "PR #1 one\nStatus: OPEN\nAssignee: unassigned" {
		t.Fatalf("got %q, want %q", s.Detail.Content, "PR #1 one\nStatus: OPEN\nAssignee: unassigned")
	}
}

func TestNavigatePRs_DiffModeDoesNotOverwriteContent(t *testing.T) {
	s := NewState()
	s.ApplyPRsResult("owner/repo", []Item{{Number: 1, Title: "one"}, {Number: 2, Title: "two"}}, nil)
	s.Detail.Content = "diff-body"
	s.SwitchToDiff()

	changed := s.NavigateDown()
	if !changed {
		t.Fatal("expected selection change")
	}
	if s.List.PRsSelected != 1 {
		t.Fatalf("got %d, want %d", s.List.PRsSelected, 1)
	}
	if s.Detail.Content != "diff-body" {
		t.Fatalf("got %q, want %q", s.Detail.Content, "diff-body")
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
			if s.Detail.Loading != LoadingDetail {
				t.Fatalf("got %v, want %v", s.Detail.Loading, LoadingDetail)
			}
			if tt.wantKind == EnterLoadPRDetail {
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
			s.Detail.Loading = LoadingDetail

			s.ApplyDetailResult(tt.content, tt.err)

			if s.Detail.Loading != LoadingNone {
				t.Fatalf("got %v, want %v", s.Detail.Loading, LoadingNone)
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
			s.Detail.Loading = LoadingDetail

			s.ApplyDiffResult(tt.content, tt.err)

			if s.Detail.Loading != LoadingNone {
				t.Fatalf("got %v, want %v", s.Detail.Loading, LoadingNone)
			}
			if s.Detail.Content != tt.want.detail {
				t.Fatalf("got %q, want %q", s.Detail.Content, tt.want.detail)
			}
		})
	}
}

func TestSwitchMode(t *testing.T) {
	s := NewState()
	s.ApplyPRsResult("owner/repo", []Item{{Number: 1, Title: "one"}}, nil)
	s.Detail.Content = "from-overview"

	if !s.SwitchToDiff() {
		t.Fatal("expected switch to diff")
	}
	if s.Detail.Mode != DetailModeDiff {
		t.Fatalf("got %v, want %v", s.Detail.Mode, DetailModeDiff)
	}
	if !s.SwitchToOverview() {
		t.Fatal("expected switch to overview")
	}
	if s.Detail.Mode != DetailModeOverview {
		t.Fatalf("got %v, want %v", s.Detail.Mode, DetailModeOverview)
	}
	if s.Detail.Content != "PR #1 one\nStatus: OPEN\nAssignee: unassigned" {
		t.Fatalf("got %q, want %q", s.Detail.Content, "PR #1 one\nStatus: OPEN\nAssignee: unassigned")
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

func TestCycleReviewEvent(t *testing.T) {
	tests := []struct {
		start ReviewEvent
		want  ReviewEvent
	}{
		{ReviewEventComment, ReviewEventApprove},
		{ReviewEventApprove, ReviewEventRequestChanges},
		{ReviewEventRequestChanges, ReviewEventComment},
	}
	for _, tt := range tests {
		s := NewState()
		s.Review.Event = tt.start
		s.CycleReviewEvent()
		if s.Review.Event != tt.want {
			t.Errorf("start=%v: got %v, want %v", tt.start, s.Review.Event, tt.want)
		}
	}
}

func TestDeleteSelectedComment(t *testing.T) {
	tests := []struct {
		name            string
		comments        []ReviewComment
		selectedIdx     int
		wantDeleted     string
		wantCount       int
		wantSelectedIdx int
	}{
		{
			name:            "delete middle comment",
			comments:        []ReviewComment{{CommentID: "c1", Body: "a"}, {CommentID: "c2", Body: "b"}, {CommentID: "c3", Body: "c"}},
			selectedIdx:     1,
			wantDeleted:     "c2",
			wantCount:       2,
			wantSelectedIdx: 1,
		},
		{
			name:            "delete last comment",
			comments:        []ReviewComment{{CommentID: "c1", Body: "a"}, {CommentID: "c2", Body: "b"}},
			selectedIdx:     1,
			wantDeleted:     "c2",
			wantCount:       1,
			wantSelectedIdx: 0,
		},
		{
			name:            "delete only comment",
			comments:        []ReviewComment{{CommentID: "c1", Body: "a"}},
			selectedIdx:     0,
			wantDeleted:     "c1",
			wantCount:       0,
			wantSelectedIdx: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewState()
			s.Review.Comments = tt.comments
			s.Review.SelectedCommentIdx = tt.selectedIdx

			deleted, ok := s.DeleteSelectedComment()
			if !ok {
				t.Fatal("expected ok=true")
			}
			if deleted.CommentID != tt.wantDeleted {
				t.Errorf("got deleted %q, want %q", deleted.CommentID, tt.wantDeleted)
			}
			if len(s.Review.Comments) != tt.wantCount {
				t.Errorf("got %d comments, want %d", len(s.Review.Comments), tt.wantCount)
			}
			if s.Review.SelectedCommentIdx != tt.wantSelectedIdx {
				t.Errorf("got selectedIdx=%d, want %d", s.Review.SelectedCommentIdx, tt.wantSelectedIdx)
			}
		})
	}
}

func TestSelectComment(t *testing.T) {
	s := NewState()
	s.Review.Comments = []ReviewComment{{Body: "a"}, {Body: "b"}, {Body: "c"}}
	s.Review.SelectedCommentIdx = 0

	s.SelectNextComment()
	if s.Review.SelectedCommentIdx != 1 {
		t.Errorf("got %d, want 1", s.Review.SelectedCommentIdx)
	}
	s.SelectNextComment()
	s.SelectNextComment() // at boundary
	if s.Review.SelectedCommentIdx != 2 {
		t.Errorf("got %d, want 2", s.Review.SelectedCommentIdx)
	}
	s.SelectPrevComment()
	if s.Review.SelectedCommentIdx != 1 {
		t.Errorf("got %d, want 1", s.Review.SelectedCommentIdx)
	}
	s.SelectPrevComment()
	s.SelectPrevComment() // at boundary
	if s.Review.SelectedCommentIdx != 0 {
		t.Errorf("got %d, want 0", s.Review.SelectedCommentIdx)
	}
}

func TestApplyEditComment(t *testing.T) {
	s := NewState()
	s.Review.Comments = []ReviewComment{{CommentID: "c1", Body: "original"}}
	s.Review.SelectedCommentIdx = 0
	s.BeginEditComment()

	if s.Review.EditingCommentIdx != 0 {
		t.Fatalf("got EditingCommentIdx=%d, want 0", s.Review.EditingCommentIdx)
	}
	if s.Review.InputMode != ReviewInputComment {
		t.Fatalf("expected ReviewInputComment mode")
	}

	s.ApplyEditComment("updated body")

	if s.Review.Comments[0].Body != "updated body" {
		t.Errorf("got %q, want %q", s.Review.Comments[0].Body, "updated body")
	}
	if s.Review.EditingCommentIdx != -1 {
		t.Errorf("got EditingCommentIdx=%d, want -1", s.Review.EditingCommentIdx)
	}
	if s.Review.InputMode != ReviewInputNone {
		t.Errorf("expected ReviewInputNone mode after edit")
	}
}

func TestPRFilterMaskLabel(t *testing.T) {
	tests := []struct {
		mask  PRFilterMask
		label string
	}{
		{PRFilterOpen, "Open"},
		{PRFilterClosed, "Closed"},
		{PRFilterMerged, "Merged"},
		{PRFilterOpen | PRFilterClosed, "Open,Closed"},
		{PRFilterOpen | PRFilterMerged, "Open,Merged"},
		{PRFilterClosed | PRFilterMerged, "Closed,Merged"},
		{PRFilterOpen | PRFilterClosed | PRFilterMerged, "All"},
		{0, "None"},
	}
	for _, tt := range tests {
		if got := tt.mask.Label(); got != tt.label {
			t.Errorf("mask=%v: got %q, want %q", tt.mask, got, tt.label)
		}
	}
}

func TestDeleteSelectedComment(t *testing.T) {
	tests := []struct {
		name            string
		comments        []ReviewComment
		selectedIdx     int
		wantDeleted     string
		wantCount       int
		wantSelectedIdx int
	}{
		{
			name:            "delete middle comment",
			comments:        []ReviewComment{{CommentID: "c1", Body: "a"}, {CommentID: "c2", Body: "b"}, {CommentID: "c3", Body: "c"}},
			selectedIdx:     1,
			wantDeleted:     "c2",
			wantCount:       2,
			wantSelectedIdx: 1,
		},
		{
			name:            "delete last comment",
			comments:        []ReviewComment{{CommentID: "c1", Body: "a"}, {CommentID: "c2", Body: "b"}},
			selectedIdx:     1,
			wantDeleted:     "c2",
			wantCount:       1,
			wantSelectedIdx: 0,
		},
		{
			name:            "delete only comment",
			comments:        []ReviewComment{{CommentID: "c1", Body: "a"}},
			selectedIdx:     0,
			wantDeleted:     "c1",
			wantCount:       0,
			wantSelectedIdx: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewState()
			s.Review.Comments = tt.comments
			s.Review.SelectedCommentIdx = tt.selectedIdx

			deleted, ok := s.DeleteSelectedComment()
			if !ok {
				t.Fatal("expected ok=true")
			}
			if deleted.CommentID != tt.wantDeleted {
				t.Errorf("got deleted %q, want %q", deleted.CommentID, tt.wantDeleted)
			}
			if len(s.Review.Comments) != tt.wantCount {
				t.Errorf("got %d comments, want %d", len(s.Review.Comments), tt.wantCount)
			}
			if s.Review.SelectedCommentIdx != tt.wantSelectedIdx {
				t.Errorf("got selectedIdx=%d, want %d", s.Review.SelectedCommentIdx, tt.wantSelectedIdx)
			}
		})
	}
}

func TestSelectComment(t *testing.T) {
	s := NewState()
	s.Review.Comments = []ReviewComment{{Body: "a"}, {Body: "b"}, {Body: "c"}}
	s.Review.SelectedCommentIdx = 0

	s.SelectNextComment()
	if s.Review.SelectedCommentIdx != 1 {
		t.Errorf("got %d, want 1", s.Review.SelectedCommentIdx)
	}
	s.SelectNextComment()
	s.SelectNextComment() // at boundary
	if s.Review.SelectedCommentIdx != 2 {
		t.Errorf("got %d, want 2", s.Review.SelectedCommentIdx)
	}
	s.SelectPrevComment()
	if s.Review.SelectedCommentIdx != 1 {
		t.Errorf("got %d, want 1", s.Review.SelectedCommentIdx)
	}
	s.SelectPrevComment()
	s.SelectPrevComment() // at boundary
	if s.Review.SelectedCommentIdx != 0 {
		t.Errorf("got %d, want 0", s.Review.SelectedCommentIdx)
	}
}

func TestApplyEditComment(t *testing.T) {
	s := NewState()
	s.Review.Comments = []ReviewComment{{CommentID: "c1", Body: "original"}}
	s.Review.SelectedCommentIdx = 0
	s.BeginEditComment()

	if s.Review.EditingCommentIdx != 0 {
		t.Fatalf("got EditingCommentIdx=%d, want 0", s.Review.EditingCommentIdx)
	}
	if s.Review.InputMode != ReviewInputComment {
		t.Fatalf("expected ReviewInputComment mode")
	}

	s.ApplyEditComment("updated body")

	if s.Review.Comments[0].Body != "updated body" {
		t.Errorf("got %q, want %q", s.Review.Comments[0].Body, "updated body")
	}
	if s.Review.EditingCommentIdx != -1 {
		t.Errorf("got EditingCommentIdx=%d, want -1", s.Review.EditingCommentIdx)
	}
	if s.Review.InputMode != ReviewInputNone {
		t.Errorf("expected ReviewInputNone mode after edit")
	}
}

func TestReviewEventLabel(t *testing.T) {
	tests := []struct {
		event ReviewEvent
		label string
	}{
		{ReviewEventComment, "COMMENT"},
		{ReviewEventApprove, "APPROVE"},
		{ReviewEventRequestChanges, "REQUEST CHANGES"},
	}
	for _, tt := range tests {
		if got := tt.event.Label(); got != tt.label {
			t.Errorf("event=%v: got %q, want %q", tt.event, got, tt.label)
		}
	}
}
