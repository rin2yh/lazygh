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
	Items     []Item
	Loading   bool
	Formatter ItemFormatter
}

func NewItemsPanel(formatter ItemFormatter) *ItemsPanel {
	if formatter == nil {
		formatter = FormatRepoItem
	}
	return &ItemsPanel{
		ListPanel: NewListPanel(),
		Items:     []Item{},
		Loading:   false,
		Formatter: formatter,
	}
}

func (p *ItemsPanel) Render(v *gocui.View) {
	if p.Loading {
		v.Clear()
		_ = v.SetCursor(0, 0)
		_, _ = v.Write([]byte("Loading...\n"))
		return
	}
	p.ListPanel.Render(v, len(p.Items), p.renderRow)
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
