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
				{keys.DiffLabel(), "Show Diff"},
				{keys.OverviewLabel(), "Show Overview"},
				{keys.ReloadLabel(), "Reload PR"},
				{keys.QuitLabel(), "Quit"},
			},
		},
		{
			Title: "Review",
			Rows: [][2]string{
				{keys.RangeLabel(), "Select Range"},
				{keys.CommentLabel(), "Add Comment"},
				{keys.SummaryLabel(), "Edit Summary"},
				{keys.SaveLabel(), "Save Comment"},
				{keys.SubmitLabel(), "Submit Review"},
				{keys.DiscardLabel(), "Discard Review"},
			},
		},
	}
}
