package help

import (
	"github.com/rin2yh/lazygh/internal/config"
	"github.com/rin2yh/lazygh/internal/help"
)

// Sections は PR レビュー操作のキーバインドセクションを返す。
func Sections(keys config.KeyBindings) []help.Section {
	return []help.Section{
		{
			Title: "View",
			Rows: [][2]string{
				{keys.PrimaryLabel(config.ActionShowDiff), "Show Diff"},
				{keys.PrimaryLabel(config.ActionShowOverview), "Show Overview"},
				{keys.PrimaryLabel(config.ActionOpenSelected), "Reload PR"},
				{keys.PrimaryLabel(config.ActionQuit), "Quit"},
			},
		},
		{
			Title: "Review",
			Rows: [][2]string{
				{keys.PrimaryLabel(config.ActionReviewRange), "Select Range"},
				{keys.PrimaryLabel(config.ActionReviewComment), "Add Comment"},
				{keys.PrimaryLabel(config.ActionReviewSummary), "Edit Summary"},
				{keys.PrimaryLabel(config.ActionReviewSave), "Save Comment"},
				{keys.PrimaryLabel(config.ActionReviewSubmit), "Submit Review"},
				{keys.PrimaryLabel(config.ActionReviewDiscard), "Discard Review"},
			},
		},
	}
}
