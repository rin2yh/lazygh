package review

import (
	"testing"
)

func TestNewReviewState(t *testing.T) {
	rs := newReviewState()
	if rs.EditingCommentIdx != noEditingComment {
		t.Fatalf("got %d, want %d", rs.EditingCommentIdx, noEditingComment)
	}
	if rs.Comments == nil {
		t.Fatal("expected non-nil Comments slice")
	}
	if len(rs.Comments) != 0 {
		t.Fatalf("got %d comments, want 0", len(rs.Comments))
	}
}

func TestHasPendingReview(t *testing.T) {
	rs := newReviewState()

	if rs.HasPendingReview() {
		t.Fatal("expected false when ReviewID is empty")
	}
	rs.ReviewID = "PRR_1"
	if !rs.HasPendingReview() {
		t.Fatal("expected true when ReviewID is set")
	}
}

func TestSetContext(t *testing.T) {
	rs := newReviewState()
	rs.SetContext(42, "PR_id", "abc123", "PRR_id")

	if rs.PRNumber != 42 {
		t.Fatalf("got %d, want 42", rs.PRNumber)
	}
	if rs.PullRequestID != "PR_id" {
		t.Fatalf("got %q, want %q", rs.PullRequestID, "PR_id")
	}
	if rs.CommitOID != "abc123" {
		t.Fatalf("got %q, want %q", rs.CommitOID, "abc123")
	}
	if rs.ReviewID != "PRR_id" {
		t.Fatalf("got %q, want %q", rs.ReviewID, "PRR_id")
	}
}

func TestOpenAndCloseDrawer(t *testing.T) {
	rs := newReviewState()

	rs.OpenDrawer()
	if !rs.DrawerOpen {
		t.Fatal("expected DrawerOpen = true after OpenDrawer")
	}

	rs.InputMode = InputComment
	rs.Notice = "some notice"
	rs.CloseDrawer()

	if rs.DrawerOpen {
		t.Fatal("expected DrawerOpen = false after CloseDrawer")
	}
	if rs.InputMode != InputNone {
		t.Fatalf("got InputMode %v, want InputNone", rs.InputMode)
	}
	if rs.Notice != "" {
		t.Fatalf("expected Notice cleared, got %q", rs.Notice)
	}
}

func TestBeginCommentInput(t *testing.T) {
	rs := newReviewState()
	rs.Notice = "old notice"

	rs.BeginCommentInput()

	if !rs.DrawerOpen {
		t.Fatal("expected DrawerOpen = true")
	}
	if rs.InputMode != InputComment {
		t.Fatalf("got %v, want InputComment", rs.InputMode)
	}
	if rs.Notice != "" {
		t.Fatalf("expected Notice cleared, got %q", rs.Notice)
	}
}

func TestBeginSummaryInput(t *testing.T) {
	rs := newReviewState()
	rs.Notice = "old notice"

	rs.BeginSummaryInput()

	if !rs.DrawerOpen {
		t.Fatal("expected DrawerOpen = true")
	}
	if rs.InputMode != InputSummary {
		t.Fatalf("got %v, want InputSummary", rs.InputMode)
	}
	if rs.Notice != "" {
		t.Fatalf("expected Notice cleared, got %q", rs.Notice)
	}
}

func TestStopInput(t *testing.T) {
	rs := newReviewState()
	rs.InputMode = InputComment

	rs.StopInput()

	if rs.InputMode != InputNone {
		t.Fatalf("got %v, want InputNone", rs.InputMode)
	}
}

func TestSetSummary(t *testing.T) {
	rs := newReviewState()
	rs.SetSummary("my summary")

	if rs.Summary != "my summary" {
		t.Fatalf("got %q, want %q", rs.Summary, "my summary")
	}
}

func TestAddComment(t *testing.T) {
	rs := newReviewState()
	c := Comment{Path: "main.go", Body: "looks good", Line: 10, Side: "RIGHT"}

	rs.AddComment(c)

	if len(rs.Comments) != 1 {
		t.Fatalf("got %d comments, want 1", len(rs.Comments))
	}
	if rs.Comments[0].Path != "main.go" {
		t.Fatalf("got path %q, want %q", rs.Comments[0].Path, "main.go")
	}
	if rs.SelectedCommentIdx != 0 {
		t.Fatalf("got selected %d, want 0", rs.SelectedCommentIdx)
	}
	if rs.InputMode != InputNone {
		t.Fatalf("expected InputNone after AddComment, got %v", rs.InputMode)
	}
	if rs.RangeStart != nil {
		t.Fatal("expected RangeStart cleared after AddComment")
	}
	if rs.Notice == "" {
		t.Fatal("expected notice set after AddComment")
	}
}

