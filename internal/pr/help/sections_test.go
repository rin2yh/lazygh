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
			{keys.PrimaryLabel(config.ActionShowDiff), "Show Diff"},
			{keys.PrimaryLabel(config.ActionShowOverview), "Show Overview"},
			{keys.PrimaryLabel(config.ActionOpenSelected), "Reload PR"},
			{keys.PrimaryLabel(config.ActionQuit), "Quit"},
		}},
		{Title: "Review", Rows: [][2]string{
			{keys.PrimaryLabel(config.ActionReviewRange), "Select Range"},
			{keys.PrimaryLabel(config.ActionReviewComment), "Add Comment"},
			{keys.PrimaryLabel(config.ActionReviewSummary), "Edit Summary"},
			{keys.PrimaryLabel(config.ActionReviewSave), "Save Comment"},
			{keys.PrimaryLabel(config.ActionReviewSubmit), "Submit Review"},
			{keys.PrimaryLabel(config.ActionReviewDiscard), "Discard Review"},
		}},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sections() mismatch (-want +got):\n%s", diff)
	}
}
