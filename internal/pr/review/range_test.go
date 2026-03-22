package review

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/gh"
	testmock "github.com/rin2yh/lazygh/pkg/test/mock"
	reviewstub "github.com/rin2yh/lazygh/pkg/test/stub/pr/review"
)

func TestToggleRangeSelection_StartsAndClearsRange(t *testing.T) {
	host := &fakeHost{diffMode: true}
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
	controller := NewController(config.Default(), host, &testmock.GHClient{}, selection, func(target FocusTarget) {
		focus = target
	})

	controller.ToggleRangeSelection()
	if controller.rs.RangeStart == nil {
		t.Fatal("expected range start")
	}
	if controller.rs.RangeStart.Line != 10 {
		t.Fatalf("got %d, want %d", controller.rs.RangeStart.Line, 10)
	}
	if focus != FocusDiffContent {
		t.Fatalf("got %v, want %v", focus, FocusDiffContent)
	}

	controller.ToggleRangeSelection()
	if controller.rs.RangeStart != nil {
		t.Fatal("expected range selection cleared")
	}
}
