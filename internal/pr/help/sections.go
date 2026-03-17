package help

import "github.com/rin2yh/lazygh/internal/config"

// Section は1グループ分のキーバインドヘルプ定義を表す。
type Section struct {
	Title string
	Rows  [][2]string // [key label, description]
}

// BuildSections は config.KeyBindings から PR レビュー操作のヘルプセクション一覧を返す。
func BuildSections(keys config.KeyBindings) []Section {
	return []Section{
		{
			Title: "Navigation",
			Rows: [][2]string{
				{keys.MoveLabel(), "Move Up/Down"},
				{keys.PanelLabel(), "Panel Prev/Next"},
				{keys.FocusLabel(), "Cycle Focus"},
				{keys.PageLabel(), "Page Up/Down"},
				{keys.TopBottomLabel(), "Go Top/Bottom"},
				{keys.Label(config.ActionCancel), "Cancel / Close"},
			},
		},
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
