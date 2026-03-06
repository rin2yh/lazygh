package panels

import (
	"fmt"

	"github.com/jesseduffield/gocui"
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

func (p *ItemsPanel) Render(v *gocui.View, active bool) {
	if p.Loading {
		v.Clear()
		_ = v.SetCursor(0, 0)
		_, _ = v.Write([]byte("Loading...\n"))
		return
	}
	p.ListPanel.Render(v, len(p.Items), p.renderRow, p.shouldShowSelection(active))
}

func (p *ItemsPanel) shouldShowSelection(active bool) bool {
	return active || p.KeepSelectionOnBlur
}

func (p *ItemsPanel) renderRow(index int) string {
	return p.Format(p.Items[index])
}

func (p *ItemsPanel) Format(item Item) string {
	if p.Formatter == nil {
		return FormatRepoItem(item)
	}
	return p.Formatter(item)
}
