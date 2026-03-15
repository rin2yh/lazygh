package review

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
	reviewstub "github.com/rin2yh/lazygh/pkg/test/stub/review"
)

func TestBuildDraft(t *testing.T) {
	commentableLine := gh.DiffLine{
		Path:        "main.go",
		NewLine:     42,
		Side:        gh.DiffSideRight,
		Commentable: true,
	}

	tests := []struct {
		name      string
		body      string
		selection reviewstub.Selection
		rangePtr  *core.ReviewRange
		wantErr   string
		wantLine  int
		wantStart int
	}{
		{
			name:      "empty body",
			body:      "",
			selection: reviewstub.Selection{Line: commentableLine},
			wantErr:   "comment body is empty",
		},
		{
			name:      "whitespace-only body",
			body:      "   \n  ",
			selection: reviewstub.Selection{Line: commentableLine},
			wantErr:   "comment body is empty",
		},
		{
			name:      "non-commentable line",
			body:      "looks good",
			selection: reviewstub.Selection{Line: gh.DiffLine{Path: "main.go", Commentable: false}},
			wantErr:   "current line is not commentable",
		},
		{
			name:      "no path (empty selection)",
			body:      "looks good",
			selection: reviewstub.Selection{},
			wantErr:   "current line is not commentable",
		},
		{
			name:      "line with zero line numbers",
			body:      "looks good",
			selection: reviewstub.Selection{Line: gh.DiffLine{Path: "main.go", Commentable: true, Side: gh.DiffSideRight}},
			wantErr:   "comment line is invalid",
		},
		{
			name:      "valid single-line comment",
			body:      "looks good",
			selection: reviewstub.Selection{Line: commentableLine, LineIndex: 5},
			wantLine:  42,
		},
		{
			name: "LEFT side uses OldLine",
			body: "old line comment",
			selection: reviewstub.Selection{Line: gh.DiffLine{
				Path:        "main.go",
				OldLine:     10,
				NewLine:     0,
				Side:        gh.DiffSideLeft,
				Commentable: true,
			}, LineIndex: 3},
			wantLine: 10,
		},
		{
			name:      "range across different files fails",
			body:      "cross-file",
			selection: reviewstub.Selection{Line: commentableLine, LineIndex: 5},
			rangePtr:  &core.ReviewRange{Path: "other.go", Index: 2, Line: 10},
			wantErr:   "range must stay within one file",
		},
		{
			name:      "valid range comment",
			body:      "range comment",
			selection: reviewstub.Selection{Line: commentableLine, LineIndex: 5},
			rangePtr:  &core.ReviewRange{Path: "main.go", Index: 2, Line: 30, Side: "RIGHT"},
			wantLine:  42,
			wantStart: 30,
		},
		{
			name:      "range with reversed indices swaps lines",
			body:      "reversed range",
			selection: reviewstub.Selection{Line: commentableLine, LineIndex: 2},
			rangePtr:  &core.ReviewRange{Path: "main.go", Index: 5, Line: 50, Side: "RIGHT"},
			wantLine:  50,
			wantStart: 42,
		},
		{
			name:      "same index range has no StartLine",
			body:      "same line range",
			selection: reviewstub.Selection{Line: commentableLine, LineIndex: 5},
			rangePtr:  &core.ReviewRange{Path: "main.go", Index: 5, Line: 42, Side: "RIGHT"},
			wantLine:  42,
			wantStart: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := core.NewState()
			c := newComment(defaultTestConfig(), state, func(FocusTarget) {})
			c.bindSelection(tt.selection)

			got, err := c.BuildDraft(tt.body, tt.rangePtr)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("got error %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Line != tt.wantLine {
				t.Errorf("Line = %d, want %d", got.Line, tt.wantLine)
			}
			if got.StartLine != tt.wantStart {
				t.Errorf("StartLine = %d, want %d", got.StartLine, tt.wantStart)
			}
		})
	}
}

func TestHandleEditorKey_EscCancelsCommentAndClearsRange(t *testing.T) {
	state := core.NewState()
	state.SwitchToDiff()
	state.BeginReviewCommentInput()
	state.MarkReviewRangeStart(core.ReviewRange{Path: "a.txt", Index: 3, Line: 10})
	focus := FocusReviewDrawer
	controller := NewController(config.Default(), state, &testmock.GHClient{}, reviewstub.Selection{}, func(target FocusTarget) {
		focus = target
	})
	controller.SetCommentValue("draft")

	_, handled := controller.HandleEditorKey(tea.KeyMsg{Type: tea.KeyEsc})
	if !handled {
		t.Fatal("expected key handled")
	}
	if state.Review.RangeStart != nil {
		t.Fatal("expected range cleared")
	}
	if state.Review.InputMode != core.ReviewInputNone {
		t.Fatalf("got %v, want %v", state.Review.InputMode, core.ReviewInputNone)
	}
	if controller.CurrentCommentValue() != "" {
		t.Fatalf("got %q, want empty", controller.CurrentCommentValue())
	}
	if focus != FocusDiffContent {
		t.Fatalf("got %v, want %v", focus, FocusDiffContent)
	}
}
