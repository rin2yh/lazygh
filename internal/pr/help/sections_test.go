package help

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rin2yh/lazygh/internal/config"
	ihelp "github.com/rin2yh/lazygh/internal/help"
)

func TestSections(t *testing.T) {
	keys := config.Default().KeyBindings
	got := Sections(keys)
	want := []ihelp.Section{
		{Title: "View", Rows: [][2]string{
			{keys.DiffLabel(), "Show Diff"},
			{keys.OverviewLabel(), "Show Overview"},
			{keys.ReloadLabel(), "Reload PR"},
			{keys.QuitLabel(), "Quit"},
		}},
		{Title: "Review", Rows: [][2]string{
			{keys.RangeLabel(), "Select Range"},
			{keys.CommentLabel(), "Add Comment"},
			{keys.SummaryLabel(), "Edit Summary"},
			{keys.SaveLabel(), "Save Comment"},
			{keys.SubmitLabel(), "Submit Review"},
			{keys.DiscardLabel(), "Discard Review"},
		}},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sections() mismatch (-want +got):\n%s", diff)
	}
}
