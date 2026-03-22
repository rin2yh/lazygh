package help

import "github.com/rin2yh/lazygh/internal/config"

// CommonSections はアプリ全体で共通のキーバインドセクションを返す。
func CommonSections(keys config.KeyBindings) []Section {
	return []Section{
		{
			Title: "Navigation",
			Rows: [][2]string{
				{keys.MoveLabel(), "Move Up/Down"},
				{keys.PanelLabel(), "Panel Prev/Next"},
				{keys.PrimaryLabel(config.ActionFocusNext), "Cycle Focus"},
				{keys.PageLabel(), "Page Up/Down"},
				{keys.TopBottomLabel(), "Go Top/Bottom"},
				{keys.Label(config.ActionCancel), "Cancel / Close"},
			},
		},
	}
}
