package help

import (
	"testing"

	"github.com/rin2yh/lazygh/internal/config"
)

func TestSections_ReturnsTwoSections(t *testing.T) {
	keys := config.Default().KeyBindings
	got := Sections(keys)
	if len(got) != 2 {
		t.Fatalf("got %d sections, want 2", len(got))
	}
}

func TestSections_Titles(t *testing.T) {
	keys := config.Default().KeyBindings
	got := Sections(keys)
	tests := []struct {
		index int
		title string
	}{
		{0, "View"},
		{1, "Review"},
	}
	for _, tt := range tests {
		if got[tt.index].Title != tt.title {
			t.Errorf("section[%d].Title = %q, want %q", tt.index, got[tt.index].Title, tt.title)
		}
	}
}

func TestSections_ViewContainsExpectedRows(t *testing.T) {
	keys := config.Default().KeyBindings
	sections := Sections(keys)
	viewSection := sections[0]

	wantDescriptions := []string{"Show Diff", "Show Overview", "Reload PR", "Quit"}
	for i, want := range wantDescriptions {
		if i >= len(viewSection.Rows) {
			t.Fatalf("View section has only %d rows, want at least %d", len(viewSection.Rows), i+1)
		}
		if viewSection.Rows[i][1] != want {
			t.Errorf("View row[%d] description = %q, want %q", i, viewSection.Rows[i][1], want)
		}
	}
}

func TestSections_ReviewContainsExpectedRows(t *testing.T) {
	keys := config.Default().KeyBindings
	sections := Sections(keys)
	reviewSection := sections[1]

	wantDescriptions := []string{"Select Range", "Add Comment", "Edit Summary", "Save Comment", "Submit Review", "Discard Review"}
	for i, want := range wantDescriptions {
		if i >= len(reviewSection.Rows) {
			t.Fatalf("Review section has only %d rows, want at least %d", len(reviewSection.Rows), i+1)
		}
		if reviewSection.Rows[i][1] != want {
			t.Errorf("Review row[%d] description = %q, want %q", i, reviewSection.Rows[i][1], want)
		}
	}
}

func TestSections_KeyLabelsReflectBindings(t *testing.T) {
	keys := config.Default().KeyBindings
	sections := Sections(keys)

	// View section: first row should use DiffLabel
	wantDiffLabel := keys.DiffLabel()
	if got := sections[0].Rows[0][0]; got != wantDiffLabel {
		t.Errorf("View row[0] key label = %q, want %q", got, wantDiffLabel)
	}

	// Review section: first row should use RangeLabel
	wantRangeLabel := keys.RangeLabel()
	if got := sections[1].Rows[0][0]; got != wantRangeLabel {
		t.Errorf("Review row[0] key label = %q, want %q", got, wantRangeLabel)
	}
}
