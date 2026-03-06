package panels

import (
	"fmt"
)

type Item struct {
	Number int
	Title  string
}

type ItemFormatter func(Item) string

func FormatRepoItem(item Item) string {
	return item.Title
}

func FormatIssueItem(item Item) string {
	return fmt.Sprintf("Issue #%d %s", item.Number, item.Title)
}

func FormatPRItem(item Item) string {
	return fmt.Sprintf("PR #%d %s", item.Number, item.Title)
}

type ItemsPanel struct {
	ListPanel
	Items               []Item
	Loading             bool
	Formatter           ItemFormatter
	KeepSelectionOnBlur bool
}

func NewItemsPanel(formatter ItemFormatter, keepSelectionOnBlur bool) *ItemsPanel {
	if formatter == nil {
		formatter = FormatRepoItem
	}
	return &ItemsPanel{
		ListPanel:           NewListPanel(),
		Items:               []Item{},
		Loading:             false,
		Formatter:           formatter,
		KeepSelectionOnBlur: keepSelectionOnBlur,
	}
}

func (p *ItemsPanel) Format(item Item) string {
	if p.Formatter == nil {
		return FormatRepoItem(item)
	}
	return p.Formatter(item)
}
