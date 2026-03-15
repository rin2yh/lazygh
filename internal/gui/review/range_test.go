package review

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/core"
	"github.com/rin2yh/lazygh/internal/gh"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
	reviewstub "github.com/rin2yh/lazygh/pkg/test/stub/review"
)

func TestToggleRangeSelection_StartsAndClearsRange(t *testing.T) {
	state := core.NewState()
	state.SwitchToDiff()
	selection := reviewstub.Selection{
		Line: gh.DiffLine{
			Path:        "a.txt",
			NewLine:     10,
			Side:        gh.DiffSideRight,
			Commentable: true,
		},
		LineIndex: 5,
	}
	focus := FocusReviewDrawer
	controller := NewController(config.Default(), state, &testmock.GHClient{}, selection, func(target FocusTarget) {
		focus = target
	})

	controller.ToggleRangeSelection()
	if state.Review.RangeStart == nil {
		t.Fatal("expected range start")
	}
	if state.Review.RangeStart.Line != 10 {
		t.Fatalf("got %d, want %d", state.Review.RangeStart.Line, 10)
	}
	if focus != FocusDiffContent {
		t.Fatalf("got %v, want %v", focus, FocusDiffContent)
	}

	controller.ToggleRangeSelection()
	if state.Review.RangeStart != nil {
		t.Fatal("expected range selection cleared")
	}
}