func TestAddComment_SelectsLastAdded(t *testing.T) {
	rs := newReviewState()
	rs.AddComment(Comment{Path: "a.go", Body: "first", Line: 1})
	rs.AddComment(Comment{Path: "b.go", Body: "second", Line: 2})

	if rs.SelectedCommentIdx != 1 {
		t.Fatalf("got %d, want 1", rs.SelectedCommentIdx)
	}
}

func TestSelectNextAndPrevComment(t *testing.T) {
	rs := newReviewState()
	rs.AddComment(Comment{Path: "a.go", Body: "c1", Line: 1})
	rs.AddComment(Comment{Path: "a.go", Body: "c2", Line: 2})
	rs.AddComment(Comment{Path: "a.go", Body: "c3", Line: 3})
	rs.SelectedCommentIdx = 0

	rs.SelectNextComment()
	if rs.SelectedCommentIdx != 1 {
		t.Fatalf("got %d, want 1", rs.SelectedCommentIdx)
	}

	rs.SelectNextComment()
	if rs.SelectedCommentIdx != 2 {
		t.Fatalf("got %d, want 2", rs.SelectedCommentIdx)
	}

	// should not go past last
	rs.SelectNextComment()
	if rs.SelectedCommentIdx != 2 {
		t.Fatalf("got %d, want 2 (at boundary)", rs.SelectedCommentIdx)
	}

	rs.SelectPrevComment()
	if rs.SelectedCommentIdx != 1 {
		t.Fatalf("got %d, want 1", rs.SelectedCommentIdx)
	}

	rs.SelectPrevComment()
	rs.SelectPrevComment()
	if rs.SelectedCommentIdx != 0 {
		t.Fatalf("got %d, want 0 (at boundary)", rs.SelectedCommentIdx)
	}
}

func TestSelectNextComment_Empty(t *testing.T) {
	rs := newReviewState()
	rs.SelectNextComment() // should not panic
}

func TestDeleteSelectedComment(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*ReviewState)
		wantLen        int
		wantSelected   int
		wantDeletedIdx int
		wantOk         bool
	}{
		{
			name:    "no comments",
			setup:   func(rs *ReviewState) {},
			wantLen: 0, wantOk: false,
		},
		{
			name: "delete only comment",
			setup: func(rs *ReviewState) {
				rs.AddComment(Comment{Body: "c1", Line: 1, Path: "a.go"})
				rs.SelectedCommentIdx = 0
			},
			wantLen: 0, wantSelected: 0, wantOk: true,
		},
		{
			name: "delete first of two",
			setup: func(rs *ReviewState) {
				rs.AddComment(Comment{Body: "c1", Line: 1, Path: "a.go"})
				rs.AddComment(Comment{Body: "c2", Line: 2, Path: "a.go"})
				rs.SelectedCommentIdx = 0
			},
			wantLen: 1, wantSelected: 0, wantOk: true,
		},
		{
			name: "delete last of two",
			setup: func(rs *ReviewState) {
				rs.AddComment(Comment{Body: "c1", Line: 1, Path: "a.go"})
				rs.AddComment(Comment{Body: "c2", Line: 2, Path: "a.go"})
				rs.SelectedCommentIdx = 1
			},
			wantLen: 1, wantSelected: 0, wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := newReviewState()
			tt.setup(rs)

			deleted, ok := rs.DeleteSelectedComment()

			if ok != tt.wantOk {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOk)
			}
			if !tt.wantOk {
				if deleted != (Comment{}) {
					t.Fatalf("expected empty comment on failure, got %+v", deleted)
				}
				return
			}
			if len(rs.Comments) != tt.wantLen {
				t.Fatalf("got %d comments, want %d", len(rs.Comments), tt.wantLen)
			}
			if rs.SelectedCommentIdx != tt.wantSelected {
				t.Fatalf("got selected %d, want %d", rs.SelectedCommentIdx, tt.wantSelected)
			}
		})
	}
}

func TestSelectedComment(t *testing.T) {
	rs := newReviewState()

	_, ok := rs.SelectedComment()
	if ok {
		t.Fatal("expected false when no comments")
	}

	rs.AddComment(Comment{Path: "a.go", Body: "hi", Line: 5})
	rs.SelectedCommentIdx = 0

	c, ok := rs.SelectedComment()
	if !ok {
		t.Fatal("expected true when comment exists")
	}
	if c.Body != "hi" {
		t.Fatalf("got %q, want %q", c.Body, "hi")
	}
}

func TestBeginEditComment(t *testing.T) {
	rs := newReviewState()
	rs.AddComment(Comment{Path: "a.go", Body: "original", Line: 1})
	rs.SelectedCommentIdx = 0

	rs.BeginEditComment()

	if rs.EditingCommentIdx != 0 {
		t.Fatalf("got %d, want 0", rs.EditingCommentIdx)
	}
	if rs.InputMode != InputComment {
		t.Fatalf("got %v, want InputComment", rs.InputMode)
	}
	if !rs.DrawerOpen {
		t.Fatal("expected DrawerOpen = true")
	}
}

