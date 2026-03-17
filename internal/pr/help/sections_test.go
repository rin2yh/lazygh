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

func TestSections_ContainsExpectedRows(t *testing.T) {
	keys := config.Default().KeyBindings
	sections := Sections(keys)

	tests := []struct {
		name  string
		index int
		want  []string
	}{
		{"View", 0, []string{"Show Diff", "Show Overview", "Reload PR", "Quit"}},
		{"Review", 1, []string{"Select Range", "Add Comment", "Edit Summary", "Save Comment", "Submit Review", "Discard Review"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sec := sections[tt.index]
			for i, wantDesc := range tt.want {
				if i >= len(sec.Rows) {
					t.Fatalf("section has only %d rows, want at least %d", len(sec.Rows), i+1)
				}
				if sec.Rows[i][1] != wantDesc {
					t.Errorf("row[%d] description = %q, want %q", i, sec.Rows[i][1], wantDesc)
				}
			}
		})
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
