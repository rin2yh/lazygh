package review

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rin2yh/lazygh/internal/core"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
	reviewstub "github.com/rin2yh/lazygh/pkg/test/stub/review"
)

func TestHandleEditorKey_EscCancelsCommentAndClearsRange(t *testing.T) {
	state := core.NewState()
	state.SwitchToDiff()
	state.BeginReviewCommentInput()
	state.MarkReviewRangeStart(core.ReviewRange{Path: "a.txt", Index: 3, Line: 10})
	focus := FocusReviewDrawer
	controller := NewController(state, &testmock.GHClient{}, reviewstub.Selection{}, func(target FocusTarget) {
		focus = target
	})
	controller.SetCommentValue("draft")

	handled := controller.HandleEditorKey(tea.KeyMsg{Type: tea.KeyEsc})
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