func TestApplyEditComment(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*ReviewState)
		newBody  string
		wantBody string
		wantMode InputMode
		wantIdx  int
	}{
		{
			name: "valid edit",
			setup: func(rs *ReviewState) {
				rs.AddComment(Comment{Path: "a.go", Body: "old", Line: 1})
				rs.SelectedCommentIdx = 0
				rs.BeginEditComment()
			},
			newBody: "new body", wantBody: "new body",
			wantMode: InputNone, wantIdx: noEditingComment,
		},
		{
			name: "out of range index does nothing",
			setup: func(rs *ReviewState) {
				rs.EditingCommentIdx = 99
			},
			newBody: "new body", wantBody: "",
			wantMode: InputNone, wantIdx: 99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := newReviewState()
			tt.setup(rs)

			rs.ApplyEditComment(tt.newBody)

			if len(rs.Comments) > 0 && rs.Comments[0].Body != tt.wantBody {
				t.Fatalf("got body %q, want %q", rs.Comments[0].Body, tt.wantBody)
			}
			if rs.InputMode != tt.wantMode {
				t.Fatalf("got InputMode %v, want %v", rs.InputMode, tt.wantMode)
			}
			if tt.wantIdx != 99 && rs.EditingCommentIdx != tt.wantIdx {
				t.Fatalf("got EditingCommentIdx %d, want %d", rs.EditingCommentIdx, tt.wantIdx)
			}
		})
	}
}

func TestNotifyAndClearNotice(t *testing.T) {
	rs := newReviewState()

	rs.Notify("hello")
	if rs.Notice != "hello" {
		t.Fatalf("got %q, want %q", rs.Notice, "hello")
	}

	rs.ClearNotice()
	if rs.Notice != "" {
		t.Fatalf("expected empty Notice, got %q", rs.Notice)
	}
}

func TestMarkRangeStart(t *testing.T) {
	rs := newReviewState()
	anchor := Range{Path: "a.go", Index: 3, Line: 10}

	rs.MarkRangeStart(anchor)

	if rs.RangeStart == nil {
		t.Fatal("expected RangeStart set")
	}
	if *rs.RangeStart != anchor {
		t.Fatalf("got %+v, want %+v", *rs.RangeStart, anchor)
	}
	if !rs.DrawerOpen {
		t.Fatal("expected DrawerOpen = true")
	}
}

func TestClearRangeStart(t *testing.T) {
	rs := newReviewState()
	rs.MarkRangeStart(Range{Path: "a.go", Index: 1, Line: 5})

	rs.ClearRangeStart()

	if rs.RangeStart != nil {
		t.Fatal("expected RangeStart cleared")
	}
}

func TestCycleEvent(t *testing.T) {
	rs := newReviewState()

	if rs.Event != EventComment {
		t.Fatalf("initial event should be EventComment, got %v", rs.Event)
	}

	rs.CycleEvent()
	if rs.Event != EventApprove {
		t.Fatalf("got %v, want EventApprove", rs.Event)
	}

	rs.CycleEvent()
	if rs.Event != EventRequestChanges {
		t.Fatalf("got %v, want EventRequestChanges", rs.Event)
	}

	rs.CycleEvent()
	if rs.Event != EventComment {
		t.Fatalf("got %v, want EventComment", rs.Event)
	}
}

func TestReset(t *testing.T) {
	rs := newReviewState()
	rs.SetContext(1, "PR_1", "abc", "PRR_1")
	rs.AddComment(Comment{Path: "a.go", Body: "hi", Line: 1})
	rs.SetSummary("summary")
	rs.DrawerOpen = true
	rs.InputMode = InputSummary

	rs.Reset()

	if rs.ReviewID != "" {
		t.Fatalf("expected ReviewID cleared, got %q", rs.ReviewID)
	}
	if len(rs.Comments) != 0 {
		t.Fatalf("expected empty Comments, got %d", len(rs.Comments))
	}
	if rs.Summary != "" {
		t.Fatalf("expected Summary cleared, got %q", rs.Summary)
	}
	if rs.DrawerOpen {
		t.Fatal("expected DrawerOpen = false after Reset")
	}
	if rs.InputMode != InputNone {
		t.Fatalf("expected InputNone after Reset, got %v", rs.InputMode)
	}
	if rs.EditingCommentIdx != noEditingComment {
		t.Fatalf("expected EditingCommentIdx = noEditingComment, got %d", rs.EditingCommentIdx)
	}
}
