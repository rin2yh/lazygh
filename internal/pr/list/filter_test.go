package list

import (
	"strings"
	"testing"

	"github.com/rin2yh/lazygh/internal/model"
)

func TestFilterPanelLines_ReturnsLinesAndWidth(t *testing.T) {
	lines, width := FilterPanelLines(model.PRFilterOpen, 0)

	if len(lines) == 0 {
		t.Fatal("expected non-empty lines")
	}
	if width <= 0 {
		t.Fatalf("expected positive width, got %d", width)
	}
}

func TestFilterPanelLines_ContainsFilterOptions(t *testing.T) {
	filter := model.PRFilterOpen | model.PRFilterMerged
	lines, _ := FilterPanelLines(filter, 0)

	joined := strings.Join(lines, "\n")

	// All filter options should appear
	for _, opt := range model.PRFilterOptions {
		label := opt.Label()
		if !strings.Contains(joined, label) {
			t.Errorf("expected label %q in panel output", label)
		}
	}
}

func TestFilterPanelLines_CheckedAndUnchecked(t *testing.T) {
	// Only Open is enabled
	filter := model.PRFilterOpen
	lines, _ := FilterPanelLines(filter, 0)
	joined := strings.Join(lines, "\n")

	if !strings.Contains(joined, "[x]") {
		t.Error("expected [x] for enabled filter")
	}
	if !strings.Contains(joined, "[ ]") {
		t.Error("expected [ ] for disabled filter")
	}
}

func TestFilterPanelLines_CursorMarked(t *testing.T) {
	filter := model.PRFilterOpen
	cursor := 1
	lines, _ := FilterPanelLines(filter, cursor)
	joined := strings.Join(lines, "\n")

	// The cursor row should have "> " marker
	if !strings.Contains(joined, ">") {
		t.Error("expected cursor marker '>' in panel output")
	}
}

func TestFilterPanelLines_ContainsFooter(t *testing.T) {
	lines, _ := FilterPanelLines(model.PRFilterOpen, 0)
	joined := strings.Join(lines, "\n")

	if !strings.Contains(joined, "toggle") {
		t.Error("expected footer with 'toggle' hint")
	}
	if !strings.Contains(joined, "apply") {
		t.Error("expected footer with 'apply' hint")
	}
	if !strings.Contains(joined, "cancel") {
		t.Error("expected footer with 'cancel' hint")
	}
}

func TestBuildFilterContent_AllEnabled(t *testing.T) {
	filter := model.PRFilterOpen | model.PRFilterClosed | model.PRFilterMerged
	lines, maxW := buildFilterContent(filter, -1)

	if maxW <= 0 {
		t.Fatalf("expected positive maxW, got %d", maxW)
	}

	joined := strings.Join(lines, "\n")
	// All should be checked
	count := strings.Count(joined, "[x]")
	if count != len(model.PRFilterOptions) {
		t.Fatalf("got %d checked, want %d", count, len(model.PRFilterOptions))
	}
}

func TestBuildFilterContent_NoneEnabled(t *testing.T) {
	var filter model.PRFilterMask
	lines, _ := buildFilterContent(filter, -1)
	joined := strings.Join(lines, "\n")

	if strings.Contains(joined, "[x]") {
		t.Error("expected no [x] when no filter is enabled")
	}
}
